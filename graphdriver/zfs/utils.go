package zfs

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"github.com/dotcloud/docker/utils"
)

/*
 * TODO: Add infrastructire to enable/disable debugging; inspired by
 * http://play.golang.org/p/mOSbdHwSYR
 */

/*
 * Check if the slice contains a string
 */
func sliceContainsString(list []string, a string) (bool) {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

/*
 * Check if ZFS is supported
 */
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

func dbg(format string, a ... interface{}) {
	utils.Debugf("[zfs] " + format, a...)
}

/*
 * TODO: Make calls of TrimSPace() optional, dependent on a parameter which is on
 * by default.
 */
func execCmd(name string, args ... string) (string, string, error) {
	cmd := exec.Command(name, args...)
	dbg("Command: %v", cmd)
	var outBuf bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err := cmd.Run()
	outString := strings.TrimSpace(outBuf.String())
	errString := strings.TrimSpace(errBuf.String())
	if outString != "" {
		dbg("outStream: %s", outString)
	}
	if errString != "" {
		dbg("errStream: %s", errString)
	}
	if err != nil {
		dbg("error: %v", err)
	}

	return outString, errString, err
}

func funcEnter() string {

	pc, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	dbg("Entering: %s", funcName)
	return funcName
}

func funcLeave(funcName string) {
	dbg("Leaving: %s", funcName)
}

func (d *Driver) getDataset(id string) string {
	return path.Join(d.root_dataset, id)
}

func (d *Driver) getPath(id string) string {
	return path.Join(d.root_mountpoint, id)
}
