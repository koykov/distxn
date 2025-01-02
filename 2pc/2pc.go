package twophasecommit

import (
	"context"

	"github.com/koykov/distnx"
)

type TPC struct {
	distnx.Jobs
}

func New() *TPC {
	return &TPC{}
}

func (tpc *TPC) Commit(ctx context.Context) (n int, err error) {
	jobs := tpc.Jobs.Jobs()
	for i := 0; i < len(jobs); i++ {
		if err = jobs[i].Commit(ctx); err != nil {
			return
		}
		n++
	}
	return
}

func (tpc *TPC) Rollback(ctx context.Context) (n int, err error) {
	jobs := tpc.Jobs.Jobs()
	for i := 0; i < len(jobs); i++ {
		if err = jobs[i].Rollback(ctx); err != nil {
			return
		}
		n++
	}
	return
}

func (tpc *TPC) Reset() {
	tpc.Jobs.Reset()
}
