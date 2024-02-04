package pipeline

import (
	"fmt"
	"io/fs"
	"sync"

	"github.com/hackborn/onefunc/lock"
)

// RegisterFs adds a named file system to the registry.
// Nodes can make use of named filesystem for, example,
// in loading operations.
func RegisterFs(name string, fsys fs.FS) error {
	return regFs.Register(name, fsys)
}

func FindFs(name string) (fs.FS, bool) {
	return regFs.Find(name)
}

// registryFs stores a list of named filesystems.
type registryFs struct {
	lock    sync.Mutex
	systems map[string]fs.FS
}

func newRegistryFs() *registryFs {
	systems := make(map[string]fs.FS)
	return &registryFs{systems: systems}
}

func (r *registryFs) Register(name string, fsys fs.FS) error {
	defer lock.Locker(&r.lock).Unlock()
	if _, ok := r.systems[name]; ok {
		return fmt.Errorf(`FS "` + name + `" already registered`)
	}
	r.systems[name] = fsys
	return nil
}

func (r *registryFs) Find(name string) (fs.FS, bool) {
	defer lock.Locker(&r.lock).Unlock()
	if f, ok := r.systems[name]; ok {
		return f, ok
	}
	return nil, false
}

var (
	regFs *registryFs = newRegistryFs()
)
