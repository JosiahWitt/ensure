package integration_test_suite_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/JosiahWitt/ensure"
	"github.com/JosiahWitt/ensure/internal/stringerr"
)

// TODO: These tests currently verify happy paths, since failing directly in these tests
// isn't possible without failing the whole test suite. Consider calling test fixtures
// and expecting them to fail.

func TestContains(t *testing.T) {
	t.Run("assertion succeeds when actual contains the expected", func(t *testing.T) {
		ensure := ensure.New(t)

		ensure("hello").Contains("el")
		ensure([]string{"apple", "orange", "banana"}).Contains("orange")
	})
}

func TestDoesNotContain(t *testing.T) {
	t.Run("assertion succeeds when actual does not contain the expected", func(t *testing.T) {
		ensure := ensure.New(t)

		ensure("hello").DoesNotContain("world")
		ensure([]string{"apple", "orange", "banana"}).DoesNotContain("pl")
	})
}

func TestEqual(t *testing.T) {
	t.Run("assertion succeeds when actual equals the expected", func(t *testing.T) {
		ensure := ensure.New(t)

		type NestedStruct struct {
			Value string
		}

		type OuterStruct struct {
			Value  string
			Nested *NestedStruct
		}

		ensure(nil).Equals(nil)
		ensure(error(nil)).Equals(nil)
		ensure("a").Equals("a")
		ensure(123).Equals(123)
		ensure(&OuterStruct{Value: "a", Nested: &NestedStruct{Value: "b"}}).Equals(&OuterStruct{Value: "a", Nested: &NestedStruct{Value: "b"}})
	})
}

func TestIsEmpty(t *testing.T) {
	t.Run("assertion succeeds when actual is empty", func(t *testing.T) {
		ensure := ensure.New(t)

		ensure("").IsEmpty()
		ensure([]string{}).IsEmpty()
		ensure(map[string]string{}).IsEmpty()
	})
}

func TestIsError(t *testing.T) {
	t.Run("assertion succeeds when actual is the expected", func(t *testing.T) {
		ensure := ensure.New(t)

		sampleErr := errors.New("bad day")

		ensure(nil).IsError(nil)
		ensure(error(nil)).IsError(error(nil))
		ensure(sampleErr).IsError(sampleErr)
		ensure(stringerr.Newf("oops")).IsError(stringerr.Newf("oops"))
		ensure(fmt.Errorf("wrapped %w", stringerr.Newf("oops"))).IsError(stringerr.Newf("oops"))
	})
}

func TestIsFalse(t *testing.T) {
	t.Run("assertion succeeds when actual is false", func(t *testing.T) {
		ensure := ensure.New(t)

		ensure(false).IsFalse()
	})
}

func TestIsNil(t *testing.T) {
	t.Run("assertion succeeds when actual is nil", func(t *testing.T) {
		ensure := ensure.New(t)

		ensure(nil).IsNil()
		ensure(error(nil)).IsNil()
	})
}

func TestIsNotEmpty(t *testing.T) {
	t.Run("assertion succeeds when actual is not empty", func(t *testing.T) {
		ensure := ensure.New(t)

		ensure("abc").IsNotEmpty()
		ensure([]string{"hello"}).IsNotEmpty()
		ensure(map[string]string{"a": "pair"}).IsNotEmpty()
	})
}

func TestIsNotError(t *testing.T) {
	t.Run("assertion succeeds when actual is nil", func(t *testing.T) {
		ensure := ensure.New(t)

		ensure(error(nil)).IsNotError()
		ensure(nil).IsNotError()
	})
}

func TestIsNotNil(t *testing.T) {
	t.Run("assertion succeeds when actual is not nil", func(t *testing.T) {
		ensure := ensure.New(t)

		val := "a"

		ensure("a").IsNotNil()
		ensure(&val).IsNotNil()
	})
}

func TestIsTrue(t *testing.T) {
	t.Run("assertion succeeds when actual is true", func(t *testing.T) {
		ensure := ensure.New(t)

		ensure(true).IsTrue()
	})
}

func TestMatchesAllErrors(t *testing.T) {
	t.Run("assertion succeeds when action matches all expected errors", func(t *testing.T) {
		ensure := ensure.New(t)

		sampleErr := errors.New("bad day")
		sampleWrappedErr := fmt.Errorf("wrapped %w", stringerr.Newf("oops"))

		ensure(nil).MatchesAllErrors()
		ensure(error(nil)).MatchesAllErrors()
		ensure(sampleErr).MatchesAllErrors(sampleErr)
		ensure(stringerr.Newf("oops")).MatchesAllErrors(stringerr.Newf("oops"))
		ensure(sampleWrappedErr).MatchesAllErrors(sampleWrappedErr, stringerr.Newf("oops"))
	})
}

func TestMatchesRegexp(t *testing.T) {
	t.Run("assertion succeeds when actual matches the expected regular expression", func(t *testing.T) {
		ensure := ensure.New(t)

		ensure("hello<world>!").MatchesRegexp(".*<.*>!")
	})
}
