package saga

import "sync"

type pool struct {
	p sync.Pool
}

var p_ pool

func Acquire() *Ctx {
	ctx := p_.p.Get()
	if ctx == nil {
		return New()
	}
	return ctx.(*Ctx)
}

func Release(ctx *Ctx) {
	if ctx == nil {
		return
	}
	p_.p.Put(ctx)
}

var _, _ = Acquire, Release
