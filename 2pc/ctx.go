package twophasecommit

import "github.com/koykov/distnx"

type Ctx struct {
	jobs []distnx.Job
}

func New() *Ctx {
	return &Ctx{}
}

func (ctx *Ctx) AddJob(job distnx.Job) distnx.DistNX {
	ctx.jobs = append(ctx.jobs, job)
	return ctx
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
	ctx.jobs = ctx.jobs[:0]
}
