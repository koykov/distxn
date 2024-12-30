package distnx

type DistNX interface {
	AddJob(Job) DistNX
	Commit() error
	Rollback() error
}

type Job interface {
	Do() error
}
