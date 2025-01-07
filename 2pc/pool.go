package twopc

import (
	"sync"

	"github.com/koykov/distxn"
)

type pool struct {
	p sync.Pool
}

var p_ pool

func Acquire(async bool) *Txn {
	raw := p_.p.Get()
	if raw == nil {
		return New(async)
	}
	dxn := raw.(*Txn)
	dxn.async = async
	return dxn
}

func AcquireWithJobs(async bool, jobs ...distxn.Job) *Txn {
	raw := p_.p.Get()
	if raw == nil {
		return NewWithJobs(async, jobs...)
	}
	dxn := raw.(*Txn)
	dxn.async = async
	for i := 0; i < len(jobs); i++ {
		dxn.AddJob(jobs[i])
	}
	return dxn
}

func Release(dxn *Txn) {
	if dxn == nil {
		return
	}
	dxn.Reset()
	p_.p.Put(dxn)
}

var _, _, _ = Acquire, AcquireWithJobs, Release
