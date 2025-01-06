package twopc

import (
	"context"
	"math"
	"sync"
	"sync/atomic"

	"github.com/koykov/distxn"
)

type TPC struct {
	distxn.Jobs
	async bool
	buf   []distxn.Txn
}

func New(async bool) *TPC {
	return &TPC{async: async}
}

func NewWithJobs(async bool, jobs ...distxn.Job) *TPC {
	dxn := &TPC{async: async}
	for i := 0; i < len(jobs); i++ {
		dxn.AddJob(jobs[i])
	}
	return dxn
}

func (dxn *TPC) Execute(ctx context.Context) error {
	jobs := dxn.Jobs.Jobs()
	n := len(jobs)
	if cap(dxn.buf) < n {
		dxn.buf = make([]distxn.Txn, n)
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
		go func(ctx context.Context, job distxn.Job, errc chan error) {
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
		close(errc)
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
	if dxn.async {
		return dxn.asyncCommit(ctx)
	}
	return dxn.commit(ctx)
}

func (dxn *TPC) commit(ctx context.Context) error {
	for i := 0; i < len(dxn.buf); i++ {
		if err := dxn.buf[i].Commit(ctx); err != nil {
			for j := len(dxn.buf); j >= 0; j-- {
				if err := dxn.buf[j].Rollback(ctx); err != nil {
					return err
				}
			}
			return err
		}
	}
	return nil
}

func (dxn *TPC) asyncCommit(ctx context.Context) error {
	n := len(dxn.buf)
	var (
		wg   sync.WaitGroup
		errc = make(chan error)
		done = make(chan struct{})
	)
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(ctx context.Context, txn distxn.Txn, errc chan error) {
			defer wg.Done()
			if err := txn.Commit(ctx); err != nil {
				errc <- err
				return
			}
		}(ctx, dxn.buf[i], errc)
	}
	go func() {
		wg.Wait()
		close(done)
		close(errc)
	}()
	var err error
	select {
	case err1 := <-errc:
		err = err1
		return err
	case <-ctx.Done():
		err = ctx.Err()
	case <-done:
		// do nothing
	}

	if err != nil {
		for i := 0; i < n; i++ {
			_ = dxn.buf[i].Rollback(ctx)
		}
		return err
	}
	return nil
}

func (dxn *TPC) Reset() {
	dxn.Jobs.Reset()
	dxn.async = false
	dxn.buf = dxn.buf[:0]
}
