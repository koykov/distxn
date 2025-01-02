package twophasecommit

import (
	"sync"

	"github.com/koykov/distnx"
)

type pool struct {
	p sync.Pool
}

var p_ pool

func Acquire() *TPC {
	raw := p_.p.Get()
	if raw == nil {
		return New()
	}
	return raw.(*TPC)
}

func AcquireWithJobs(jobs ...distnx.Job) *TPC {
	raw := p_.p.Get()
	if raw == nil {
		return NewWithJobs(jobs...)
	}
	dxn := raw.(*TPC)
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
