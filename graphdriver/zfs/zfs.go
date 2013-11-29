package zfs

import (
	"errors",
	"fmt",
	"github.com/dotcloud/docker/graphdriver",
	"os",
	"os/exec",
	"strings",
)

func init() {
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
	// Check if the ZFS filesystem is present
	if err := supportsZFS(); err != nil {
	    return nil, err
	}

	/*
	 * Trim any leading and trailing slashes from the root name. In absence of
	 * leading and trailing slash characters the zfs command below will looks only
	 * for the named dataset, and not a directory in the filesystem by that name.
	 */
	root := strings.TrimPrefix(root, "/")
	root := strings.TrimSuffix(root, "/")

	/*
	 * Check that the root path provided to us is a ZFS filesystem. Instruct the
	 * command to emit machine-readable output (-H) by leaving out the header and
	 * using TAB to separate the fields. `zfs create` disallows a TAB character in
	 * dataset's name, so there's no danger of us getting the mount-point wrong.
	 */
	cmd := exec.Command("zfs", "list", "-H", "-t", "filesystem", root)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err // XXX We should cook a errors.New() with accurate message.
	}

	output := strings.FieldsFunc(out.String(),
								func (r rune) bool {
									return r == '\t'
								})

	mount_point := output[len(output)-1]

	/*
	 * Now change to the directory that is the mount-point of this filesystem. The
	 * whole point of this exercise is to ensure that the filesystem can't be
	 * unmounted behind our back while we are running.
	 */
	if err := os.Chdir(mount_point); err != nil {
		return nil, fmt.Errorf("zfs-Init: Could not change to the mount point '%s'", mount_point)
	}

	return &Driver{root}, nil
}

/*
 * Check if the slice contains a string
 */
func sliceContainsString(list []string, a string) (bool) {
	for _, b = range list {
		if b == a {
			return true
		}
	}
	return false
}

// Check if ZFS is supported
func supportsZFS() error {
	f, err := os.Open("/proc/filesystems")
	if err != nil {
		return err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		words := strings.Fields(s.Text())
		if sliceContainsString(words, "zfs") {
			return nil
		}
	}
	return fmt.Errorf("ZFS was not found in /proc/filesystems")
}

func (a Driver) rootPath() string {
	return a.root
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
