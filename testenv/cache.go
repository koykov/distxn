package testenv

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

func NewTestCacheClient(cache *TestCache) *TestCacheClient {
	return &TestCacheClient{cache}
}

func (c *TestCacheClient) Put(ctx context.Context, key string, value any) error {
	_, _ = key, value
	if c.svc.fail {
		return errors.New("unexpected error")
	}
	if c.svc.timeout {
		<-ctx.Done()
		return context.DeadlineExceeded
	}
	return nil
}

func (c *TestCacheClient) Delete(_ context.Context, key string) error {
	_ = key
	return nil
}

type TestCacheJob struct {
	cln *TestCacheClient
}

func NewTestCacheJob(cln *TestCacheClient) *TestCacheJob {
	return &TestCacheJob{cln}
}

func (j *TestCacheJob) Prepare(_ context.Context) error { return nil }

func (j *TestCacheJob) Commit(ctx context.Context) error {
	return j.cln.Put(ctx, "foobar", "lorem ipsum...")
}

func (j *TestCacheJob) Rollback(ctx context.Context) error { return j.cln.Delete(ctx, "foobar") }
