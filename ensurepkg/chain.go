package ensurepkg

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/JosiahWitt/erk"
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

// IsNil ensures the actual value is nil or a nil pointer.
func (c Chain) IsNil() {
	c.t.Helper()

	if !isNil(c.actual) {
		c.t.Errorf("Got %+v, expected nil", c.actual)
	}
}

// Equals ensures the actual value equals the expected value.
// Equals uses deep.Equal to print easy to read diffs.
func (c Chain) Equals(expected interface{}) {
	c.t.Helper()

	// If we expect nil, return early if actual is nil or if it is a nil pointer
	if expected == nil && isNil(c.actual) {
		return
	}

	deep.CompareUnexportedFields = true
	deep.NilMapsAreEmpty = true
	deep.NilSlicesAreEmpty = true
	results := deep.Equal(c.actual, expected)
	if len(results) > 0 {
		errors := "Actual does not equal expected:"
		for _, result := range results {
			errors += "\n - " + result
		}

		c.t.Errorf("\n%s\n\nActual:   %+v\nExpected: %+v", errors, c.actual, expected)
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
		actualOutput := buildActualErrorOutput(actual)
		expectedOutput := buildExpectedErrorOutput(expected)
		c.t.Errorf("\nGot:      %s\nExpected: %s", actualOutput, expectedOutput)
	}
}

// IsNotError ensures that the actual value is nil.
// It is analogous to IsError(nil).
func (c Chain) IsNotError() {
	c.t.Helper()
	c.IsError(nil)
}

func buildActualErrorOutput(actual error) string {
	actualErk, isActualErk := actual.(erk.Erkable) //nolint:errorlint // Want to output the top level error
	if !isActualErk {
		return fmt.Sprintf("%v", actual)
	}

	return fmt.Sprintf(
		"{KIND: \"%s\", MESSAGE: \"%s\", PARAMS: %+v}",
		erk.GetKindString(actualErk),
		actualErk.Error(),
		actualErk.Params(),
	)
}

func buildExpectedErrorOutput(expected error) string {
	expectedErk, isExpectedErk := expected.(erk.Erkable) //nolint:errorlint // Want to output the top level error
	if !isExpectedErk {
		return fmt.Sprintf("%v", expected)
	}

	return fmt.Sprintf(
		"{KIND: \"%s\", RAW MESSAGE: \"%s\", PARAMS: %+v}",
		erk.GetKindString(expectedErk),
		expectedErk.ExportRawMessage(),
		expectedErk.Params(),
	)
}

func isNil(value interface{}) bool {
	if value == nil {
		return true
	}

	reflection := reflect.ValueOf(value)
	isNilPointer := reflection.Kind() == reflect.Ptr && reflection.IsNil()
	return isNilPointer
}
