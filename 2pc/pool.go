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
		raw = New()
	}
	txn := raw.(*Txn)
	txn.async = async
	return txn
}

func AcquireWithJobs(async bool, jobs ...distxn.Job) *Txn {
	raw := p_.p.Get()
	if raw == nil {
		raw = NewWithJobs(jobs...)
	}
	txn := raw.(*Txn)
	txn.async = async
	for i := 0; i < len(jobs); i++ {
		txn.AddJob(jobs[i])
	}
	return txn
}

func Release(dxn *Txn) {
	if dxn == nil {
		return
	}
	dxn.Reset()
	p_.p.Put(dxn)
}

var _, _, _ = Acquire, AcquireWithJobs, Release
