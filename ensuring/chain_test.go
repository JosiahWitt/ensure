package ensuring_test

import (
	"sync"
	"testing"

	"github.com/JosiahWitt/ensure"
	"github.com/golang/mock/gomock"
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
	const failureFormat = "Got %+v, expected nil"

	t.Run("when nil", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)
		mockT.EXPECT().Helper()

		ensure := ensure.New(mockT)
		ensure(nil).IsNil()
	})

	t.Run("not nilable", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)

		const val = "not nil"
		mockT.EXPECT().Fatalf(failureFormat, val).After(
			mockT.EXPECT().Helper(),
		)

		ensure := ensure.New(mockT)
		ensure(val).IsNil()
	})

	t.Run("pointer", func(t *testing.T) {
		t.Run("when nil", func(t *testing.T) {
			mockT := setupMockTWithCleanupCheck(t)
			mockT.EXPECT().Helper()

			var ptr *string

			ensure := ensure.New(mockT)
			ensure(ptr).IsNil()
		})

		t.Run("when not nil", func(t *testing.T) {
			mockT := setupMockTWithCleanupCheck(t)

			val := "not nil"
			ptr := &val

			mockT.EXPECT().Fatalf(failureFormat, ptr).After(
				mockT.EXPECT().Helper(),
			)

			ensure := ensure.New(mockT)
			ensure(ptr).IsNil()
		})
	})

	t.Run("slice", func(t *testing.T) {
		t.Run("when nil", func(t *testing.T) {
			mockT := setupMockTWithCleanupCheck(t)
			mockT.EXPECT().Helper()

			var nilSlice []string

			ensure := ensure.New(mockT)
			ensure(nilSlice).IsNil()
		})

		t.Run("when not nil", func(t *testing.T) {
			mockT := setupMockTWithCleanupCheck(t)

			slice := []string{}

			mockT.EXPECT().Fatalf(failureFormat, slice).After(
				mockT.EXPECT().Helper(),
			)

			ensure := ensure.New(mockT)
			ensure(slice).IsNil()
		})
	})

	t.Run("map", func(t *testing.T) {
		t.Run("when nil", func(t *testing.T) {
			mockT := setupMockTWithCleanupCheck(t)
			mockT.EXPECT().Helper()

			var nilMap map[string]string

			ensure := ensure.New(mockT)
			ensure(nilMap).IsNil()
		})

		t.Run("when not nil", func(t *testing.T) {
			mockT := setupMockTWithCleanupCheck(t)

			m := map[string]string{}

			mockT.EXPECT().Fatalf(failureFormat, m).After(
				mockT.EXPECT().Helper(),
			)

			ensure := ensure.New(mockT)
			ensure(m).IsNil()
		})
	})

	t.Run("func", func(t *testing.T) {
		t.Run("when nil", func(t *testing.T) {
			mockT := setupMockTWithCleanupCheck(t)
			mockT.EXPECT().Helper()

			var nilFunc func(string) string

			ensure := ensure.New(mockT)
			ensure(nilFunc).IsNil()
		})

		t.Run("when not nil", func(t *testing.T) {
			mockT := setupMockTWithCleanupCheck(t)

			f := func(s string) string { return "hello, " + s }

			mockT.EXPECT().Fatalf(failureFormat, gomock.Any()).After(
				mockT.EXPECT().Helper(),
			)

			ensure := ensure.New(mockT)
			ensure(f).IsNil()
		})
	})

	t.Run("chan", func(t *testing.T) {
		t.Run("when nil", func(t *testing.T) {
			mockT := setupMockTWithCleanupCheck(t)
			mockT.EXPECT().Helper()

			var nilChan chan string

			ensure := ensure.New(mockT)
			ensure(nilChan).IsNil()
		})

		t.Run("when not nil", func(t *testing.T) {
			mockT := setupMockTWithCleanupCheck(t)

			c := make(chan string)

			mockT.EXPECT().Fatalf(failureFormat, c).After(
				mockT.EXPECT().Helper(),
			)

			ensure := ensure.New(mockT)
			ensure(c).IsNil()
		})
	})

	t.Run("interface", func(t *testing.T) {
		t.Run("when nil", func(t *testing.T) {
			mockT := setupMockTWithCleanupCheck(t)
			mockT.EXPECT().Helper()

			var nilInterface interface{ Hello(string) string }

			ensure := ensure.New(mockT)
			ensure(nilInterface).IsNil()
		})

		t.Run("when not nil", func(t *testing.T) {
			mockT := setupMockTWithCleanupCheck(t)

			var iface interface{ Hello(string) string } = &ExampleGreeter{}

			mockT.EXPECT().Fatalf(failureFormat, iface).After(
				mockT.EXPECT().Helper(),
			)

			ensure := ensure.New(mockT)
			ensure(iface).IsNil()
		})
	})
}

func TestChainIsNotNil(t *testing.T) {
	t.Run("when not nil", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)
		mockT.EXPECT().Helper()

		const val = "not nil"
		ensure := ensure.New(mockT)
		ensure(val).IsNotNil()
	})

	t.Run("when nil", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)

		mockT.EXPECT().Fatalf("Got nil of type %T, expected it not to be nil", nil).After(
			mockT.EXPECT().Helper(),
		)

		ensure := ensure.New(mockT)
		ensure(nil).IsNotNil()
	})

	t.Run("when nil pointer", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)

		var nilPtr *string

		mockT.EXPECT().Fatalf("Got nil of type %T, expected it not to be nil", nilPtr).After(
			mockT.EXPECT().Helper(),
		)

		ensure := ensure.New(mockT)
		ensure(nilPtr).IsNotNil()
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

	t.Run("when nil pointer equals nil", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)
		mockT.EXPECT().Helper()

		var nilPtr *string

		ensure := ensure.New(mockT)
		ensure(nilPtr).Equals(nil)
	})

	t.Run("when nil map equals empty map", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)
		mockT.EXPECT().Fatalf(errorMessageFormat,
			"Actual does not equal expected:\n - <nil map> != map[]",
			"  map[string]string{}",
			"  map[string]string{}",
		).After(
			mockT.EXPECT().Helper(),
		)

		var nilMap map[string]string

		ensure := ensure.New(mockT)
		ensure(nilMap).Equals(map[string]string{})
	})

	t.Run("when nil slice equals empty slice", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)
		mockT.EXPECT().Fatalf(errorMessageFormat,
			"Actual does not equal expected:\n - <nil slice> != []",
			"  []string(nil)",
			"  []string{}",
		).After(
			mockT.EXPECT().Helper(),
		)

		var nilSlice []string

		ensure := ensure.New(mockT)
		ensure(nilSlice).Equals([]string{})
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

	t.Run("when concurrent", func(t *testing.T) {
		mockT := setupMockT(t)

		mockT.EXPECT().Helper().Times(3 * 2)
		mockT.EXPECT().Cleanup(gomock.Any()).Do(func(fn func()) {
			t.Cleanup(fn)
		}).Times(3)

		ensure := ensure.New(mockT)

		var wg sync.WaitGroup
		wg.Add(3)

		go func() {
			ensure("abc").Equals("abc")
			wg.Done()
		}()

		go func() {
			ensure("xyz").Equals("xyz")
			wg.Done()
		}()

		go func() {
			ensure("qwerty").Equals("qwerty")
			wg.Done()
		}()

		wg.Wait()
	})

	t.Run("string-like equality", func(t *testing.T) {
		const (
			differentTypesErrorFormat = "\nTypes provided to Equals are different: got %s, expected %s\n\nACTUAL:\n%s\n\nEXPECTED:\n%s"
			inequalStringErrorFormat  = "\nActual %s does not equal expected %s:\n\nACTUAL:\n%s\n\nEXPECTED:\n%s"
		)

		t.Run("strings", func(t *testing.T) {
			t.Run("when strings are equal", func(t *testing.T) {
				mockT := setupMockTWithCleanupCheck(t)
				mockT.EXPECT().Helper()

				ensure := ensure.New(mockT)
				ensure("Hello, world!").Equals("Hello, world!")
			})

			t.Run("when strings not equal: expected empty string", func(t *testing.T) {
				mockT := setupMockTWithCleanupCheck(t)
				mockT.EXPECT().Fatalf(inequalStringErrorFormat, "string", "string",
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
				mockT.EXPECT().Fatalf(inequalStringErrorFormat, "string", "string",
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
				mockT.EXPECT().Fatalf(inequalStringErrorFormat, "string", "string",
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
				mockT.EXPECT().Fatalf(inequalStringErrorFormat, "string", "string",
					`  "abc"`,
					`  "abc\n\"xyz\"\n\tqwerty"`, // Formatted with quotes and escaped control characters
				).After(
					mockT.EXPECT().Helper(),
				)

				ensure := ensure.New(mockT)
				ensure("abc").Equals("abc\n\"xyz\"\n\tqwerty")
			})
		})

		t.Run("[]bytes", func(t *testing.T) {
			t.Run("when non-string byte slices are equal", func(t *testing.T) {
				mockT := setupMockTWithCleanupCheck(t)
				mockT.EXPECT().Helper()

				ensure := ensure.New(mockT)
				ensure([]byte{1, 2, 128}).Equals([]byte{1, 2, 128})
			})

			t.Run("when non-string byte slices are not equal", func(t *testing.T) {
				mockT := setupMockTWithCleanupCheck(t)
				mockT.EXPECT().Fatalf(errorMessageFormat, "Actual does not equal expected:\n - slice[2]: 128 != 129",
					"  []uint8{0x1, 0x2, 0x80}",
					"  []uint8{0x1, 0x2, 0x81}",
				).After(
					mockT.EXPECT().Helper(),
				)

				ensure := ensure.New(mockT)
				ensure([]byte{1, 2, 128}).Equals([]byte{1, 2, 129})
			})

			t.Run("when string byte slices are equal", func(t *testing.T) {
				mockT := setupMockTWithCleanupCheck(t)
				mockT.EXPECT().Helper()

				ensure := ensure.New(mockT)
				ensure([]byte("Hello, world!")).Equals([]byte("Hello, world!"))
			})

			t.Run("when string byte slices not equal: expected empty byte slice", func(t *testing.T) {
				mockT := setupMockTWithCleanupCheck(t)
				mockT.EXPECT().Fatalf(inequalStringErrorFormat, "[]byte", "[]byte",
					`  []byte("abc")`,
					"  (empty []byte)",
				).After(
					mockT.EXPECT().Helper(),
				)

				ensure := ensure.New(mockT)
				ensure([]byte("abc")).Equals([]byte(""))
			})

			t.Run("when string byte slices not equal: actual empty byte slice", func(t *testing.T) {
				mockT := setupMockTWithCleanupCheck(t)
				mockT.EXPECT().Fatalf(inequalStringErrorFormat, "[]byte", "[]byte",
					"  (empty []byte)",
					`  []byte("abc")`,
				).After(
					mockT.EXPECT().Helper(),
				)

				ensure := ensure.New(mockT)
				ensure([]byte("")).Equals([]byte("abc"))
			})

			t.Run("when string byte slices not equal: expected contains double quotes, newlines, and tabs", func(t *testing.T) {
				mockT := setupMockTWithCleanupCheck(t)
				mockT.EXPECT().Fatalf(inequalStringErrorFormat, "[]byte", "[]byte",
					`  []byte("abc\n\"xyz\"\n\tqwerty")`, // Formatted with quotes and escaped control characters
					`  []byte("abc")`,
				).After(
					mockT.EXPECT().Helper(),
				)

				ensure := ensure.New(mockT)
				ensure([]byte("abc\n\"xyz\"\n\tqwerty")).Equals([]byte("abc"))
			})

			t.Run("when string byte slices not equal: actual contains double quotes, newlines, and tabs", func(t *testing.T) {
				mockT := setupMockTWithCleanupCheck(t)
				mockT.EXPECT().Fatalf(inequalStringErrorFormat, "[]byte", "[]byte",
					`  []byte("abc")`,
					`  []byte("abc\n\"xyz\"\n\tqwerty")`, // Formatted with quotes and escaped control characters
				).After(
					mockT.EXPECT().Helper(),
				)

				ensure := ensure.New(mockT)
				ensure([]byte("abc")).Equals([]byte("abc\n\"xyz\"\n\tqwerty"))
			})
		})

		t.Run("mixed strings and []bytes", func(t *testing.T) {
			t.Run("when received []byte and expected string", func(t *testing.T) {
				t.Run("and both are empty", func(t *testing.T) {
					mockT := setupMockTWithCleanupCheck(t)
					mockT.EXPECT().Fatalf(differentTypesErrorFormat, "[]byte", "string",
						"  (empty []byte)",
						"  (empty string)",
					).After(
						mockT.EXPECT().Helper(),
					)

					ensure := ensure.New(mockT)
					ensure([]byte("")).Equals("")
				})

				t.Run("and actual is empty", func(t *testing.T) {
					mockT := setupMockTWithCleanupCheck(t)
					mockT.EXPECT().Fatalf(differentTypesErrorFormat, "[]byte", "string",
						"  (empty []byte)",
						`  "Hello, World!"`,
					).After(
						mockT.EXPECT().Helper(),
					)

					ensure := ensure.New(mockT)
					ensure([]byte("")).Equals("Hello, World!")
				})

				t.Run("and expected is empty", func(t *testing.T) {
					mockT := setupMockTWithCleanupCheck(t)
					mockT.EXPECT().Fatalf(differentTypesErrorFormat, "[]byte", "string",
						`  []byte("Hello, World!")`,
						"  (empty string)",
					).After(
						mockT.EXPECT().Helper(),
					)

					ensure := ensure.New(mockT)
					ensure([]byte("Hello, World!")).Equals("")
				})

				t.Run("and both are present and equal", func(t *testing.T) {
					mockT := setupMockTWithCleanupCheck(t)
					mockT.EXPECT().Fatalf(differentTypesErrorFormat, "[]byte", "string",
						`  []byte("Hello, World!")`,
						`  "Hello, World!"`,
					).After(
						mockT.EXPECT().Helper(),
					)

					ensure := ensure.New(mockT)
					ensure([]byte("Hello, World!")).Equals("Hello, World!")
				})

				t.Run("and both are present and not equal", func(t *testing.T) {
					mockT := setupMockTWithCleanupCheck(t)
					mockT.EXPECT().Fatalf(differentTypesErrorFormat, "[]byte", "string",
						`  []byte("Hello")`,
						`  "World"`,
					).After(
						mockT.EXPECT().Helper(),
					)

					ensure := ensure.New(mockT)
					ensure([]byte("Hello")).Equals("World")
				})
			})

			t.Run("when received string and expected []byte", func(t *testing.T) {
				t.Run("and both are empty", func(t *testing.T) {
					mockT := setupMockTWithCleanupCheck(t)
					mockT.EXPECT().Fatalf(differentTypesErrorFormat, "string", "[]byte",
						"  (empty string)",
						"  (empty []byte)",
					).After(
						mockT.EXPECT().Helper(),
					)

					ensure := ensure.New(mockT)
					ensure("").Equals([]byte(""))
				})

				t.Run("and actual is empty", func(t *testing.T) {
					mockT := setupMockTWithCleanupCheck(t)
					mockT.EXPECT().Fatalf(differentTypesErrorFormat, "string", "[]byte",
						"  (empty string)",
						`  []byte("Hello, World!")`,
					).After(
						mockT.EXPECT().Helper(),
					)

					ensure := ensure.New(mockT)
					ensure("").Equals([]byte("Hello, World!"))
				})

				t.Run("and expected is empty", func(t *testing.T) {
					mockT := setupMockTWithCleanupCheck(t)
					mockT.EXPECT().Fatalf(differentTypesErrorFormat, "string", "[]byte",
						`  "Hello, World!"`,
						"  (empty []byte)",
					).After(
						mockT.EXPECT().Helper(),
					)

					ensure := ensure.New(mockT)
					ensure("Hello, World!").Equals([]byte(""))
				})

				t.Run("and both are present and equal", func(t *testing.T) {
					mockT := setupMockTWithCleanupCheck(t)
					mockT.EXPECT().Fatalf(differentTypesErrorFormat, "string", "[]byte",
						`  "Hello, World!"`,
						`  []byte("Hello, World!")`,
					).After(
						mockT.EXPECT().Helper(),
					)

					ensure := ensure.New(mockT)
					ensure("Hello, World!").Equals([]byte("Hello, World!"))
				})

				t.Run("and both are present and not equal", func(t *testing.T) {
					mockT := setupMockTWithCleanupCheck(t)
					mockT.EXPECT().Fatalf(differentTypesErrorFormat, "string", "[]byte",
						`  "Hello"`,
						`  []byte("World")`,
					).After(
						mockT.EXPECT().Helper(),
					)

					ensure := ensure.New(mockT)
					ensure("Hello").Equals([]byte("World"))
				})
			})
		})
	})
}

func TestChainIsEmpty(t *testing.T) {
	testEmptyChain(t, func(t *testing.T, valueLength int, value interface{}) {
		mockT := setupMockTWithCleanupCheck(t)

		if valueLength == 0 {
			mockT.EXPECT().Helper()
		} else {
			mockT.EXPECT().Fatalf("Got %+v with length %d, expected it to be empty", value, valueLength).After(
				mockT.EXPECT().Helper(),
			)
		}

		ensure := ensure.New(mockT)
		ensure(value).IsEmpty()
	})

	t.Run("when not valid type", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)

		mockT.EXPECT().Fatalf("Got type int, expected array, slice, string, or map").After(
			mockT.EXPECT().Helper(),
		)

		ensure := ensure.New(mockT)
		ensure(1234).IsEmpty()
	})
}

func TestChainIsNotEmpty(t *testing.T) {
	testEmptyChain(t, func(t *testing.T, valueLength int, value interface{}) {
		mockT := setupMockTWithCleanupCheck(t)

		if valueLength == 0 {
			mockT.EXPECT().Fatalf("Got %+v, expected it to not be empty", value).After(
				mockT.EXPECT().Helper(),
			)
		} else {
			mockT.EXPECT().Helper()
		}

		ensure := ensure.New(mockT)
		ensure(value).IsNotEmpty()
	})

	t.Run("when not valid type", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)

		mockT.EXPECT().Fatalf("Got type int, expected array, slice, string, or map").After(
			mockT.EXPECT().Helper(),
		)

		ensure := ensure.New(mockT)
		ensure(1234).IsNotEmpty()
	})
}

func testEmptyChain(t *testing.T, run func(t *testing.T, valueLength int, value interface{})) {
	table := []struct {
		Name        string
		ValueLength int
		Value       interface{}
	}{
		{
			Name:        "when empty: array",
			ValueLength: 0,
			Value:       [0]string{},
		},
		{
			Name:        "when not empty: array",
			ValueLength: 2,
			Value:       [2]string{"1", "2"},
		},
		{
			Name:        "when empty: slice",
			ValueLength: 0,
			Value:       []string{},
		},
		{
			Name:        "when not empty: slice",
			ValueLength: 1,
			Value:       []string{"1"},
		},
		{
			Name:        "when empty: string",
			ValueLength: 0,
			Value:       "",
		},
		{
			Name:        "when not empty: string",
			ValueLength: len("not empty"),
			Value:       "not empty",
		},
		{
			Name:        "when empty: map",
			ValueLength: 0,
			Value:       map[string]string{},
		},
		{
			Name:        "when not empty: map",
			ValueLength: 1,
			Value:       map[string]string{"hello": "world"},
		},
	}

	for _, entry := range table {
		entry := entry // Pin range variable

		t.Run(entry.Name, func(t *testing.T) {
			run(t, entry.ValueLength, entry.Value)
		})
	}
}

func TestChainContains(t *testing.T) {
	testContainsChain(t, func(t *testing.T, doesContain bool, actual, expected interface{}, formattedActual, formattedExpected string) {
		mockT := setupMockTWithCleanupCheck(t)

		if doesContain {
			mockT.EXPECT().Helper()
		} else {
			mockT.EXPECT().Fatalf("Actual does not contain expected:\n\nACTUAL:\n%s\n\nEXPECTED TO CONTAIN:\n%s", formattedActual, formattedExpected).After(
				mockT.EXPECT().Helper(),
			)
		}

		ensure := ensure.New(mockT)
		ensure(actual).Contains(expected)
	})

	t.Run("when not valid type", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)

		mockT.EXPECT().Fatalf("Got type int, expected string, array, or slice").After(
			mockT.EXPECT().Helper(),
		)

		ensure := ensure.New(mockT)
		ensure(1234).Contains(2)
	})

	t.Run("when string is expected to contain a non-string type", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)

		mockT.EXPECT().Fatalf("Got string, but expected is a int, and a string can only contain other strings").After(
			mockT.EXPECT().Helper(),
		)

		ensure := ensure.New(mockT)
		ensure("hello").DoesNotContain(123)
	})
}

func TestChainDoesNotContain(t *testing.T) {
	testContainsChain(t, func(t *testing.T, doesContain bool, actual, expected interface{}, formattedActual, formattedExpected string) {
		mockT := setupMockTWithCleanupCheck(t)

		if doesContain {
			mockT.EXPECT().Fatalf("Actual contains expected, but did not expect it to:\n\nACTUAL:\n%s\n\nEXPECTED NOT TO CONTAIN:\n%s", formattedActual, formattedExpected).After(
				mockT.EXPECT().Helper(),
			)
		} else {
			mockT.EXPECT().Helper()
		}

		ensure := ensure.New(mockT)
		ensure(actual).DoesNotContain(expected)
	})

	t.Run("when not valid type", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)

		mockT.EXPECT().Fatalf("Got type int, expected string, array, or slice").After(
			mockT.EXPECT().Helper(),
		)

		ensure := ensure.New(mockT)
		ensure(1234).DoesNotContain(2)
	})

	t.Run("when string is expected to not contain a non-string type", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)

		mockT.EXPECT().Fatalf("Got string, but expected is a int, and a string can only contain other strings").After(
			mockT.EXPECT().Helper(),
		)

		ensure := ensure.New(mockT)
		ensure("hello").DoesNotContain(123)
	})
}

func testContainsChain(t *testing.T, run func(t *testing.T, doesContain bool, actual, expected interface{}, formattedActual, formattedExpected string)) {
	table := []struct {
		Name        string
		Actual      interface{}
		Expected    interface{}
		DoesContain bool

		FormattedActual   string
		FormattedExpected string
	}{
		{
			Name:        "when contains: string",
			Actual:      "hello",
			Expected:    "ell",
			DoesContain: true,

			FormattedActual:   `  "hello"`, // Indented
			FormattedExpected: `  "ell"`,   // Indented
		},
		{
			Name:        "when does not contain: string",
			Actual:      "hello",
			Expected:    "zzz",
			DoesContain: false,

			FormattedActual:   `  "hello"`, // Indented
			FormattedExpected: `  "zzz"`,   // Indented
		},
		{
			Name:        "when contains: array",
			Actual:      [2]string{"abc", "xyz"},
			Expected:    "xyz",
			DoesContain: true,

			FormattedActual:   `  [2]string{"abc", "xyz"}`, // Indented
			FormattedExpected: `  "xyz"`,                   // Indented
		},
		{
			Name:        "when does not contain: array",
			Actual:      [2]string{"abc", "xyz"},
			Expected:    "qwerty",
			DoesContain: false,

			FormattedActual:   `  [2]string{"abc", "xyz"}`, // Indented
			FormattedExpected: `  "qwerty"`,                // Indented
		},
		{
			Name:        "when contains: slice",
			Actual:      []int{123, 456},
			Expected:    123,
			DoesContain: true,

			FormattedActual:   `  []int{123, 456}`, // Indented
			FormattedExpected: `  int(123)`,        // Indented
		},
		{
			Name:        "when does not contain: slice",
			Actual:      []int{123, 456},
			Expected:    789,
			DoesContain: false,

			FormattedActual:   `  []int{123, 456}`, // Indented
			FormattedExpected: `  int(789)`,        // Indented
		},
	}

	for _, entry := range table {
		entry := entry // Pin range variable

		t.Run(entry.Name, func(t *testing.T) {
			run(t, entry.DoesContain, entry.Actual, entry.Expected, entry.FormattedActual, entry.FormattedExpected)
		})
	}
}

func TestChainMatchesRegexp(t *testing.T) {
	t.Run("with valid complete match", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)
		mockT.EXPECT().Helper()

		ensure := ensure.New(mockT)
		ensure("hello 123 world").MatchesRegexp("^hello [1-3]+ world$")
	})

	t.Run("with valid partial match", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)
		mockT.EXPECT().Helper()

		ensure := ensure.New(mockT)
		ensure("hello 123 world").MatchesRegexp("[1-3]+")
	})

	t.Run("with no match", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)

		mockT.EXPECT().Fatalf(
			"Actual does not match regular expression:\n\nACTUAL:\n%s\n\nEXPECTED TO MATCH:\n%s",
			`  "hello 1-3 world"`,      // Indented
			`  "^hello [1-3]+ world$"`, // Indented
		).After(
			mockT.EXPECT().Helper(),
		)

		ensure := ensure.New(mockT)
		ensure("hello 1-3 world").MatchesRegexp("^hello [1-3]+ world$")
	})

	t.Run("when pattern is empty", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)

		mockT.EXPECT().Fatalf("Cannot match against an empty pattern").After(
			mockT.EXPECT().Helper(),
		)

		ensure := ensure.New(mockT)
		ensure("hello").MatchesRegexp("")
	})

	t.Run("when actual is not a string", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)

		mockT.EXPECT().Fatalf("Actual is not a string, it's a %T", 123).After(
			mockT.EXPECT().Helper(),
		)

		ensure := ensure.New(mockT)
		ensure(123).MatchesRegexp("hello")
	})

	t.Run("when regular expression is invalid", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)

		mockT.EXPECT().Fatalf("Unable to compile regular expression: %s\nERROR: %v", "[", gomock.Any()).After(
			mockT.EXPECT().Helper(),
		)

		ensure := ensure.New(mockT)
		ensure("hello").MatchesRegexp("[") // Missing closing ]
	})
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

type ExampleGreeter struct{}

func (*ExampleGreeter) Hello(s string) string {
	return "hello, " + s
}
