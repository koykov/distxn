package saga

import "sync"

type pool struct {
	p sync.Pool
}

var p_ pool

func Acquire() *Txn {
	ctx := p_.p.Get()
	if ctx == nil {
		return New()
	}
	return ctx.(*Txn)
}

func Release(ctx *Txn) {
	if ctx == nil {
		return
	}
	p_.p.Put(ctx)
}

var _, _ = Acquire, Release
