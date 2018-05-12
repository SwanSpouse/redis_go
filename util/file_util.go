package util

import "os"

func FileExists(path string) bool {
	if info, err := os.Stat(path); err != nil || os.IsNotExist(err) {
		return false
	} else {
		return !info.IsDir()
	}
}
