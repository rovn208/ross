package util

import (
	"os"
	"path/filepath"
)

func IsFileAlreadyExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}

func CreateDirectory(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		err = CreateDirectory(filepath.Dir(path))
		if err != nil {
			return err
		}
		// Create the directory
		err = os.Mkdir(path, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}
