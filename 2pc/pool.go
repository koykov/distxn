package twophasecommit

import "sync"

type pool struct {
	p sync.Pool
}

var p_ pool

func Acquire() *TPC {
	tpc := p_.p.Get()
	if tpc == nil {
		return New()
	}
	return tpc.(*TPC)
}

func Release(tpc *TPC) {
	if tpc == nil {
		return
	}
	tpc.Reset()
	p_.p.Put(tpc)
}

var _, _ = Acquire, Release
