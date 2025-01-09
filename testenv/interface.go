package testenv

import "context"

type TestServiceInterface interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	SetFail(bool)
	SetTimeout(bool)
}
