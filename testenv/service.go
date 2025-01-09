package testenv

import (
	"context"
	"errors"
)

type TestService struct {
	TestBase
}

type TestServiceClient struct {
	svc *TestService
}

func NewTestServiceClient(svc *TestService) *TestServiceClient {
	return &TestServiceClient{svc}
}

func (c *TestServiceClient) Save(ctx context.Context, tuple any) error {
	_ = tuple
	if c.svc.fail {
		return errors.New("unexpected error")
	}
	if c.svc.timeout {
		<-ctx.Done()
		return context.DeadlineExceeded
	}
	return nil
}

func (c *TestServiceClient) Remove(_ context.Context) error {
	return nil
}

type TestServiceJob struct {
	cln *TestServiceClient
}

func NewTestServiceJob(cln *TestServiceClient) *TestServiceJob {
	return &TestServiceJob{cln}
}

func (j *TestServiceJob) Prepare(_ context.Context) error { return nil }

func (j *TestServiceJob) Commit(ctx context.Context) error {
	return j.cln.Save(ctx, "foobar")
}

func (j *TestServiceJob) Rollback(ctx context.Context) error { return nil }
