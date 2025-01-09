package testenv

import "context"

type TestEnvironment []TestServiceInterface

func NewTestEnvironment(services ...TestServiceInterface) *TestEnvironment {
	env := TestEnvironment(services[:])
	return &env
}

func (env *TestEnvironment) Setup(ctx context.Context) error {
	for i := 0; i < len(*env); i++ {
		go func(svc TestServiceInterface) { _ = svc.Start(ctx) }((*env)[i])
	}
	return nil
}

func (env *TestEnvironment) TearDown(ctx context.Context) error {
	for i := 0; i < len(*env); i++ {
		_ = (*env)[i].Stop(ctx)
	}
	return nil
}
