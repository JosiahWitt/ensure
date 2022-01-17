// Package fswrite provides a simple, mockable wrapper for dealing with FS IO.
package fswrite

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// Writable provides a simple, mockable wrapper for dealing with FS IO.
type Writable interface {
	WriteFile(filename string, data string, perm os.FileMode) error
	MkdirAll(path string, perm os.FileMode) error
	ListRecursive(dir string) ([]string, error)
	RemoveAll(paths string) error
}

// FSWrite provides a simple, mockable wrapper for dealing with FS IO.
type FSWrite struct{}

var _ Writable = &FSWrite{}

// WriteFile wraps ioutil.WriteFile.
func (*FSWrite) WriteFile(filename string, data string, perm os.FileMode) error {
	return ioutil.WriteFile(filename, []byte(data), perm)
}

// MkdirAll wraps os.MkdirAll.
func (*FSWrite) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

// ListRecursive returns every path that can be recursively found in the provided directory.
func (*FSWrite) ListRecursive(dir string) ([]string, error) {
	paths := []string{}
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		paths = append(paths, path)
		return err
	})
	if err != nil {
		return nil, err
	}

	return paths, nil
}

// RemoveAll deletes the path and any sub paths.
func (*FSWrite) RemoveAll(path string) error {
	return os.RemoveAll(path)
}
