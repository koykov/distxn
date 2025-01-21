package saga

import (
	"context"

	"github.com/koykov/distxn"
)

type Txn struct {
	distxn.Jobs
}

func New() *Txn {
	return &Txn{}
}

func NewWithJobs(jobs ...distxn.Job) *Txn {
	txn := New()
	for i := 0; i < len(jobs); i++ {
		txn.AddJob(jobs[i])
	}
	return txn
}

func (txn *Txn) Execute(ctx context.Context) error {
	jobs := txn.Jobs.Jobs()
	n := len(jobs)
	_ = n
	return nil
}

func (txn *Txn) Reset() {
	txn.Jobs.Reset()
}
