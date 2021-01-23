package ensurepkg_test

import (
	"testing"

	"github.com/JosiahWitt/ensure"
	"github.com/JosiahWitt/ensure/ensurepkg"
	"github.com/golang/mock/gomock"
)

func TestEnsureRunTableByIndex(t *testing.T) {
	type (
		TwoValidMocks struct {
			Valid1 *ExampleMockValid1
			Valid2 *ExampleMockValid2
		}

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

	table := []struct {
		Name                 string
		ExpectedNames        []string
		ExpectedFatalMessage string
		Table                interface{}
		CheckEntry           func(t *testing.T, rawEntry interface{})
	}{
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

		// ********** Mocks field ********** //

		{
			Name:          "with mocks: when valid",
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
					isTrue(t, entry.Mocks.Valid1.WasInitialized)
					isTrue(t, entry.Mocks.Valid2.WasInitialized)
					isTrue(t, entry.Mocks.Valid1.GoMockController == entry.Mocks.Valid2.GoMockController) // Ensure GoMock Controller is memoized
				}
			},
		},

		{
			Name:                 "with mocks: when not pointer to mock struct",
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
			Name:                 "with mocks: when pointer to non struct",
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
			Name: "with mocks: when missing NEW method",
			ExpectedFatalMessage: "\nMocks.Invalid is missing the NEW method. Expected:\n\tfunc(*gomock.Controller) *struct { Nothing bool }" +
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
			Name:                 "with mocks: when NEW method has zero params",
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
			Name:                 "with mocks: when NEW method has incorrect param",
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
			Name:                 "with mocks: when NEW method has zero returns",
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
			Name:                 "with mocks: when NEW method has incorrect return",
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
			Name:                 "with mocks: when mock is not a pointer",
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
			Name:                 "with mocks: with duplicate mock",
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

		// ********** SetupMocks field ********** //

		{
			Name:          "with mocks: with valid SetupMocks function",
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
					isTrue(t, entry.Mocks.Valid1.WasInitialized)
					isTrue(t, entry.Mocks.Valid2.WasInitialized)
					isTrue(t, entry.Mocks.Valid1.GoMockController == entry.Mocks.Valid2.GoMockController) // Ensure GoMock Controller is memoized

					isTrue(t, entry.Mocks.Valid1.CustomField == "updated "+entry.Name)
				}
			},
		},

		{
			Name:          "with mocks: with SetupMocks function not present for one",
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
			Name:                 "with mocks: SetupMocks without Mocks",
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
			Name:                 "with mocks: SetupMocks with no param",
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
			Name:                 "with mocks: SetupMocks with invalid param",
			ExpectedFatalMessage: "\nSetupMocks has this function signature:\n\tfunc(*ensurepkg_test.OneMockNotPointer)\nExpected:\n\tfunc(*ensurepkg_test.TwoValidMocks)",
			Table: []struct {
				Name       string
				Mocks      *TwoValidMocks
				SetupMocks func(*OneMockNotPointer)
			}{
				{
					Name:       "name 1",
					SetupMocks: func(*OneMockNotPointer) {},
				},
				{
					Name:       "name 2",
					SetupMocks: func(*OneMockNotPointer) {},
				},
			},
		},

		{
			Name:                 "with mocks: SetupMocks with a return",
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

func isTrue(t *testing.T, value bool) {
	t.Helper()

	if !value {
		t.Errorf("value is not true")
	}
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
