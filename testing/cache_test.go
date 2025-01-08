package testing

import (
	"context"
	"errors"
)

type TestCache struct {
	TestBase
}

type TestCacheClient struct {
	svc *TestCache
}

func (c *TestCacheClient) Put(ctx context.Context, key string, value any) error {
	if c.svc.fail {
		return errors.New("unexpected error")
	}
	if c.svc.timeout {
		<-ctx.Done()
		return context.DeadlineExceeded
	}
	return nil
}

func (c *TestCacheClient) Delete(ctx context.Context, key string) error {
	return nil
}

type TestCacheJob struct {
	cln *TestCacheClient
}

func (j *TestCacheJob) Prepare(ctx context.Context) error { return nil }

func (j *TestCacheJob) Commit(ctx context.Context) error {
	return j.cln.Put(ctx, "foobar", "lorem ipsum...")
}

func (j *TestCacheJob) Rollback(ctx context.Context) error { return j.cln.Delete(ctx, "foobar") }
