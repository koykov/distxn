package saga

import "github.com/koykov/distxn"

type Ctx struct {
	distxn.Jobs
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
