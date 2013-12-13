package zfs

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"github.com/dotcloud/docker/utils"
)

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

func execCmd(name string, args ... string) (string, string, error) {
	cmd := exec.Command(name, args...)
	var outBuf bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err := cmd.Run()

	return outBuf.String(), errBuf.String(), err
}
