package ensurepkg

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/JosiahWitt/erk"
	"github.com/go-test/deep"
	"github.com/kr/pretty"
	"github.com/kr/text"
)

// IsTrue ensures the actual value is the boolean "true".
func (c *Chain) IsTrue() {
	c.t.Helper()
	c.markRun()

	actual, ok := c.actual.(bool)
	if !ok {
		c.t.Fatalf("Got type %T, expected boolean", c.actual)
		return
	}

	if !actual {
		c.t.Fatalf("Got false, expected true")
	}
}

// IsFalse ensures the actual value is the boolean "false".
func (c *Chain) IsFalse() {
	c.t.Helper()
	c.markRun()

	actual, ok := c.actual.(bool)
	if !ok {
		c.t.Fatalf("Got type %T, expected boolean", c.actual)
		return
	}

	if actual {
		c.t.Fatalf("Got true, expected false")
	}
}

// IsNil ensures the actual value is nil or a nil pointer.
func (c *Chain) IsNil() {
	c.t.Helper()
	c.markRun()

	if !isNil(c.actual) {
		c.t.Fatalf("Got %+v, expected nil", c.actual)
	}
}

// IsNotNil ensures the actual value is not nil and not a nil pointer.
func (c *Chain) IsNotNil() {
	c.t.Helper()
	c.markRun()

	if isNil(c.actual) {
		c.t.Fatalf("Got nil of type %T, expected it not to be nil", c.actual)
	}
}

// Equals ensures the actual value equals the expected value.
// Equals uses deep.Equal to print easy to read diffs.
func (c *Chain) Equals(expected interface{}) {
	c.t.Helper()
	c.markRun()

	// If we expect nil, return early if actual is nil or if it is a nil pointer
	if expected == nil && isNil(c.actual) {
		return
	}

	deep.CompareUnexportedFields = true
	deep.NilMapsAreEmpty = false
	deep.NilSlicesAreEmpty = false
	results := deep.Equal(c.actual, expected)
	if len(results) > 0 {
		errors := "Actual does not equal expected:"
		for _, result := range results {
			errors += "\n - " + result
		}

		c.t.Fatalf("\n%s\n\nACTUAL:\n%s\n\nEXPECTED:\n%s", errors, prettyFormat(c.actual), prettyFormat(expected))
	}
}

// IsError ensures the actual value equals the expected error.
// IsError uses errors.Is to support Go 1.13+ error comparisons.
func (c *Chain) IsError(expected error) {
	c.t.Helper()
	c.markRun()

	if c.actual == nil && expected == nil {
		return
	}

	actual, ok := c.actual.(error)
	if !ok && c.actual != nil {
		c.t.Fatalf("Got type %T, expected error: \"%v\"", c.actual, expected)
		return
	}

	if !errors.Is(actual, expected) {
		actualOutput := buildActualErrorOutput(actual)
		expectedOutput := buildExpectedErrorOutput(expected)
		c.t.Fatalf("\nActual error is not the expected error:\n\tActual:   %s\n\tExpected: %s", actualOutput, expectedOutput)
	}
}

// IsNotError ensures that the actual value is nil.
// It is analogous to IsError(nil).
func (c *Chain) IsNotError() {
	c.t.Helper()
	c.IsError(nil)
}

// IsEmpty ensures that the actual value is empty.
// It only supports arrays, slices, strings, or maps.
func (c *Chain) IsEmpty() {
	c.t.Helper()
	c.markRun()

	length, err := lengthOf(c.actual)
	if err != nil {
		c.t.Fatalf(err.Error())
		return
	}

	if length > 0 {
		c.t.Fatalf("Got %+v with length %d, expected it to be empty", c.actual, length)
	}
}

// IsNotEmpty ensures that the actual value is not empty.
// It only supports arrays, slices, strings, or maps.
func (c *Chain) IsNotEmpty() {
	c.t.Helper()
	c.markRun()

	length, err := lengthOf(c.actual)
	if err != nil {
		c.t.Fatalf(err.Error())
		return
	}

	if length == 0 {
		c.t.Fatalf("Got %+v, expected it to not be empty", c.actual)
	}
}

// Contains ensures that the actual value contains the expected value.
// It only supports searching strings, arrays, or slices for the expected value.
// If both the actual and expected are strings, strings.Contains(...) is used.
//
// For example:
//  ensure("abc").Contains("b") // Succeeds
//  ensure("abc").Contains("z") // Fails
//
//  ensure([]string{"abc", "xyz"}).Contains("xyz") // Succeeds
//  ensure([]string{"abc", "xyz"}).Contains("y") // Fails
func (c *Chain) Contains(expected interface{}) {
	c.t.Helper()
	c.markRun()

	doesContain, err := contains(c.actual, expected)
	if err != nil {
		c.t.Fatalf(err.Error())
		return
	}

	if !doesContain {
		c.t.Fatalf(
			"Actual does not contain expected:\n\nACTUAL:\n%s\n\nEXPECTED TO CONTAIN:\n%s",
			prettyFormat(c.actual),
			prettyFormat(expected),
		)
	}
}

// DoesNotContain ensures that the actual value does not contain the expected value.
// It only supports verifying that strings, arrays, or slices do not contain the expected value.
// If both the actual and expected are strings, strings.Contains(...) is used.
//
// For example:
//  ensure("abc").DoesNotContain("b") // Fails
//  ensure("abc").DoesNotContain("z") // Succeeds
//
//  ensure([]string{"abc", "xyz"}).DoesNotContain("xyz") // Fails
//  ensure([]string{"abc", "xyz"}).DoesNotContain("y") // Succeeds
func (c *Chain) DoesNotContain(expected interface{}) {
	c.t.Helper()
	c.markRun()

	doesContain, err := contains(c.actual, expected)
	if err != nil {
		c.t.Fatalf(err.Error())
		return
	}

	if doesContain {
		c.t.Fatalf(
			"Actual contains expected, but did not expect it to:\n\nACTUAL:\n%s\n\nEXPECTED NOT TO CONTAIN:\n%s",
			prettyFormat(c.actual),
			prettyFormat(expected),
		)
	}
}

// MatchesRegexp ensures that the actual value matches the regular expression pattern provided.
// It only supports strings as actual values.
func (c *Chain) MatchesRegexp(pattern string) {
	c.t.Helper()
	c.markRun()

	if pattern == "" {
		c.t.Fatalf("Cannot match against an empty pattern")
		return
	}

	actualStr, ok := c.actual.(string)
	if !ok {
		c.t.Fatalf("Actual is not a string, it's a %T", c.actual)
		return
	}

	patternRegexp, err := regexp.Compile(pattern)
	if err != nil {
		c.t.Fatalf("Unable to compile regular expression: %s\nERROR: %v", pattern, err)
		return
	}

	isMatch := patternRegexp.MatchString(actualStr)
	if !isMatch {
		c.t.Fatalf(
			"Actual does not match regular expression:\n\nACTUAL:\n%s\n\nEXPECTED TO MATCH:\n%s",
			prettyFormat(c.actual),
			prettyFormat(pattern),
		)
	}
}

func (c *Chain) markRun() {
	c.wasRun = true
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

func lengthOf(value interface{}) (int, error) {
	reflectValue := reflect.ValueOf(value)
	reflectKind := reflectValue.Kind()
	if reflectKind != reflect.Array && reflectKind != reflect.Slice && reflectKind != reflect.String && reflectKind != reflect.Map {
		return 0, fmt.Errorf("Got type %T, expected array, slice, string, or map", value) //nolint:goerr113,stylecheck // Only used internally
	}

	return reflectValue.Len(), nil
}

func contains(items, value interface{}) (bool, error) {
	if str, strOk := items.(string); strOk {
		substr, substrOk := value.(string)
		if !substrOk {
			return false, fmt.Errorf("Got string, but expected is a %T, and a string can only contain other strings", value) //nolint:goerr113,stylecheck,lll // Only used internally
		}

		return strings.Contains(str, substr), nil
	}

	itemsReflectValue := reflect.ValueOf(items)
	itemsReflectKind := itemsReflectValue.Kind()
	if itemsReflectKind != reflect.Array && itemsReflectKind != reflect.Slice {
		return false, fmt.Errorf("Got type %T, expected string, array, or slice", value) //nolint:goerr113,stylecheck // Only used internally
	}

	for i := 0; i < itemsReflectValue.Len(); i++ {
		item := itemsReflectValue.Index(i)
		if reflect.DeepEqual(item.Interface(), value) {
			return true, nil
		}
	}

	return false, nil
}

func prettyFormat(value interface{}) string {
	return text.Indent(prettyFormatValue(value), "  ")
}

func prettyFormatValue(value interface{}) string {
	if str, ok := value.(string); ok {
		return prettyFormatString(str)
	}

	return pretty.Sprint(value)
}

func prettyFormatString(str string) string {
	if str == "" {
		return "(empty string)"
	}

	return strconv.Quote(str)
}
