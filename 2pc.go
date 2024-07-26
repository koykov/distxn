package distnx

type TwoPhaseCommit struct {
    jobs
}

func (x *TwoPhaseCommit) Commit() error {
    // todo implement me
    return nil
}

func (x *TwoPhaseCommit) Rollback() error {
    // todo implement me
    return nil
}
