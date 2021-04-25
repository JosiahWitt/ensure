package ensurepkg

import (
	"errors"
	"fmt"

	"github.com/JosiahWitt/erk"
)

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

// MatchesAllErrors ensures the actual value is all of the expected errors.
// This is useful for validating various levels of a wrapped error.
//
// If no errors are provided, this method assumes no error is expected.
//
// MatchesAllErrors uses errors.Is with each expected error to support Go 1.13+ error comparisons.
func (c *Chain) MatchesAllErrors(expectedErrors ...error) {
	c.t.Helper()
	c.markRun()

	actual, ok := c.actual.(error)
	if !ok && c.actual != nil {
		c.t.Fatalf("Got type %T, expected an error", c.actual)
		return
	}

	if len(expectedErrors) == 0 {
		if c.actual != nil {
			c.t.Fatalf("\nExpected no error, but got: %s", buildActualErrorOutput(actual))
		}

		return
	}

	if len(expectedErrors) == 1 {
		c.IsError(expectedErrors[0])
		return
	}

	failed := false
	failureDetails := ""
	for _, expected := range expectedErrors {
		status := "✅"
		if !errors.Is(actual, expected) {
			failed = true
			status = "❌"
		}

		failureDetails += fmt.Sprintf("\n\t  %s %s", status, buildExpectedErrorOutput(expected))
	}

	if failed {
		actualOutput := buildActualErrorOutput(actual)
		c.t.Fatalf("\nActual error is not all of the expected errors:\n\tActual:\n\t     %s\n\n\tExpected all of:%s",
			actualOutput,
			failureDetails,
		)
	}
}

// IsNotError ensures that the actual value is nil.
// It is analogous to IsError(nil).
func (c *Chain) IsNotError() {
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
