package ensurepkg

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/go-test/deep"
	"github.com/kr/pretty"
	"github.com/kr/text"
)

// Mutex to synchronize accessing deep.
//nolint:gochecknoglobals // Deep global variables need a global mutex.
var deepGlobalMu sync.Mutex

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

	results := checkEquality(c.actual, expected)
	if len(results) > 0 {
		errors := "Actual does not equal expected:"
		for _, result := range results {
			errors += "\n - " + result
		}

		c.t.Fatalf("\n%s\n\nACTUAL:\n%s\n\nEXPECTED:\n%s", errors, prettyFormat(c.actual), prettyFormat(expected))
	}
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

func checkEquality(actual, expected interface{}) []string {
	// Since deep only supports global settings, we wrap setting them
	// and the equality check in a mutex for concurrency safety.

	deepGlobalMu.Lock()
	defer deepGlobalMu.Unlock()

	deep.CompareUnexportedFields = true
	deep.NilMapsAreEmpty = false
	deep.NilSlicesAreEmpty = false

	return deep.Equal(actual, expected)
}
