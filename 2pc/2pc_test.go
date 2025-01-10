package twopc

import (
	"context"
	"testing"
	"time"

	"github.com/koykov/distxn"
	"github.com/koykov/distxn/testenv"
)

func Test2PC(t *testing.T) {
	newEnv := func() (*testenv.TestEnvironment, []distxn.Job) {
		cache, storage, service := testenv.TestCache{}, testenv.TestStorage{}, testenv.TestService{}
		cacheClient, storageClient, serviceClient := testenv.NewTestCacheClient(&cache), testenv.NewTestStorageClient(&storage), testenv.NewTestServiceClient(&service)
		cacheJob, storageJob, serviceJob := testenv.NewTestCacheJob(cacheClient), testenv.NewTestStorageJob(storageClient), testenv.NewTestServiceJob(serviceClient)

		env := testenv.NewTestEnvironment(&cache, &storage, &service)
		return env, []distxn.Job{cacheJob, storageJob, serviceJob}
	}
	t.Run("success", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		env, jobs := newEnv()
		_ = env.Setup(ctx)
		defer func() { _ = env.TearDown(ctx) }()

		txn := NewWithJobs(jobs...)
		xctx, xcacnel := context.WithTimeout(context.Background(), time.Millisecond*10)
		defer xcacnel()
		if err := txn.Execute(xctx); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("fail", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		env, jobs := newEnv()
		(*env)[2].SetFail(true)
		_ = env.Setup(ctx)
		defer func() { _ = env.TearDown(ctx) }()

		txn := NewWithJobs(jobs...)
		xctx, xcacnel := context.WithCancel(context.Background())
		defer xcacnel()
		if err := txn.Execute(xctx); err.Error() != "unexpected error" {
			t.Fatal(err)
		}
	})
	t.Run("timeout", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		env, jobs := newEnv()
		(*env)[1].SetTimeout(true)
		_ = env.Setup(ctx)
		defer func() { _ = env.TearDown(ctx) }()

		txn := NewWithJobs(jobs...)
		xctx, xcacnel := context.WithTimeout(context.Background(), time.Millisecond*10)
		defer xcacnel()
		if err := txn.Execute(xctx); err != context.DeadlineExceeded {
			t.Fatal(err)
		}
	})
}

func Test2PCAsync(t *testing.T) {
	newEnv := func() (*testenv.TestEnvironment, []distxn.Job) {
		cache, storage, service := testenv.TestCache{}, testenv.TestStorage{}, testenv.TestService{}
		cacheClient, storageClient, serviceClient := testenv.NewTestCacheClient(&cache), testenv.NewTestStorageClient(&storage), testenv.NewTestServiceClient(&service)
		cacheJob, storageJob, serviceJob := testenv.NewTestCacheJob(cacheClient), testenv.NewTestStorageJob(storageClient), testenv.NewTestServiceJob(serviceClient)

		env := testenv.NewTestEnvironment(&cache, &storage, &service)
		return env, []distxn.Job{cacheJob, storageJob, serviceJob}
	}
	t.Run("success", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		env, jobs := newEnv()
		_ = env.Setup(ctx)
		defer func() { _ = env.TearDown(ctx) }()

		txn := NewWithJobs(jobs...).WithAsync()
		xctx, xcacnel := context.WithTimeout(context.Background(), time.Millisecond*10)
		defer xcacnel()
		if err := txn.Execute(xctx); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("fail", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		env, jobs := newEnv()
		(*env)[2].SetFail(true)
		_ = env.Setup(ctx)
		defer func() { _ = env.TearDown(ctx) }()

		txn := NewWithJobs(jobs...).WithAsync()
		xctx, xcacnel := context.WithTimeout(context.Background(), time.Millisecond*10)
		defer xcacnel()
		if err := txn.Execute(xctx); err.Error() != "unexpected error" {
			t.Fatal(err)
		}
	})
	t.Run("timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
		defer cancel()
		env, jobs := newEnv()
		(*env)[1].SetTimeout(true)
		_ = env.Setup(ctx)
		defer func() { _ = env.TearDown(ctx) }()

		txn := NewWithJobs(jobs...).WithAsync()
		xctx, xcacnel := context.WithTimeout(context.Background(), time.Millisecond*10)
		defer xcacnel()
		if err := txn.Execute(xctx); err != context.DeadlineExceeded {
			t.Fatal(err)
		}
	})
}
