package distxn

import "context"

type DistXN interface {
	AddJob(Job)
	Execute(ctx context.Context) error
}

type Job interface {
	Begin(ctx context.Context) (Txn, error)
}

type Txn interface {
	Prepare(ctx context.Context) error
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
