package distnx

type Saga struct {
    jobs []Job
}

func (x *Saga) AddJob(job Job) {
    x.jobs = append(x.jobs, job)
}

func (x *Saga) Commit() error {
    // todo implement me
    return nil
}

func (x *Saga) Rollback() error {
    // todo implement me
    return nil
}
