package fswrite

import (
	"io/ioutil"
	"os"
)

//nolint:golint // Ignore stuttering concerns
type FSWriteIface interface {
	WriteFile(filename string, data string, perm os.FileMode) error
	MkdirAll(path string, perm os.FileMode) error
}

type FSWrite struct{}

var _ FSWriteIface = &FSWrite{}

// WriteFile wraps ioutil.WriteFile.
func (*FSWrite) WriteFile(filename string, data string, perm os.FileMode) error {
	return ioutil.WriteFile(filename, []byte(data), perm)
}

// MkdirAll wraps os.MkdirAll.
func (*FSWrite) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}
