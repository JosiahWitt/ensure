package ensurepkg_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/JosiahWitt/ensure"
	"github.com/JosiahWitt/erk"
	"github.com/kr/pretty"
	"github.com/kr/text"
)

func TestChainIsTrue(t *testing.T) {
	t.Run("when true", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)
		mockT.EXPECT().Helper()

		ensure := ensure.New(mockT)
		ensure(true).IsTrue()
	})

	t.Run("when false", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)

		mockT.EXPECT().Fatalf("Got false, expected true").After(
			mockT.EXPECT().Helper(),
		)

		ensure := ensure.New(mockT)
		ensure(false).IsTrue()
	})

	t.Run("when not a boolean", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)

		const val = "not a boolean"
		mockT.EXPECT().Fatalf("Got type %T, expected boolean", val).After(
			mockT.EXPECT().Helper(),
		)

		ensure := ensure.New(mockT)
		ensure(val).IsTrue()
	})
}

func TestChainIsFalse(t *testing.T) {
	t.Run("when false", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)
		mockT.EXPECT().Helper()

		ensure := ensure.New(mockT)
		ensure(false).IsFalse()
	})

	t.Run("when true", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)

		mockT.EXPECT().Fatalf("Got true, expected false").After(
			mockT.EXPECT().Helper(),
		)

		ensure := ensure.New(mockT)
		ensure(true).IsFalse()
	})

	t.Run("when not a boolean", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)

		const val = "not a boolean"
		mockT.EXPECT().Fatalf("Got type %T, expected boolean", val).After(
			mockT.EXPECT().Helper(),
		)

		ensure := ensure.New(mockT)
		ensure(val).IsFalse()
	})
}

func TestChainIsNil(t *testing.T) {
	t.Run("when nil", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)
		mockT.EXPECT().Helper()

		ensure := ensure.New(mockT)
		ensure(nil).IsNil()
	})

	t.Run("when nil pointer", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)
		mockT.EXPECT().Helper()

		var nilPtr *string

		ensure := ensure.New(mockT)
		ensure(nilPtr).IsNil()
	})

	t.Run("when not nil", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)

		const val = "not nil"
		mockT.EXPECT().Fatalf("Got %+v, expected nil", val).After(
			mockT.EXPECT().Helper(),
		)

		ensure := ensure.New(mockT)
		ensure(val).IsNil()
	})
}

func TestChainEquals(t *testing.T) {
	const errorMessageFormat = "\n%s\n\nACTUAL:\n%s\n\nEXPECTED:\n%s"

	t.Run("when equal", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)
		mockT.EXPECT().Helper()

		ensure := ensure.New(mockT)
		ensure(ExamplePerson{Name: "John", Email: "john@test"}).Equals(ExamplePerson{Name: "John", Email: "john@test"})
	})

	t.Run("when unexported field is equal", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)
		mockT.EXPECT().Helper()

		ensure := ensure.New(mockT)
		ensure(ExamplePerson{Name: "John", Email: "john@test", ssn: "123456789"}).Equals(ExamplePerson{Name: "John", Email: "john@test", ssn: "123456789"})
	})

	t.Run("when both are nil", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)
		mockT.EXPECT().Helper()

		ensure := ensure.New(mockT)
		ensure(nil).Equals(nil)
	})

	t.Run("when strings are equal", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)
		mockT.EXPECT().Helper()

		ensure := ensure.New(mockT)
		ensure("abc").Equals("abc")
	})

	t.Run("when nil pointer equals nil", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)
		mockT.EXPECT().Helper()

		var nilPtr *string

		ensure := ensure.New(mockT)
		ensure(nilPtr).Equals(nil)
	})

	t.Run("when nil map equals empty map", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)
		mockT.EXPECT().Helper()

		var nilMap map[string]string

		ensure := ensure.New(mockT)
		ensure(nilMap).Equals(map[string]string{})
	})

	t.Run("when nil slice equals empty slice", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)
		mockT.EXPECT().Helper()

		var nilMap []string

		ensure := ensure.New(mockT)
		ensure(nilMap).Equals([]string{})
	})

	t.Run("when nil array equals empty array", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)
		mockT.EXPECT().Helper()

		var nilMap [2]string

		ensure := ensure.New(mockT)
		ensure(nilMap).Equals([2]string{})
	})

	t.Run("when one field is not equal", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)
		mockT.EXPECT().Fatalf(errorMessageFormat,
			"Actual does not equal expected:\n - Name: John != Sam",
			ExamplePerson{Name: "John", Email: "john@test"}.ExpectedOutput(),
			ExamplePerson{Name: "Sam", Email: "john@test"}.ExpectedOutput(),
		).After(
			mockT.EXPECT().Helper(),
		)

		ensure := ensure.New(mockT)
		ensure(ExamplePerson{Name: "John", Email: "john@test"}).Equals(ExamplePerson{Name: "Sam", Email: "john@test"})
	})

	t.Run("when not equal: expected is nil", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)
		mockT.EXPECT().Fatalf(errorMessageFormat,
			"Actual does not equal expected:\n - {John john@test  []} != <nil pointer>",
			ExamplePerson{Name: "John", Email: "john@test"}.ExpectedOutput(),
			"  nil",
		).After(
			mockT.EXPECT().Helper(),
		)

		ensure := ensure.New(mockT)
		ensure(ExamplePerson{Name: "John", Email: "john@test"}).Equals(nil)
	})

	t.Run("when not equal: actual is nil", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)
		mockT.EXPECT().Fatalf(errorMessageFormat,
			"Actual does not equal expected:\n - <nil pointer> != {John john@test  []}",
			"  nil",
			ExamplePerson{Name: "John", Email: "john@test"}.ExpectedOutput(),
		).After(
			mockT.EXPECT().Helper(),
		)

		ensure := ensure.New(mockT)
		ensure(nil).Equals(ExamplePerson{Name: "John", Email: "john@test"})
	})

	t.Run("when unexported field is not equal", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)
		mockT.EXPECT().Fatalf(errorMessageFormat,
			"Actual does not equal expected:\n - ssn: 123456789 != 123456780",
			ExamplePerson{Name: "John", Email: "john@test", ssn: "123456789"}.ExpectedOutput(),
			ExamplePerson{Name: "John", Email: "john@test", ssn: "123456780"}.ExpectedOutput(),
		).After(
			mockT.EXPECT().Helper(),
		)

		ensure := ensure.New(mockT)
		ensure(ExamplePerson{Name: "John", Email: "john@test", ssn: "123456789"}).Equals(ExamplePerson{Name: "John", Email: "john@test", ssn: "123456780"})
	})

	t.Run("when two fields are not equal", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)
		mockT.EXPECT().Fatalf(errorMessageFormat,
			"Actual does not equal expected:\n - Name: John != Sam\n - Messages.slice[1].Body: Hello != Greetings",
			ExamplePerson{
				Name:  "John",
				Email: "john@test",
				Messages: []ExampleMessage{
					{Body: "Hi"},
					{Body: "Hello"},
				},
			}.ExpectedOutput(),
			ExamplePerson{
				Name:  "Sam",
				Email: "john@test",
				Messages: []ExampleMessage{
					{Body: "Hi"},
					{Body: "Greetings"},
				},
			}.ExpectedOutput(),
		).After(
			mockT.EXPECT().Helper(),
		)

		ensure := ensure.New(mockT)
		ensure(ExamplePerson{
			Name:  "John",
			Email: "john@test",
			Messages: []ExampleMessage{
				{Body: "Hi"},
				{Body: "Hello"},
			},
		}).
			Equals(ExamplePerson{
				Name:  "Sam",
				Email: "john@test",
				Messages: []ExampleMessage{
					{Body: "Hi"},
					{Body: "Greetings"},
				},
			})
	})

	t.Run("when strings not equal: expected empty string", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)
		mockT.EXPECT().Fatalf(errorMessageFormat, "Actual does not equal expected:\n - abc != ",
			`  "abc"`,
			"  (empty string)",
		).After(
			mockT.EXPECT().Helper(),
		)

		ensure := ensure.New(mockT)
		ensure("abc").Equals("")
	})

	t.Run("when strings not equal: actual empty string", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)
		mockT.EXPECT().Fatalf(errorMessageFormat, "Actual does not equal expected:\n -  != abc",
			"  (empty string)",
			`  "abc"`,
		).After(
			mockT.EXPECT().Helper(),
		)

		ensure := ensure.New(mockT)
		ensure("").Equals("abc")
	})

	t.Run("when strings not equal: expected contains double quotes, newlines, and tabs", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)
		mockT.EXPECT().Fatalf(errorMessageFormat, "Actual does not equal expected:\n - abc\n\"xyz\"\n\tqwerty != abc",
			`  "abc\n\"xyz\"\n\tqwerty"`, // Formatted with quotes and escaped control characters
			`  "abc"`,
		).After(
			mockT.EXPECT().Helper(),
		)

		ensure := ensure.New(mockT)
		ensure("abc\n\"xyz\"\n\tqwerty").Equals("abc")
	})

	t.Run("when strings not equal: actual contains double quotes, newlines, and tabs", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)
		mockT.EXPECT().Fatalf(errorMessageFormat, "Actual does not equal expected:\n - abc != abc\n\"xyz\"\n\tqwerty",
			`  "abc"`,
			`  "abc\n\"xyz\"\n\tqwerty"`, // Formatted with quotes and escaped control characters
		).After(
			mockT.EXPECT().Helper(),
		)

		ensure := ensure.New(mockT)
		ensure("abc").Equals("abc\n\"xyz\"\n\tqwerty")
	})
}

func TestChainIsError(t *testing.T) {
	const errorFormat = "\nActual error is not the expected error:\n\tActual:   %s\n\tExpected: %s"

	t.Run("when equal error by reference", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)
		mockT.EXPECT().Helper()

		err := errors.New("my error")

		ensure := ensure.New(mockT)
		ensure(err).IsError(err)
	})

	t.Run("when equal error by Is method", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)
		mockT.EXPECT().Helper()

		const errMsg = "my error"

		ensure := ensure.New(mockT)
		ensure(TestError{Message: errMsg}).IsError(TestError{Message: errMsg})
	})

	t.Run("when both are nil", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)
		mockT.EXPECT().Helper()

		ensure := ensure.New(mockT)
		ensure(nil).IsError(nil)
	})

	t.Run("when not error type", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)

		err := errors.New("my error")
		const val = "not an error"
		mockT.EXPECT().Fatalf("Got type %T, expected error: \"%v\"", val, err).After(
			mockT.EXPECT().Helper(),
		)

		ensure := ensure.New(mockT)
		ensure(val).IsError(err)
	})

	t.Run("when not equal: two different errors by reference", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)

		err1 := errors.New("my error")
		err2 := errors.New("my error")

		mockT.EXPECT().Fatalf(errorFormat, err1.Error(), err2.Error()).After(
			mockT.EXPECT().Helper(),
		)

		ensure := ensure.New(mockT)
		ensure(err1).IsError(err2)
	})

	t.Run("when not equal: two different errors by Is method", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)

		err1 := TestError{Message: "error message 1"}
		err2 := TestError{Message: "error message 2"}
		mockT.EXPECT().Fatalf(errorFormat, err1.Error(), err2.Error()).After(
			mockT.EXPECT().Helper(),
		)

		ensure := ensure.New(mockT)
		ensure(err1).IsError(err2)
	})

	t.Run("when not equal: expected nil", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)

		err := errors.New("my error")
		mockT.EXPECT().Fatalf(errorFormat, err.Error(), "<nil>").After(
			mockT.EXPECT().Helper(),
		)

		ensure := ensure.New(mockT)
		ensure(err).IsError(nil)
	})

	t.Run("when not equal: got nil", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)

		err := errors.New("my error")
		mockT.EXPECT().Fatalf(errorFormat, "<nil>", err.Error()).After(
			mockT.EXPECT().Helper(),
		)

		ensure := ensure.New(mockT)
		ensure(nil).IsError(err)
	})

	t.Run("when not equal: erk errors: different kinds", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)

		type kind1 struct{ erk.DefaultKind }
		type kind2 struct{ erk.DefaultKind }

		expectedError := erk.New(kind1{}, "expected {{.a}}")
		actualError := erk.NewWith(kind2{}, "actual {{.a}}", erk.Params{"a": "hi"})
		mockT.EXPECT().Fatalf(
			errorFormat,
			fmt.Sprintf("{KIND: \"%s\", MESSAGE: \"actual hi\", PARAMS: map[a:hi]}", erk.GetKindString(actualError)),
			fmt.Sprintf("{KIND: \"%s\", RAW MESSAGE: \"expected {{.a}}\", PARAMS: map[]}", erk.GetKindString(expectedError)),
		).After(
			mockT.EXPECT().Helper(),
		)

		ensure := ensure.New(mockT)
		ensure(actualError).IsError(expectedError)
	})

	t.Run("when not equal: erk errors: same kind", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)

		type kind1 struct{ erk.DefaultKind }

		expectedError := erk.New(kind1{}, "expected {{.a}}")
		actualError := erk.NewWith(kind1{}, "actual {{.a}}", erk.Params{"a": "hi"})
		mockT.EXPECT().Fatalf(
			errorFormat,
			fmt.Sprintf("{KIND: \"%s\", MESSAGE: \"actual hi\", PARAMS: map[a:hi]}", erk.GetKindString(actualError)),
			fmt.Sprintf("{KIND: \"%s\", RAW MESSAGE: \"expected {{.a}}\", PARAMS: map[]}", erk.GetKindString(expectedError)),
		).After(
			mockT.EXPECT().Helper(),
		)

		ensure := ensure.New(mockT)
		ensure(actualError).IsError(expectedError)
	})

	t.Run("when not equal: erk errors: only expected is erk error", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)

		type kind1 struct{ erk.DefaultKind }

		expectedError := erk.New(kind1{}, "expected {{.a}}")
		actualError := errors.New("actual")
		mockT.EXPECT().Fatalf(
			errorFormat,
			actualError.Error(),
			fmt.Sprintf("{KIND: \"%s\", RAW MESSAGE: \"expected {{.a}}\", PARAMS: map[]}", erk.GetKindString(expectedError)),
		).After(
			mockT.EXPECT().Helper(),
		)

		ensure := ensure.New(mockT)
		ensure(actualError).IsError(expectedError)
	})

	t.Run("when not equal: erk errors: only actual is erk error", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)

		type kind1 struct{ erk.DefaultKind }

		expectedError := errors.New("expected")
		actualError := erk.NewWith(kind1{}, "actual {{.a}}", erk.Params{"a": "hi"})
		mockT.EXPECT().Fatalf(
			errorFormat,
			fmt.Sprintf("{KIND: \"%s\", MESSAGE: \"actual hi\", PARAMS: map[a:hi]}", erk.GetKindString(actualError)),
			expectedError.Error(),
		).After(
			mockT.EXPECT().Helper(),
		)

		ensure := ensure.New(mockT)
		ensure(actualError).IsError(expectedError)
	})
}

func TestChainIsNotError(t *testing.T) {
	t.Run("when no error", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)
		mockT.EXPECT().Helper().Times(2)

		ensure := ensure.New(mockT)
		ensure(nil).IsNotError()
	})

	t.Run("when error", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)

		err := errors.New("my error")
		mockT.EXPECT().Fatalf("\nActual error is not the expected error:\n\tActual:   %s\n\tExpected: %s", err.Error(), "<nil>").After(
			mockT.EXPECT().Helper().Times(2),
		)

		ensure := ensure.New(mockT)
		ensure(err).IsNotError()
	})
}

func TestChainIsEmpty(t *testing.T) {
	const notEmptyFormat = "Got %+v with length %d, expected it to be empty"

	t.Run("when empty: array", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)
		mockT.EXPECT().Helper()

		ensure := ensure.New(mockT)
		ensure([0]string{}).IsEmpty()
	})

	t.Run("when not empty: array", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)

		mockT.EXPECT().Fatalf(notEmptyFormat, [2]string{"1", "2"}, 2).After(
			mockT.EXPECT().Helper(),
		)

		ensure := ensure.New(mockT)
		ensure([2]string{"1", "2"}).IsEmpty()
	})

	t.Run("when empty: slice", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)
		mockT.EXPECT().Helper()

		ensure := ensure.New(mockT)
		ensure([]string{}).IsEmpty()
	})

	t.Run("when not empty: slice", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)

		mockT.EXPECT().Fatalf(notEmptyFormat, []string{"1"}, 1).After(
			mockT.EXPECT().Helper(),
		)

		ensure := ensure.New(mockT)
		ensure([]string{"1"}).IsEmpty()
	})

	t.Run("when empty: string", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)
		mockT.EXPECT().Helper()

		ensure := ensure.New(mockT)
		ensure("").IsEmpty()
	})

	t.Run("when not empty: string", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)

		mockT.EXPECT().Fatalf(notEmptyFormat, "not empty", 9).After(
			mockT.EXPECT().Helper(),
		)

		ensure := ensure.New(mockT)
		ensure("not empty").IsEmpty()
	})

	t.Run("when empty: map", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)
		mockT.EXPECT().Helper()

		ensure := ensure.New(mockT)
		ensure(map[string]string{}).IsEmpty()
	})

	t.Run("when not empty: map", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)

		mockT.EXPECT().Fatalf(notEmptyFormat, map[string]string{"hello": "world"}, 1).After(
			mockT.EXPECT().Helper(),
		)

		ensure := ensure.New(mockT)
		ensure(map[string]string{"hello": "world"}).IsEmpty()
	})

	t.Run("when not valid type", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)

		mockT.EXPECT().Fatalf("Got type %T, expected array, slice, string, or map", 1234).After(
			mockT.EXPECT().Helper(),
		)

		ensure := ensure.New(mockT)
		ensure(1234).IsEmpty()
	})
}

type TestError struct {
	Message string
}

func (t TestError) Is(err error) bool {
	inputErr := TestError{}
	if errors.As(err, &inputErr) {
		return inputErr.Message == t.Message
	}

	return false
}

func (t TestError) Error() string {
	return t.Message
}

type ExampleMessage struct {
	Body string
}

type ExamplePerson struct {
	Name  string
	Email string
	ssn   string

	Messages []ExampleMessage
}

func (p ExamplePerson) ExpectedOutput() string {
	return text.Indent(pretty.Sprint(p), "  ")
}
