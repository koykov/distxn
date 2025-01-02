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

func (dxn *TPC) Execute(ctx context.Context) error {
	jobs := dxn.Jobs.Jobs()
	for i := 0; i < len(jobs); i++ {
		txn, err := jobs[i].Begin(ctx)
		if err != nil {
			return err
		}
		dxn.buf = append(dxn.buf, txn)
	}

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
