package distnx

type Saga struct {
	Jobs
}

func (x *Saga) Commit() error {
	// todo implement me
	return nil
}

func (x *Saga) Rollback() error {
	// todo implement me
	return nil
}
