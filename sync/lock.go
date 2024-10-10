package sync

import (
	"sync"
)

// Scoped-lock style conveniences, for one line lock/auto-unlock.

// sync.RWMutex examples:
//		defer lock.Read(&rwmutex).Unlock()
//		defer lock.Write(&rwmutex).Unlock()

// sync.Locker examples (including Mutex):
//		defer lock.Locker(locker).Unlock()

// struct ReadAuto manages an automatic RWMutex read lock.
type ReadAuto struct {
	rw *sync.RWMutex
}

func (l ReadAuto) Unlock() {
	l.rw.RUnlock()
}

func Read(rw *sync.RWMutex) ReadAuto {
	rw.RLock()
	return ReadAuto{rw}
}

// struct WriteAuto manages an automatic RWMutex write lock.
type WriteAuto struct {
	rw *sync.RWMutex
}

func (l WriteAuto) Unlock() {
	l.rw.Unlock()
}

func Write(rw *sync.RWMutex) WriteAuto {
	rw.Lock()
	return WriteAuto{rw}
}

// AutoLock manages an automatic Locker lock.
type AutoLock struct {
	m sync.Locker
}

func (l AutoLock) Unlock() {
	l.m.Unlock()
}

func Lock(m sync.Locker) AutoLock {
	m.Lock()
	return AutoLock{m}
}
