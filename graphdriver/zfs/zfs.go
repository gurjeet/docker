package zfs

import (
	"errors"
	"github.com/dotcloud/docker/graphdriver"	
)

func init() {
	graphdriver.Register("zfs", Init)
}

func Init(root string) (graphdriver.Driver, error) {
	d := &Driver{
		root: root,
	}
	return d, return errors.New("zfs-Init: not supported yet")
}

type Driver struct {
	root string
}

func (d *Driver) String() string {
	return "zfs"
}

func (d *Driver) Status() [][2]string {
	ids, _ := loadIds(path.Join(a.rootPath(), "layers"))
	return [][2]string{
		{"Root Dir", d.root},
		// TODO: Emulate AUFS driver-like output,
	}
}

func (d *Driver) Cleanup() error {
	return errors.New("zfs-Cleanup: not supported yet")
}

func (d *Driver) Create(id string, parent string) error {
	return errors.New("zfs-Create: not supported yet")
}

func (d *Driver) Remove(id string) error {
	return errors.New("zfs-Remove: not supported yet")
}

func (d *Driver) Get(id string) (string, error) {
	return nil, errors.New("zfs-Get: not supported yet")
}

func (d *Driver) Exists(id string) bool {
	return false
}
