package ensurepkg_test

import (
	"testing"

	"github.com/JosiahWitt/ensure"
	"github.com/JosiahWitt/ensure/ensurepkg"
	"github.com/golang/mock/gomock"
)

type runTableTestEntryGroup struct {
	Prefix  string
	Entries []runTableTestEntry
}

type runTableTestEntry struct {
	Name                 string
	ExpectedNames        []string
	ExpectedFatalMessage string
	Table                interface{}
	CheckEntry           func(t *testing.T, rawEntry interface{})
}

func TestEnsureRunTableByIndex(t *testing.T) {
	runTableTests := runTableTests{}

	groups := []runTableTestEntryGroup{
		runTableTests.general(),
		runTableTests.mocksField(),
		runTableTests.setupMocksField(),
		runTableTests.subjectField(),
	}

	table := []runTableTestEntry{}
	for _, group := range groups {
		for _, entry := range group.Entries {
			entry.Name = group.Prefix + ": " + entry.Name
			table = append(table, entry)
		}
	}

	for _, entry := range table {
		entry := entry // Pin range variable

		t.Run(entry.Name, func(t *testing.T) {
			mockT := setupMockT(t)
			expectedTableSize := len(entry.ExpectedNames)

			mockT.EXPECT().Helper().Times(3 * expectedTableSize) // 3 = RunTableByIndex + run + before Cleanup call
			mockT.EXPECT().Cleanup(gomock.Any()).Times(expectedTableSize)

			// Build expected Run calls only if there's no expected error
			expectedTestingInputs := []*testing.T{}
			if entry.ExpectedFatalMessage == "" {
				expectedRunCalls := []*gomock.Call{}

				for _, name := range entry.ExpectedNames {
					providedTestingInput := &testing.T{}
					expectedTestingInputs = append(expectedTestingInputs, providedTestingInput)

					expectedRunCalls = append(expectedRunCalls,
						mockT.EXPECT().Run(name, gomock.Any()).
							Do(func(name string, fn func(t *testing.T)) {
								fn(providedTestingInput)
							}),
					)
				}

				// Run calls should be in order
				gomock.InOrder(expectedRunCalls...)
			} else {
				gomock.InOrder(
					mockT.EXPECT().Helper(),
					mockT.EXPECT().Cleanup(gomock.Any()),
					mockT.EXPECT().Helper(),
					mockT.EXPECT().Fatalf(entry.ExpectedFatalMessage),
				)
			}

			type Params struct {
				ensure ensurepkg.Ensure
				i      int
			}

			// Run table and save parameters
			actualParams := []Params{}
			ensure := ensure.New(mockT)
			ensure.RunTableByIndex(entry.Table, func(ensure ensurepkg.Ensure, i int) {
				actualParams = append(actualParams, Params{ensure: ensure, i: i})
			})

			// Verify call count
			if len(actualParams) != expectedTableSize {
				t.Fatalf("len(actualParams) != expectedTableSize")
			}

			// Verify parameters are correct
			for i, actualParam := range actualParams {
				if actualParam.ensure.T() != expectedTestingInputs[i] {
					t.Fatalf("actualParams[%d].ensure.T() != expectedTestingInputs[%d]", i, i)
				}

				if actualParam.i != i {
					t.Fatalf("actualParams[%d].i != %d", i, i)
				}
			}

			if entry.CheckEntry != nil {
				entry.CheckEntry(t, entry.Table)
			}
		})
	}
}

type runTableTests struct{}

func (runTableTests) general() runTableTestEntryGroup {
	return runTableTestEntryGroup{
		Prefix: "general",
		Entries: []runTableTestEntry{
			{
				Name:          "with valid table: slice",
				ExpectedNames: []string{"name 1", "name 2", "name 3"},
				Table: []struct {
					Name  string
					Value string
				}{
					{
						Name:  "name 1",
						Value: "item 1",
					},
					{
						Name:  "name 2",
						Value: "item 2",
					},
					{
						Name:  "name 3",
						Value: "item 3",
					},
				},
			},

			{
				Name:          "with valid table: array",
				ExpectedNames: []string{"name 1", "name 2"},
				Table: [2]struct {
					Name  string
					Value string
				}{
					{
						Name:  "name 1",
						Value: "item 1",
					},
					{
						Name:  "name 2",
						Value: "item 2",
					},
				},
			},

			{
				Name:                 "with invalid table type: not array or slice",
				ExpectedFatalMessage: "Expected a slice or array for the table, got string",
				Table:                "my table",
			},

			{
				Name:                 "with invalid table type: not array or slice of stucts",
				ExpectedFatalMessage: "Expected entry in table to be a struct, got string",
				Table: []string{
					"item 1",
					"item 2",
				},
			},

			{
				Name:                 "with missing name",
				ExpectedFatalMessage: "Name field does not exist on struct in table",
				Table: []struct {
					Value string
				}{
					{
						Value: "item 1",
					},
					{
						Value: "item 2",
					},
				},
			},

			{
				Name:                 "with name with invalid type",
				ExpectedFatalMessage: "Name field in struct in table is not a string",
				Table: []struct {
					Name  int
					Value string
				}{
					{
						Name:  1,
						Value: "item 1",
					},
					{
						Name:  2,
						Value: "item 2",
					},
				},
			},

			{
				Name:                 "with missing name for one item",
				ExpectedFatalMessage: "Errors encountered while building table:\n - table[1]: Name not set for item",
				Table: []struct {
					Name  string
					Value string
				}{
					{
						Name:  "name 1",
						Value: "item 1",
					},
					{
						Name:  "",
						Value: "item 2",
					},
				},
			},

			{
				Name:                 "with duplicate name",
				ExpectedFatalMessage: "Errors encountered while building table:\n - table[2]: duplicate Name found; first occurrence was table[0].Name: name 1",
				Table: []struct {
					Name  string
					Value string
				}{
					{
						Name:  "name 1",
						Value: "item 1",
					},
					{
						Name:  "name 2",
						Value: "item 2",
					},
					{
						Name:  "name 1",
						Value: "item 3",
					},
				},
			},

			{
				Name:                 "with double duplicate name",
				ExpectedFatalMessage: "Errors encountered while building table:\n - table[2]: duplicate Name found; first occurrence was table[0].Name: name 1\n - table[3]: duplicate Name found; first occurrence was table[0].Name: name 1",
				Table: []struct {
					Name  string
					Value string
				}{
					{
						Name:  "name 1",
						Value: "item 1",
					},
					{
						Name:  "name 2",
						Value: "item 2",
					},
					{
						Name:  "name 1",
						Value: "item 3",
					},
					{
						Name:  "name 1",
						Value: "item 4",
					},
				},
			},
		},
	}
}

func (runTableTests) mocksField() runTableTestEntryGroup {
	type (
		OneMockMissingNEWMethod struct {
			Valid1  *ExampleMockValid1
			Invalid *struct{ Nothing bool }
			Valid2  *ExampleMockValid2
		}

		OneMockNEWMethodZeroParams struct {
			Valid1  *ExampleMockValid1
			Invalid *ExampleMockNEWMethodZeroParams
			Valid2  *ExampleMockValid2
		}

		OneMockNEWMethodIncorrectParam struct {
			Valid1  *ExampleMockValid1
			Invalid *ExampleMockNEWMethodIncorrectParam
			Valid2  *ExampleMockValid2
		}

		OneMockNEWMethodZeroReturns struct {
			Valid1  *ExampleMockValid1
			Invalid *ExampleMockNEWMethodZeroReturns
			Valid2  *ExampleMockValid2
		}

		OneMockNEWMethodIncorrectReturn struct {
			Valid1  *ExampleMockValid1
			Invalid *ExampleMockNEWMethodIncorrectReturn
			Valid2  *ExampleMockValid2
		}

		OneMockNotPointer struct {
			Valid1  *ExampleMockValid1
			Invalid ExampleMockValid1
			Valid2  *ExampleMockValid2
		}

		DuplicateMocks struct {
			Valid1          *ExampleMockValid1
			Valid1Duplicate *ExampleMockValid1
		}
	)

	return runTableTestEntryGroup{
		Prefix: "Mocks field",
		Entries: []runTableTestEntry{
			{
				Name:          "when valid",
				ExpectedNames: []string{"name 1", "name 2"},
				Table: []struct {
					Name  string
					Mocks *TwoValidMocks
				}{
					{
						Name: "name 1",
					},
					{
						Name: "name 2",
					},
				},

				CheckEntry: func(t *testing.T, rawTable interface{}) {
					table := rawTable.([]struct {
						Name  string
						Mocks *TwoValidMocks
					})

					for _, entry := range table {
						entry.Mocks.check(t)
					}
				},
			},

			{
				Name:                 "when not pointer to mock struct",
				ExpectedFatalMessage: "Mocks field should be a pointer to a struct, got ensurepkg_test.TwoValidMocks",
				Table: []struct {
					Name  string
					Mocks TwoValidMocks
				}{
					{
						Name: "name 1",
					},
					{
						Name: "name 2",
					},
				},
			},

			{
				Name:                 "when pointer to non struct",
				ExpectedFatalMessage: "Mocks field should be a pointer to a struct, got *string",
				Table: []struct {
					Name  string
					Mocks *string
				}{
					{
						Name: "name 1",
					},
					{
						Name: "name 2",
					},
				},
			},

			{
				Name: "when missing NEW method",
				ExpectedFatalMessage: "\nMocks.Invalid is missing the NEW method. Expected:\n\tfunc (*struct { Nothing bool }) NEW(*gomock.Controller) *struct { Nothing bool }" +
					"\nPlease ensure you generated the mocks using the `ensure generate mocks` command.",
				Table: []struct {
					Name  string
					Mocks *OneMockMissingNEWMethod
				}{
					{
						Name: "name 1",
					},
					{
						Name: "name 2",
					},
				},
			},

			{
				Name:                 "when NEW method has zero params",
				ExpectedFatalMessage: "\nMocks.Invalid.NEW has this method signature:\n\tfunc() *ensurepkg_test.ExampleMockNEWMethodZeroParams\nExpected:\n\tfunc(*gomock.Controller) *ensurepkg_test.ExampleMockNEWMethodZeroParams",
				Table: []struct {
					Name  string
					Mocks *OneMockNEWMethodZeroParams
				}{
					{
						Name: "name 1",
					},
					{
						Name: "name 2",
					},
				},
			},

			{
				Name:                 "when NEW method has incorrect param",
				ExpectedFatalMessage: "\nMocks.Invalid.NEW has this method signature:\n\tfunc(string) *ensurepkg_test.ExampleMockNEWMethodIncorrectParam\nExpected:\n\tfunc(*gomock.Controller) *ensurepkg_test.ExampleMockNEWMethodIncorrectParam",
				Table: []struct {
					Name  string
					Mocks *OneMockNEWMethodIncorrectParam
				}{
					{
						Name: "name 1",
					},
					{
						Name: "name 2",
					},
				},
			},

			{
				Name:                 "when NEW method has zero returns",
				ExpectedFatalMessage: "\nMocks.Invalid.NEW has this method signature:\n\tfunc(*gomock.Controller)\nExpected:\n\tfunc(*gomock.Controller) *ensurepkg_test.ExampleMockNEWMethodZeroReturns",
				Table: []struct {
					Name  string
					Mocks *OneMockNEWMethodZeroReturns
				}{
					{
						Name: "name 1",
					},
					{
						Name: "name 2",
					},
				},
			},

			{
				Name:                 "when NEW method has incorrect return",
				ExpectedFatalMessage: "\nMocks.Invalid.NEW has this method signature:\n\tfunc(*gomock.Controller) string\nExpected:\n\tfunc(*gomock.Controller) *ensurepkg_test.ExampleMockNEWMethodIncorrectReturn",
				Table: []struct {
					Name  string
					Mocks *OneMockNEWMethodIncorrectReturn
				}{
					{
						Name: "name 1",
					},
					{
						Name: "name 2",
					},
				},
			},

			{
				Name:                 "when mock is not a pointer",
				ExpectedFatalMessage: "Mocks.Invalid should be a pointer to a struct, got ensurepkg_test.ExampleMockValid1",
				Table: []struct {
					Name  string
					Mocks *OneMockNotPointer
				}{
					{
						Name: "name 1",
					},
					{
						Name: "name 2",
					},
				},
			},

			{
				Name:                 "with duplicate mock",
				ExpectedFatalMessage: "Found multiple mocks with type '*ensurepkg_test.ExampleMockValid1'; only one mock of each type is allowed",
				Table: []struct {
					Name  string
					Mocks *DuplicateMocks
				}{
					{
						Name: "name 1",
					},
					{
						Name: "name 2",
					},
				},
			},
		},
	}
}

func (runTableTests) setupMocksField() runTableTestEntryGroup {
	return runTableTestEntryGroup{
		Prefix: "SetupMocks field",
		Entries: []runTableTestEntry{
			{
				Name:          "with valid function",
				ExpectedNames: []string{"name 1", "name 2"},
				Table: []struct {
					Name       string
					Mocks      *TwoValidMocks
					SetupMocks func(*TwoValidMocks)
				}{
					{
						Name: "name 1",
						SetupMocks: func(tvm *TwoValidMocks) {
							tvm.Valid1.CustomField = "updated name 1"
						},
					},
					{
						Name: "name 2",
						SetupMocks: func(tvm *TwoValidMocks) {
							tvm.Valid1.CustomField = "updated name 2"
						},
					},
				},

				CheckEntry: func(t *testing.T, rawTable interface{}) {
					table := rawTable.([]struct {
						Name       string
						Mocks      *TwoValidMocks
						SetupMocks func(*TwoValidMocks)
					})

					for _, entry := range table {
						entry.Mocks.check(t)
						isTrue(t, entry.Mocks.Valid1.CustomField == "updated "+entry.Name)
					}
				},
			},

			{
				Name:          "with function not present for one",
				ExpectedNames: []string{"name 1", "name 2"},
				Table: []struct {
					Name       string
					Mocks      *TwoValidMocks
					SetupMocks func(*TwoValidMocks)
				}{
					{
						Name: "name 1",
						SetupMocks: func(tvm *TwoValidMocks) {
							tvm.Valid1.CustomField = "updated name 1"
						},
					},
					{
						Name: "name 2",
					},
				},

				CheckEntry: func(t *testing.T, rawTable interface{}) {
					table := rawTable.([]struct {
						Name       string
						Mocks      *TwoValidMocks
						SetupMocks func(*TwoValidMocks)
					})

					isTrue(t, table[0].Mocks.Valid1.CustomField == "updated name 1")
					isTrue(t, table[1].Mocks.Valid1.CustomField == "")
				},
			},

			{
				Name:                 "without Mocks field",
				ExpectedFatalMessage: "SetupMocks field requires the Mocks field",
				Table: []struct {
					Name       string
					SetupMocks func(*TwoValidMocks)
				}{
					{
						Name:       "name 1",
						SetupMocks: func(*TwoValidMocks) {},
					},
					{
						Name:       "name 2",
						SetupMocks: func(*TwoValidMocks) {},
					},
				},
			},

			{
				Name:                 "function missing param",
				ExpectedFatalMessage: "\nSetupMocks has this function signature:\n\tfunc()\nExpected:\n\tfunc(*ensurepkg_test.TwoValidMocks)",
				Table: []struct {
					Name       string
					Mocks      *TwoValidMocks
					SetupMocks func()
				}{
					{
						Name:       "name 1",
						SetupMocks: func() {},
					},
					{
						Name:       "name 2",
						SetupMocks: func() {},
					},
				},
			},

			{
				Name:                 "function with invalid param",
				ExpectedFatalMessage: "\nSetupMocks has this function signature:\n\tfunc(*string)\nExpected:\n\tfunc(*ensurepkg_test.TwoValidMocks)",
				Table: []struct {
					Name       string
					Mocks      *TwoValidMocks
					SetupMocks func(*string)
				}{
					{
						Name:       "name 1",
						SetupMocks: func(*string) {},
					},
					{
						Name:       "name 2",
						SetupMocks: func(*string) {},
					},
				},
			},

			{
				Name:                 "function with a return",
				ExpectedFatalMessage: "\nSetupMocks has this function signature:\n\tfunc(*ensurepkg_test.TwoValidMocks) error\nExpected:\n\tfunc(*ensurepkg_test.TwoValidMocks)",
				Table: []struct {
					Name       string
					Mocks      *TwoValidMocks
					SetupMocks func(*TwoValidMocks) error
				}{
					{
						Name:       "name 1",
						SetupMocks: func(*TwoValidMocks) error { return nil },
					},
					{
						Name:       "name 2",
						SetupMocks: func(*TwoValidMocks) error { return nil },
					},
				},
			},
		},
	}
}

func (runTableTests) subjectField() runTableTestEntryGroup {
	type (
		IntAdder interface {
			Add(a, b int) int
		}

		AdderSubject struct {
			Adder IntAdder
		}

		AdderSubjectWithDuplicate struct {
			Adder1 IntAdder
			Adder2 IntAdder
		}

		AdderSubjectWithExtraField struct {
			Adder      IntAdder
			ExtraField string
		}

		AdderSubjectWithUnmockedInterface struct {
			Adder             IntAdder
			UnmockedInterface interface{ Multiply(a, b int) int }
		}

		SubjectMatchingMultipleMocks struct {
			Subber interface{ Sub(a, b int) int }
		}
	)

	return runTableTestEntryGroup{
		Prefix: "Subject field",
		Entries: []runTableTestEntry{
			{
				Name:          "when valid",
				ExpectedNames: []string{"name 1", "name 2"},
				Table: []struct {
					Name    string
					Mocks   *TwoValidMocks
					Subject *AdderSubject
				}{
					{
						Name: "name 1",
					},
					{
						Name: "name 2",
					},
				},

				CheckEntry: func(t *testing.T, rawTable interface{}) {
					table := rawTable.([]struct {
						Name    string
						Mocks   *TwoValidMocks
						Subject *AdderSubject
					})

					for _, entry := range table {
						entry.Mocks.check(t)
						isTrue(t, entry.Subject.Adder.Add(1, 2) == 3)
					}
				},
			},

			{
				Name:          "when duplicate interfaces",
				ExpectedNames: []string{"name 1", "name 2"},
				Table: []struct {
					Name    string
					Mocks   *TwoValidMocks
					Subject *AdderSubjectWithDuplicate
				}{
					{
						Name: "name 1",
					},
					{
						Name: "name 2",
					},
				},

				CheckEntry: func(t *testing.T, rawTable interface{}) {
					table := rawTable.([]struct {
						Name    string
						Mocks   *TwoValidMocks
						Subject *AdderSubjectWithDuplicate
					})

					for _, entry := range table {
						entry.Mocks.check(t)

						isTrue(t, entry.Subject.Adder1.Add(1, 2) == 3)
						isTrue(t, entry.Subject.Adder2.Add(1, 2) == 3)
						isTrue(t, entry.Subject.Adder1 == entry.Subject.Adder2) // Should point to the same mock
					}
				},
			},

			{
				Name:                 "when not pointer to struct",
				ExpectedFatalMessage: "Subject field should be a pointer to a struct, got ensurepkg_test.AdderSubject",
				Table: []struct {
					Name    string
					Mocks   *TwoValidMocks
					Subject AdderSubject
				}{
					{
						Name: "name 1",
					},
					{
						Name: "name 2",
					},
				},
			},

			{
				Name:                 "when pointer to non struct",
				ExpectedFatalMessage: "Subject field should be a pointer to a struct, got *string",
				Table: []struct {
					Name    string
					Mocks   *TwoValidMocks
					Subject *string
				}{
					{
						Name: "name 1",
					},
					{
						Name: "name 2",
					},
				},
			},

			{
				Name:          "when field contains a non-interface field",
				ExpectedNames: []string{"name 1", "name 2"},
				Table: []struct {
					Name    string
					Mocks   *TwoValidMocks
					Subject *AdderSubjectWithExtraField
				}{
					{
						Name: "name 1",
					},
					{
						Name: "name 2",
					},
				},

				CheckEntry: func(t *testing.T, rawTable interface{}) {
					table := rawTable.([]struct {
						Name    string
						Mocks   *TwoValidMocks
						Subject *AdderSubjectWithExtraField
					})

					for _, entry := range table {
						entry.Mocks.check(t)

						isTrue(t, entry.Subject.Adder.Add(1, 2) == 3)
						isTrue(t, entry.Subject.ExtraField == "")
					}
				},
			},

			{
				Name:          "when field contains a non-mocked interface",
				ExpectedNames: []string{"name 1", "name 2"},
				Table: []struct {
					Name    string
					Mocks   *TwoValidMocks
					Subject *AdderSubjectWithUnmockedInterface
				}{
					{
						Name: "name 1",
					},
					{
						Name: "name 2",
					},
				},

				CheckEntry: func(t *testing.T, rawTable interface{}) {
					table := rawTable.([]struct {
						Name    string
						Mocks   *TwoValidMocks
						Subject *AdderSubjectWithUnmockedInterface
					})

					for _, entry := range table {
						entry.Mocks.check(t)

						isTrue(t, entry.Subject.Adder.Add(1, 2) == 3)
						isTrue(t, entry.Subject.UnmockedInterface == nil)
					}
				},
			},

			{
				Name:                 "when entry matches multiple mocks",
				ExpectedFatalMessage: "Subject.Subber matches multiple mocks; only one mock should exist for each interface: *ensurepkg_test.ExampleMockValid1, *ensurepkg_test.ExampleMockValid2",
				Table: []struct {
					Name    string
					Mocks   *TwoValidMocks
					Subject *SubjectMatchingMultipleMocks
				}{
					{
						Name: "name 1",
					},
					{
						Name: "name 2",
					},
				},
			},
		},
	}
}

func isTrue(t *testing.T, value bool) {
	t.Helper()

	if !value {
		t.Errorf("value is not true")
	}
}

type TwoValidMocks struct {
	Valid1 *ExampleMockValid1
	Valid2 *ExampleMockValid2
}

func (tvm *TwoValidMocks) check(t *testing.T) {
	t.Helper()

	isTrue(t, tvm.Valid1.WasInitialized)
	isTrue(t, tvm.Valid2.WasInitialized)
	isTrue(t, tvm.Valid1.GoMockController == tvm.Valid2.GoMockController) // Ensure GoMock Controller is memoized
}

type ExampleMockValid1 struct {
	WasInitialized   bool
	GoMockController *gomock.Controller
	CustomField      string
}

func (m *ExampleMockValid1) NEW(ctrl *gomock.Controller) *ExampleMockValid1 {
	if ctrl == nil {
		panic("GoMock controller is nil")
	}

	return &ExampleMockValid1{WasInitialized: true, GoMockController: ctrl}
}

func (m *ExampleMockValid1) Add(a, b int) int {
	return a + b
}

func (m *ExampleMockValid1) Sub(a, b int) int {
	return a - b
}

type ExampleMockValid2 struct {
	WasInitialized   bool
	GoMockController *gomock.Controller
	CustomField      string
}

func (m *ExampleMockValid2) NEW(ctrl *gomock.Controller) *ExampleMockValid2 {
	if ctrl == nil {
		panic("GoMock controller is nil")
	}

	return &ExampleMockValid2{WasInitialized: true, GoMockController: ctrl}
}

func (m *ExampleMockValid2) Sub(a, b int) int {
	return a - b
}

type ExampleMockNEWMethodZeroParams struct{}

func (m *ExampleMockNEWMethodZeroParams) NEW() *ExampleMockNEWMethodZeroParams { return nil }

type ExampleMockNEWMethodIncorrectParam struct{}

func (m *ExampleMockNEWMethodIncorrectParam) NEW(notGoMockCtrl string) *ExampleMockNEWMethodIncorrectParam {
	return nil
}

type ExampleMockNEWMethodZeroReturns struct{}

func (m *ExampleMockNEWMethodZeroReturns) NEW(ctrl *gomock.Controller) {}

type ExampleMockNEWMethodIncorrectReturn struct{}

func (m *ExampleMockNEWMethodIncorrectReturn) NEW(ctrl *gomock.Controller) string { return "" }
