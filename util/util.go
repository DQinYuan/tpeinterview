package util

import (
	"sync"
	"sync/atomic"
)


type OnceInterruptable struct {
	m    sync.Mutex
	done uint32
}


func (o *OnceInterruptable) Do(f func()) {
	if atomic.LoadUint32(&o.done) == 1 {
		return
	}
	o.m.Lock()
	defer o.m.Unlock()
	if o.done == 0 {
		f()
		atomic.StoreUint32(&o.done, 1)
	}
}
