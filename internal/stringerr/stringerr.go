// Package stringerr creates errors to be printed by tests.
package stringerr

import (
	"fmt"
	"strings"
)

// This is three spaces because the prefix is indented one by the previous level.
const levelIndentation = "   "

type stringError string

// Newf creates a formatted string error.
func Newf(format string, args ...interface{}) error {
	return stringError(fmt.Sprintf(format, args...))
}

// Error returns the error message.
func (err stringError) Error() string {
	return string(err)
}

// Is compares errors as strings.
func (err stringError) Is(target error) bool {
	return string(err) == target.Error()
}

type groupError struct {
	prefix  string
	rawErrs []error
	footer  string
}

// NewGroup creates a group of errors, which are each printed on their own line.
func NewGroup(prefix string, rawErrs []error) error {
	return NewBlock(prefix, rawErrs, "")
}

// NewBlock creates a group of errors ending with a footer. Each error is printed on its own line.
func NewBlock(prefix string, rawErrs []error, footer string) error {
	return &groupError{
		prefix:  prefix,
		rawErrs: rawErrs,
		footer:  footer,
	}
}

// Error returns the error message.
func (g *groupError) Error() string {
	return g.errorWithIndentation("")
}

func (g *groupError) errorWithIndentation(indentation string) string {
	errs := make([]string, 0, len(g.rawErrs))
	for _, err := range g.rawErrs {
		// We don't want to recurse past the first one, so we don't want to use errors.As
		if withIndentation, ok := err.(interface{ errorWithIndentation(string) string }); ok { //nolint:inamedparam
			errs = append(errs, withIndentation.errorWithIndentation(indentation+levelIndentation))
		} else {
			errs = append(errs, err.Error())
		}
	}

	msg := fmt.Sprintf("%s:\n%s - %s", g.prefix, indentation, strings.Join(errs, "\n"+indentation+" - "))
	if g.footer != "" {
		msg += "\n" + indentation + g.footer
	}

	return msg
}

// Is compares errors as strings.
func (g *groupError) Is(target error) bool {
	return g.Error() == target.Error()
}
