package testenv

import (
	"context"
	"math"
)

type TestStorage struct {
	TestBase
}

type TestStorageClient struct {
	stor *TestStorage
}

func NewTestStorageClient(stor *TestStorage) *TestStorageClient {
	return &TestStorageClient{stor}
}

func (c *TestStorageClient) Beginx(_ context.Context) error {
	return nil
}

func (c *TestStorageClient) Insert(ctx context.Context, values ...any) error {
	_ = values
	return c.stor.emulate(ctx)
}

func (c *TestStorageClient) Delete(_ context.Context) error {
	return nil
}

type TestStorageJob struct {
	cln *TestStorageClient
}

func NewTestStorageJob(cln *TestStorageClient) *TestStorageJob {
	return &TestStorageJob{cln}
}

func (j *TestStorageJob) Prepare(ctx context.Context) error { return j.cln.Beginx(ctx) }

func (j *TestStorageJob) Commit(ctx context.Context) error {
	return j.cln.Insert(ctx, "lorem", "ipsum", math.Pi)
}

func (j *TestStorageJob) Rollback(ctx context.Context) error { return j.cln.Delete(ctx) }
