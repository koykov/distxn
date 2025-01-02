package distnx

import "context"

type DistNX interface {
	AddJob(Job) DistNX
	Execute(ctx context.Context) error
}

type Job interface {
	Begin(ctx context.Context) (Txn, error)
}

type Txn interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
