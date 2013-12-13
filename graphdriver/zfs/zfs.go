package zfs

/*
 * This file contains the public interface of the ZFS driver
 */

import (
	"errors"
	"fmt"
	"github.com/dotcloud/docker/graphdriver"
	"os"
	"strings"
)

func init() {
	dbg("ZFS init") // XXX: This line never shows up in -D output!

	graphdriver.Register("zfs", Init)
}

type Driver struct {
	root string // Path to the root of the graph storage (as seen by Docker daemon)
	root_dataset_name string // Name of the ZFS dataset mount at 'root'
	root_mountpoint string // Filesystem mountpoint; must be the same as 'root'
}

/*
 * Initialize the driver.
 *
 * An error is returned if ZFS is not available on the system.
 */
func Init(root string) (graphdriver.Driver, error) {
	dbg("Init")

	// Check if the ZFS filesystem is present
	if err := supportsZFS(); err != nil {
		dbg("ZFS is not supported")
		return nil, err
	}

	dbg("ZFS is supported")

	dbg("root: %s", root)

	/*
	 * Check that the root path provided to us is a ZFS filesystem. Instruct the
	 * command to emit machine-readable output (-H) by leaving out the header and
	 * using TAB to separate the fields. `zfs create` disallows a TAB character in
	 * dataset's name, so there's no danger of us getting the mount-point wrong.
	 */
	outStream, errStream, err := execCmd("zfs", "list", "-H", "-o", "name,mountpoint", "-t", "filesystem", root)
	if err != nil {
		dbg("`zfs list` error: %s", err)
		dbg("`zfs list` stderr: %s", errStream)

		return nil, err // XXX We should cook a errors.New() with accurate message.
	}

	dbg("`zfs list` output: %s", outStream)

	/*
	 * Split the output on tab characters.
	 */
	outSplice := strings.FieldsFunc(outStream,
								func (r rune) bool {
									return r == '\t'
								})

	dataset_name := outSplice[0];
	mount_point := outSplice[1]
	// Strip the trailing newline character
	mount_point = strings.TrimSuffix(mount_point, "\n")
	driver := Driver{root, dataset_name, mount_point}
	dbg("status: %v", driver.Status())

	/*
	 * Now change to the directory that is the mount-point of this filesystem. The
	 * whole point of this exercise is to ensure that the filesystem can't be
	 * unmounted behind our back while we are running. The Docker daemon should not
	 * change its directory past this point, or else we lose this protection.
	 */
	if err := os.Chdir(mount_point); err != nil {
		return nil, fmt.Errorf("zfs-Init: Could not change to the mount point '%s'", mount_point)
	}

	return &driver, nil
}

func (d Driver) rootPath() string {
	return d.root
}

func (d *Driver) String() string {
	return "zfs"
}

func (d *Driver) Status() [][2]string {
	return [][2]string{
		{"Root Dir", d.root},
		{"Dataset", d.root_dataset_name},
		{"Mount Point", d.root_mountpoint},
		// TODO: Emulate AUFS driver-like output; not necessary, but see what more info can help the user.
	}
}

/*
 * This is called when unmounting the driver. The driver is supposed to unmount the
 * filesystems of all the containers that it has in its registry.
 */
func (d *Driver) Cleanup() error {
	return errors.New("zfs-Cleanup: not supported yet")
}

/*
 * Create the on-disk structures for  the container's storage. Use the parent's
 * storage contents to populate the base image of this container.
 */
func (d *Driver) Create(id string, parent string) error {
	return errors.New("zfs-Create: not supported yet")
}

/*
 * Remove the on-disk structures of the container's storage.
 */
func (d *Driver) Remove(id string) error {
	return errors.New("zfs-Remove: not supported yet")
}

/*
 * Mount the storage of the container, and return the resulting (read-write capable)
 * path to it.
 */
func (d *Driver) Get(id string) (string, error) {
	return "", errors.New("zfs-Get: not supported yet")
}

/*
 * Exists returns true if the given id is registered with this driver.
 */

func (d *Driver) Exists(id string) bool {
	return false
}

