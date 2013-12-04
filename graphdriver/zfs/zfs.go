package zfs

/*
 * This file contains the public interface of the ZFS driver
 */

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/dotcloud/docker/graphdriver"
	"os"
	"os/exec"
	"strings"
)

func init() {
	dbg("ZFS init") // XXX: This line never shows up in -D output!

	graphdriver.Register("zfs", Init)
}

type Driver struct {
	root string
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
	cmd := exec.Command("zfs", "list", "-H", "-o", "mountpoint", "-t", "filesystem", root)
	var outBuf bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		dbg("`zfs list` error: %s", err)
		dbg("`zfs list` stderr: %s", errBuf.String())

		return nil, err // XXX We should cook a errors.New() with accurate message.
	}

	dbg("`zfs list` output: %s", outBuf.String())

	/*
	 * Split the output on tab characters.
	 */
	output := strings.FieldsFunc(outBuf.String(),
								func (r rune) bool {
									return r == '\t'
								})

	mount_point := output[len(output)-1]
	// Strip the trailing newline character
	mount_point = strings.TrimSuffix(mount_point, "\n")
	dbg("Mount point: %s", mount_point)

	/*
	 * Now change to the directory that is the mount-point of this filesystem. The
	 * whole point of this exercise is to ensure that the filesystem can't be
	 * unmounted behind our back while we are running. The Docker daemon should not
	 * change its directory past this point, or else we lose this protection.
	 */
	if err := os.Chdir(mount_point); err != nil {
		return nil, fmt.Errorf("zfs-Init: Could not change to the mount point '%s'", mount_point)
	}

	return &Driver{root}, nil
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
	return "", errors.New("zfs-Get: not supported yet")
}

func (d *Driver) Exists(id string) bool {
	return false
}

