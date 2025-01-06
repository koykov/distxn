package twopc

import (
	"context"
	"sync"

	"github.com/koykov/distxn"
)

type TPC struct {
	distxn.Jobs
	async bool
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

	// Phase #1: prepare
	var (
		wg   sync.WaitGroup
		errc = make(chan error)
		done = make(chan struct{})
	)
	wg.Add(n)
	for i := 0; i < len(jobs); i++ {
		go func(ctx context.Context, job distxn.Job, errc chan error) {
			defer wg.Done()

			if err := job.Prepare(ctx); err != nil {
				errc <- err
				return
			}
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
	jobs := dxn.Jobs.Jobs()
	for i := 0; i < len(jobs); i++ {
		if err := jobs[i].Commit(ctx); err != nil {
			for j := len(jobs); j >= 0; j-- {
				if err := jobs[j].Rollback(ctx); err != nil {
					return err
				}
			}
			return err
		}
	}
	return nil
}

func (dxn *TPC) asyncCommit(ctx context.Context) error {
	jobs := dxn.Jobs.Jobs()
	n := len(jobs)
	var (
		wg   sync.WaitGroup
		errc = make(chan error)
		done = make(chan struct{})
	)
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(ctx context.Context, txn distxn.Job, errc chan error) {
			defer wg.Done()
			if err := txn.Commit(ctx); err != nil {
				errc <- err
				return
			}
		}(ctx, jobs[i], errc)
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
			_ = jobs[i].Rollback(ctx)
		}
		return err
	}
	return nil
}

func (dxn *TPC) Reset() {
	dxn.Jobs.Reset()
	dxn.async = false
}
