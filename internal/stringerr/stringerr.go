// Package stringerr creates errors to be printed by tests.
package stringerr

import (
	"fmt"
	"strings"
)

type stringError string

// Newf creates a formatted string error.
func Newf(format string, args ...interface{}) error {
	return stringError(fmt.Sprintf(format, args...))
}

// NewGroup creates a group of errors, which are each printed on their own line.
func NewGroup(prefix string, rawErrs []error) error {
	errs := make([]string, 0, len(rawErrs))
	for _, err := range rawErrs {
		errs = append(errs, err.Error())
	}

	return Newf("%s:\n - %s", prefix, strings.Join(errs, "\n - "))
}

// Error returns the error message.
func (err stringError) Error() string {
	return string(err)
}

// Is compares errors as strings.
func (err stringError) Is(target error) bool {
	return string(err) == target.Error()
}
