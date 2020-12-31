package ensurepkg

import (
	"errors"

	"github.com/go-test/deep"
)

// IsTrue ensures the actual value is the boolean "true".
func (c Chain) IsTrue() {
	c.t.Helper()

	actual, ok := c.actual.(bool)
	if !ok {
		c.t.Errorf("Got type %T, expected boolean", c.actual)
		return
	}

	if !actual {
		c.t.Errorf("Got false, expected true")
	}
}

// IsFalse ensures the actual value is the boolean "false".
func (c Chain) IsFalse() {
	c.t.Helper()

	actual, ok := c.actual.(bool)
	if !ok {
		c.t.Errorf("Got type %T, expected boolean", c.actual)
		return
	}

	if actual {
		c.t.Errorf("Got true, expected false")
	}
}

// Equals ensures the actual value equals the expected value.
// Equals uses deep.Equal to print easy to read diffs.
func (c Chain) Equals(expected interface{}) {
	c.t.Helper()

	deep.CompareUnexportedFields = true
	results := deep.Equal(c.actual, expected)
	if len(results) > 0 {
		errors := "Actual does not equal expected:"
		for _, result := range results {
			errors += "\n - " + result
		}

		c.t.Errorf(errors)
	}
}

// IsError ensures the actual value equals the expected error.
// IsError uses errors.Is to support Go 1.13+ error comparisons.
func (c Chain) IsError(expected error) {
	c.t.Helper()

	if c.actual == nil && expected == nil {
		return
	}

	actual, ok := c.actual.(error)
	if !ok && c.actual != nil {
		c.t.Errorf("Got type %T, expected error: \"%v\"", c.actual, expected)
		return
	}

	if !errors.Is(actual, expected) {
		c.t.Errorf("Got \"%v\", expected \"%v\"", actual, expected)
	}
}

// IsNotError ensures that the actual value is nil.
// It is analogous to IsError(nil).
func (c Chain) IsNotError() {
	c.t.Helper()
	c.IsError(nil)
}
