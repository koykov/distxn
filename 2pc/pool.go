package twopc

import (
	"sync"

	"github.com/koykov/distnx"
)

type pool struct {
	p sync.Pool
}

var p_ pool

func Acquire(async bool) *TPC {
	raw := p_.p.Get()
	if raw == nil {
		return New(async)
	}
	dxn := raw.(*TPC)
	dxn.async = async
	return dxn
}

func AcquireWithJobs(async bool, jobs ...distnx.Job) *TPC {
	raw := p_.p.Get()
	if raw == nil {
		return NewWithJobs(async, jobs...)
	}
	dxn := raw.(*TPC)
	dxn.async = async
	for i := 0; i < len(jobs); i++ {
		dxn.AddJob(jobs[i])
	}
	return dxn
}

func Release(dxn *TPC) {
	if dxn == nil {
		return
	}
	dxn.Reset()
	p_.p.Put(dxn)
}

var _, _, _ = Acquire, AcquireWithJobs, Release
