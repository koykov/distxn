package distnx

import "context"

type DistNX interface {
	AddJob(Job) DistNX
	Commit(ctx context.Context) (int, error)
	Rollback(ctx context.Context) (int, error)
}

type Job interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
