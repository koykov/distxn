package testenv

import (
	"context"
	"errors"
)

type TestBase struct {
	fail, timeout bool
}

func (t *TestBase) Start(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

func (t *TestBase) Stop(ctx context.Context) error {
	_ = ctx
	return nil
}

func (t *TestBase) SetFail(value bool) {
	t.fail = value
}

func (t *TestBase) SetTimeout(value bool) {
	t.timeout = value
}

func (t *TestBase) emulate(ctx context.Context) error {
	if t.fail {
		return errUnexpected
	}
	if t.timeout {
		<-ctx.Done()
		return context.DeadlineExceeded
	}
	return nil
}

var errUnexpected = errors.New("unexpected error")
