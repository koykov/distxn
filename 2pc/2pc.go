package twophasecommit

import (
	"context"

	"github.com/koykov/distnx"
)

type TPC struct {
	distnx.Jobs
	buf []distnx.Txn
}

func New() *TPC {
	return &TPC{}
}

func (tpc *TPC) Execute(ctx context.Context) error {
	jobs := tpc.Jobs.Jobs()
	for i := 0; i < len(jobs); i++ {
		txn, err := jobs[i].Begin(ctx)
		if err != nil {
			return err
		}
		tpc.buf = append(tpc.buf, txn)
	}

	for i := 0; i < len(tpc.buf); i++ {
		if err := tpc.buf[i].Commit(ctx); err != nil {
			for j := i; j >= 0; j-- {
				if err := tpc.buf[j].Rollback(ctx); err != nil {
					return err
				}
			}
			return err
		}
	}
	return nil
}

func (tpc *TPC) Reset() {
	tpc.Jobs.Reset()
	tpc.buf = tpc.buf[:0]
}
