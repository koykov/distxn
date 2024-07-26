package distnx

type TwoPhaseCommit struct {
    jobs []Job
}

func (x *TwoPhaseCommit) AddJob(job Job) {
    x.jobs = append(x.jobs, job)
}

func (x *TwoPhaseCommit) Commit() error {
    // todo implement me
    return nil
}

func (x *TwoPhaseCommit) Rollback() error {
    // todo implement me
    return nil
}
