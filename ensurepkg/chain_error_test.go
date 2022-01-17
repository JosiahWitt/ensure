package ensurepkg_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/JosiahWitt/ensure"
	"github.com/JosiahWitt/ensure/ensurepkg"
	"github.com/JosiahWitt/ensure/internal/mocks/github.com/JosiahWitt/ensure/mock_ensurepkg"
	"github.com/JosiahWitt/erk"
)

func TestChainIsError(t *testing.T) {
	t.Run("when actual is not error type", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)

		err := errors.New("my error")
		const val = "not an error"
		mockT.EXPECT().Fatalf("Got type %T, expected error: \"%v\"", val, err).After(
			mockT.EXPECT().Helper(),
		)

		ensure := ensure.New(mockT)
		ensure(val).IsError(err)
	})

	sharedIsErrorTests(t, func(mockT *mock_ensurepkg.MockT, chain *ensurepkg.Chain, expected error) {
		chain.IsError(expected)
	})
}

func TestChainMatchesAllErrors(t *testing.T) {
	t.Run("when actual is not error type", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)

		const val = "not an error"
		mockT.EXPECT().Fatalf("Got type %T, expected an error", val).After(
			mockT.EXPECT().Helper(),
		)

		ensure := ensure.New(mockT)
		ensure(val).MatchesAllErrors(errors.New("something"), errors.New("else"))
	})

	t.Run("when no expected errors", func(t *testing.T) {
		t.Run("when got error and expected empty slice", func(t *testing.T) {
			mockT := setupMockTWithCleanupCheck(t)

			mockT.EXPECT().Fatalf("\nExpected no error, but got: %s", "hi").After(
				mockT.EXPECT().Helper(),
			)

			ensure := ensure.New(mockT)
			ensure(errors.New("hi")).MatchesAllErrors([]error{}...)
		})

		t.Run("when no error and expected empty slice", func(t *testing.T) {
			mockT := setupMockTWithCleanupCheck(t)
			mockT.EXPECT().Helper()

			ensure := ensure.New(mockT)
			ensure(nil).MatchesAllErrors([]error{}...)
		})

		t.Run("when got error and expected nil slice", func(t *testing.T) {
			mockT := setupMockTWithCleanupCheck(t)

			mockT.EXPECT().Fatalf("\nExpected no error, but got: %s", "hi").After(
				mockT.EXPECT().Helper(),
			)

			var errs []error // nil error slice
			ensure := ensure.New(mockT)
			ensure(errors.New("hi")).MatchesAllErrors(errs...)
		})

		t.Run("when no error and expected nil slice", func(t *testing.T) {
			mockT := setupMockTWithCleanupCheck(t)
			mockT.EXPECT().Helper()

			var errs []error // nil error slice
			ensure := ensure.New(mockT)
			ensure(nil).MatchesAllErrors(errs...)
		})
	})

	t.Run("when one expected error", func(t *testing.T) {
		sharedIsErrorTests(t, func(mockT *mock_ensurepkg.MockT, chain *ensurepkg.Chain, expected error) {
			mockT.EXPECT().Helper()

			chain.MatchesAllErrors(expected)
		})
	})

	t.Run("when two expected errors", func(t *testing.T) {
		const errorFormat = "\nActual error is not all of the expected errors:\n\tActual:\n\t     %s\n\n\tExpected all of:%s"

		t.Run("when equal errors by reference", func(t *testing.T) {
			mockT := setupMockTWithCleanupCheck(t)
			mockT.EXPECT().Helper()

			err := errors.New("my error")

			ensure := ensure.New(mockT)
			ensure(err).MatchesAllErrors(err, err)
		})

		t.Run("when equal errors by Is method", func(t *testing.T) {
			mockT := setupMockTWithCleanupCheck(t)
			mockT.EXPECT().Helper()

			const errMsg = "my error"

			ensure := ensure.New(mockT)
			ensure(TestError{Unique: 1, Message: errMsg}).MatchesAllErrors(
				TestError{Unique: 2, Message: errMsg},
				TestError{Unique: 3, Message: errMsg},
			)
		})

		t.Run("when all are nil", func(t *testing.T) {
			mockT := setupMockTWithCleanupCheck(t)
			mockT.EXPECT().Helper()

			ensure := ensure.New(mockT)
			ensure(nil).MatchesAllErrors(nil, nil)
		})

		t.Run("when not equal: different errors by reference", func(t *testing.T) {
			mockT := setupMockTWithCleanupCheck(t)

			err1 := errors.New("my error")
			err2 := errors.New("my error")
			err3 := errors.New("my error")

			mockT.EXPECT().Fatalf(errorFormat, err1.Error(), "\n\t  ❌ my error\n\t  ❌ my error").After(
				mockT.EXPECT().Helper(),
			)

			ensure := ensure.New(mockT)
			ensure(err1).MatchesAllErrors(err2, err3)
		})

		t.Run("when some not equal: different errors by reference", func(t *testing.T) {
			mockT := setupMockTWithCleanupCheck(t)

			err1 := errors.New("my error")
			err2 := errors.New("my error")

			mockT.EXPECT().Fatalf(errorFormat, err1.Error(), "\n\t  ✅ my error\n\t  ❌ my error").After(
				mockT.EXPECT().Helper(),
			)

			ensure := ensure.New(mockT)
			ensure(err1).MatchesAllErrors(err1, err2)
		})

		t.Run("when not equal: different errors by Is method", func(t *testing.T) {
			mockT := setupMockTWithCleanupCheck(t)

			err1 := TestError{Unique: 1, Message: "error message 1"}
			err2 := TestError{Unique: 2, Message: "error message 2"}
			err3 := TestError{Unique: 3, Message: "error message 3"}

			mockT.EXPECT().Fatalf(errorFormat, err1.Error(), "\n\t  ❌ error message 2\n\t  ❌ error message 3").After(
				mockT.EXPECT().Helper(),
			)

			ensure := ensure.New(mockT)
			ensure(err1).MatchesAllErrors(err2, err3)
		})

		t.Run("when some not equal: different errors by Is method", func(t *testing.T) {
			mockT := setupMockTWithCleanupCheck(t)

			err1 := TestError{Unique: 1, Message: "error message 1"}
			err2 := TestError{Unique: 2, Message: "error message 2"}

			mockT.EXPECT().Fatalf(errorFormat, err1.Error(), "\n\t  ❌ error message 2\n\t  ✅ error message 1").After(
				mockT.EXPECT().Helper(),
			)

			ensure := ensure.New(mockT)
			ensure(err1).MatchesAllErrors(err2, err1)
		})

		t.Run("when not equal: actual is nil", func(t *testing.T) {
			mockT := setupMockTWithCleanupCheck(t)

			err1 := errors.New("my error 1")
			err2 := errors.New("my error 2")
			mockT.EXPECT().Fatalf(errorFormat, "<nil>", "\n\t  ❌ my error 1\n\t  ❌ my error 2").After(
				mockT.EXPECT().Helper(),
			)

			ensure := ensure.New(mockT)
			ensure(nil).MatchesAllErrors(err1, err2)
		})

		t.Run("when some not equal: actual is nil", func(t *testing.T) {
			mockT := setupMockTWithCleanupCheck(t)

			err := errors.New("my error")
			mockT.EXPECT().Fatalf(errorFormat, "<nil>", "\n\t  ❌ my error\n\t  ✅ <nil>").After(
				mockT.EXPECT().Helper(),
			)

			ensure := ensure.New(mockT)
			ensure(nil).MatchesAllErrors(err, nil)
		})

		t.Run("when not equal: erk errors: different kinds", func(t *testing.T) {
			mockT := setupMockTWithCleanupCheck(t)

			type kind1 struct{ erk.DefaultKind }
			type kind2 struct{ erk.DefaultKind }
			type kind3 struct{ erk.DefaultKind }

			expectedError1 := erk.New(kind1{}, "expected 1 {{.a}}")
			expectedError2 := erk.New(kind2{}, "expected 2 {{.a}}")
			actualError := erk.NewWith(kind3{}, "actual {{.a}}", erk.Params{"a": "hi"})
			mockT.EXPECT().Fatalf(
				errorFormat,
				fmt.Sprintf("{KIND: \"%s\", MESSAGE: \"actual hi\", PARAMS: map[a:hi]}", erk.GetKindString(actualError)),

				fmt.Sprintf("\n\t  ❌ %s\n\t  ❌ %s",
					fmt.Sprintf("{KIND: \"%s\", RAW MESSAGE: \"expected 1 {{.a}}\", PARAMS: map[]}", erk.GetKindString(expectedError1)),
					fmt.Sprintf("{KIND: \"%s\", RAW MESSAGE: \"expected 2 {{.a}}\", PARAMS: map[]}", erk.GetKindString(expectedError2)),
				),
			).After(
				mockT.EXPECT().Helper(),
			)

			ensure := ensure.New(mockT)
			ensure(actualError).MatchesAllErrors(expectedError1, expectedError2)
		})

		t.Run("when some not equal: erk errors: different kinds", func(t *testing.T) {
			mockT := setupMockTWithCleanupCheck(t)

			type kind1 struct{ erk.DefaultKind }
			type kind2 struct{ erk.DefaultKind }

			expectedError1 := erk.New(kind1{}, "expected {{.a}}")
			err := erk.New(kind2{}, "actual {{.a}}")
			actualError := erk.WithParams(err, erk.Params{"a": "hi"})
			mockT.EXPECT().Fatalf(
				errorFormat,
				fmt.Sprintf("{KIND: \"%s\", MESSAGE: \"actual hi\", PARAMS: map[a:hi]}", erk.GetKindString(actualError)),

				fmt.Sprintf("\n\t  ✅ %s\n\t  ❌ %s",
					fmt.Sprintf("{KIND: \"%s\", RAW MESSAGE: \"actual {{.a}}\", PARAMS: map[]}", erk.GetKindString(err)),
					fmt.Sprintf("{KIND: \"%s\", RAW MESSAGE: \"expected {{.a}}\", PARAMS: map[]}", erk.GetKindString(expectedError1)),
				),
			).After(
				mockT.EXPECT().Helper(),
			)

			ensure := ensure.New(mockT)
			ensure(actualError).MatchesAllErrors(err, expectedError1)
		})

		t.Run("when not equal: erk errors: only one expected is erk error", func(t *testing.T) {
			mockT := setupMockTWithCleanupCheck(t)

			type kind1 struct{ erk.DefaultKind }

			expectedError1 := erk.New(kind1{}, "expected 1 {{.a}}")
			expectedError2 := errors.New("expected 2")
			actualError := errors.New("actual")
			mockT.EXPECT().Fatalf(
				errorFormat,
				actualError.Error(),

				fmt.Sprintf("\n\t  ❌ %s\n\t  ❌ %s",
					fmt.Sprintf("{KIND: \"%s\", RAW MESSAGE: \"expected 1 {{.a}}\", PARAMS: map[]}", erk.GetKindString(expectedError1)),
					"expected 2",
				),
			).After(
				mockT.EXPECT().Helper(),
			)

			ensure := ensure.New(mockT)
			ensure(actualError).MatchesAllErrors(expectedError1, expectedError2)
		})

		t.Run("when some not equal: erk errors: only one expected is erk error", func(t *testing.T) {
			mockT := setupMockTWithCleanupCheck(t)

			type kind1 struct{ erk.DefaultKind }

			actualError := errors.New("actual")
			expectedError1 := erk.New(kind1{}, "expected 1 {{.a}}")
			expectedError2 := actualError
			mockT.EXPECT().Fatalf(
				errorFormat,
				actualError.Error(),

				fmt.Sprintf("\n\t  ❌ %s\n\t  ✅ %s",
					fmt.Sprintf("{KIND: \"%s\", RAW MESSAGE: \"expected 1 {{.a}}\", PARAMS: map[]}", erk.GetKindString(expectedError1)),
					"actual",
				),
			).After(
				mockT.EXPECT().Helper(),
			)

			ensure := ensure.New(mockT)
			ensure(actualError).MatchesAllErrors(expectedError1, expectedError2)
		})

		t.Run("when not equal: erk errors: only actual is erk error", func(t *testing.T) {
			mockT := setupMockTWithCleanupCheck(t)

			type kind1 struct{ erk.DefaultKind }

			expectedError1 := errors.New("expected 1")
			expectedError2 := errors.New("expected 2")
			actualError := erk.NewWith(kind1{}, "actual {{.a}}", erk.Params{"a": "hi"})
			mockT.EXPECT().Fatalf(
				errorFormat,
				fmt.Sprintf("{KIND: \"%s\", MESSAGE: \"actual hi\", PARAMS: map[a:hi]}", erk.GetKindString(actualError)),

				fmt.Sprintf("\n\t  ❌ %s\n\t  ❌ %s",
					expectedError1.Error(),
					expectedError2.Error(),
				),
			).After(
				mockT.EXPECT().Helper(),
			)

			ensure := ensure.New(mockT)
			ensure(actualError).MatchesAllErrors(expectedError1, expectedError2)
		})
	})
}

func sharedIsErrorTests(t *testing.T, run func(mockT *mock_ensurepkg.MockT, chain *ensurepkg.Chain, expected error)) {
	const errorFormat = "\nActual error is not the expected error:\n\tActual:   %s\n\tExpected: %s"

	t.Run("when equal error by reference", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)
		mockT.EXPECT().Helper()

		err := errors.New("my error")

		ensure := ensure.New(mockT)
		run(mockT, ensure(err), err)
	})

	t.Run("when equal error by Is method", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)
		mockT.EXPECT().Helper()

		const errMsg = "my error"

		ensure := ensure.New(mockT)
		run(mockT, ensure(TestError{Unique: 1, Message: errMsg}), TestError{Unique: 2, Message: errMsg})
	})

	t.Run("when both are nil", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)
		mockT.EXPECT().Helper()

		ensure := ensure.New(mockT)
		run(mockT, ensure(nil), nil)
	})

	t.Run("when not equal: two different errors by reference", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)

		err1 := errors.New("my error")
		err2 := errors.New("my error")

		mockT.EXPECT().Fatalf(errorFormat, err1.Error(), err2.Error()).After(
			mockT.EXPECT().Helper(),
		)

		ensure := ensure.New(mockT)
		run(mockT, ensure(err1), err2)
	})

	t.Run("when not equal: two different errors by Is method", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)

		err1 := TestError{Unique: 1, Message: "error message 1"}
		err2 := TestError{Unique: 2, Message: "error message 2"}
		mockT.EXPECT().Fatalf(errorFormat, err1.Error(), err2.Error()).After(
			mockT.EXPECT().Helper(),
		)

		ensure := ensure.New(mockT)
		run(mockT, ensure(err1), err2)
	})

	t.Run("when not equal: expected nil", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)

		err := errors.New("my error")
		mockT.EXPECT().Fatalf(errorFormat, err.Error(), "<nil>").After(
			mockT.EXPECT().Helper(),
		)

		ensure := ensure.New(mockT)
		run(mockT, ensure(err), nil)
	})

	t.Run("when not equal: got nil", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)

		err := errors.New("my error")
		mockT.EXPECT().Fatalf(errorFormat, "<nil>", err.Error()).After(
			mockT.EXPECT().Helper(),
		)

		ensure := ensure.New(mockT)
		run(mockT, ensure(nil), err)
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
		run(mockT, ensure(actualError), expectedError)
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
		run(mockT, ensure(actualError), expectedError)
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
		run(mockT, ensure(actualError), expectedError)
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
		run(mockT, ensure(actualError), expectedError)
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

type TestError struct {
	Message string
	Unique  int
}

func (t TestError) Is(err error) bool {
	inputErr := TestError{Unique: 1}
	if errors.As(err, &inputErr) {
		return inputErr.Message == t.Message
	}

	return false
}

func (t TestError) Error() string {
	return t.Message
}
