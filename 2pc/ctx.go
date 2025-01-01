package twophasecommit

import "github.com/koykov/distnx"

type Ctx struct {
	distnx.Jobs
}

func New() *Ctx {
	return &Ctx{}
}

func (ctx *Ctx) Commit() error {
	// todo implement me
	return nil
}

func (ctx *Ctx) Rollback() error {
	// todo implement me
	return nil
}

func (ctx *Ctx) Reset() {
	ctx.Jobs.Reset()
}
