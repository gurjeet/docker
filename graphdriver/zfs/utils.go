package zfs

import (
	"bufio"
	"fmt"
	"os"
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