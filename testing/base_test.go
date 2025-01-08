package testing

import "context"

type TestBase struct {
	fail, timeout, deny bool
}

func (t *TestBase) Start(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

func (t *TestBase) Stop(ctx context.Context) error {
	_ = ctx
	return nil
}
