package util

import (
	"os"
)

func CreateDirectoryIfNotExists(path string) error {
	if !DirectoryExists(path) {
		return os.MkdirAll(path, 0766)
	}
	return nil
}

func DirectoryExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}
