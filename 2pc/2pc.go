package twophasecommit

import (
	"context"
	"math"
	"sync"
	"sync/atomic"

	"github.com/koykov/distnx"
)

type TPC struct {
	distnx.Jobs
	buf []distnx.Txn
}

func New() *TPC {
	return &TPC{}
}

func NewWithJobs(jobs ...distnx.Job) *TPC {
	dxn := &TPC{}
	for i := 0; i < len(jobs); i++ {
		dxn.AddJob(jobs[i])
	}
	return dxn
}

func (dxn *TPC) Execute(ctx context.Context) error {
	jobs := dxn.Jobs.Jobs()
	n := len(jobs)
	if cap(dxn.buf) < n {
		dxn.buf = make([]distnx.Txn, n)
	}
	dxn.buf = dxn.buf[:n:n]

	// Phase #1: prepare
	var (
		wg   sync.WaitGroup
		errc        = make(chan error)
		done        = make(chan struct{})
		idx  uint32 = math.MaxUint32
	)
	wg.Add(n)
	for i := 0; i < len(jobs); i++ {
		go func(ctx context.Context, job distnx.Job, errc chan error) {
			defer wg.Done()
			txn, err := job.Begin(ctx)
			if err != nil {
				errc <- err
				return
			}
			if err = txn.Prepare(ctx); err != nil {
				errc <- err
				return
			}
			dxn.buf[atomic.AddUint32(&idx, 1)] = txn
		}(ctx, jobs[i], errc)
	}
	go func(wg *sync.WaitGroup) {
		wg.Wait()
		close(done)
	}(&wg)
	select {
	case err := <-errc:
		return err
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		// do nothing
	}

	// Phase #2: commit
	for i := 0; i < len(dxn.buf); i++ {
		if err := dxn.buf[i].Commit(ctx); err != nil {
			for j := i; j >= 0; j-- {
				if err := dxn.buf[j].Rollback(ctx); err != nil {
					return err
				}
			}
			return err
		}
	}
	return nil
}

func (dxn *TPC) Reset() {
	dxn.Jobs.Reset()
	dxn.buf = dxn.buf[:0]
}
