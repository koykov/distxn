package distnx

type DistNX interface {
    AddJob(Job)
    Commit() error
    Rollback() error
}

type Job interface {
    Do() error
}
