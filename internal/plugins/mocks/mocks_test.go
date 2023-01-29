package mocks_test

import (
	"reflect"
	"testing"

	"github.com/JosiahWitt/ensure"
	"github.com/JosiahWitt/ensure/ensurepkg"
	"github.com/JosiahWitt/ensure/internal/plugins/internal/mocks"
	"github.com/JosiahWitt/ensure/internal/plugins/internal/testhelper"
	mocksplugin "github.com/JosiahWitt/ensure/internal/plugins/mocks"
	"github.com/JosiahWitt/ensure/internal/stringerr"
	"github.com/JosiahWitt/ensure/internal/testctx"
	"github.com/golang/mock/gomock"
)

func TestParseEntryType(t *testing.T) {
	ensure := ensure.New(t)

	table := []struct {
		Name string

		MocksInput *mocks.All
		Entry      interface{}

		ExpectedMocks *mocks.All
		ExpectedError error
	}{
		{
			Name: "is a no-op when Mocks is not provided",

			MocksInput: &mocks.All{},
			Entry:      struct{ Name string }{},

			ExpectedMocks: &mocks.All{},
		},
		{
			Name: "returns error when Mocks is not a pointer",

			MocksInput: &mocks.All{},
			Entry: struct {
				Name  string
				Mocks struct{}
			}{},

			ExpectedError: stringerr.Newf("expected Mocks field to be a pointer to a struct, got struct {}"),
		},
		{
			Name: "returns error when Mocks is not a pointer to a struct",

			MocksInput: &mocks.All{},
			Entry: struct {
				Name  string
				Mocks *string
			}{},

			ExpectedError: stringerr.Newf("expected Mocks field to be a pointer to a struct, got *string"),
		},
		{
			Name: "is a no-op when Mocks has no fields",

			MocksInput: &mocks.All{},
			Entry: struct {
				Name  string
				Mocks *struct{}
			}{},

			ExpectedMocks: &mocks.All{},
		},
		{
			Name: "is a no-op when Mocks has only unexported fields",

			MocksInput: &mocks.All{},
			Entry: struct {
				Name  string
				Mocks *TableWithUnexportedMock
			}{},

			ExpectedMocks: &mocks.All{},
		},
		{
			Name: "identifies mocks when valid mocks are provided",

			MocksInput: &mocks.All{},
			Entry: struct {
				Name  string
				Mocks *TableWithAllMocks
			}{},

			ExpectedMocks: testhelper.BuildMocks([]*testhelper.MockData{
				{
					Path: "Mocks.M1",
					Mock: &MockNoInsNEW{},
				},
				{
					Path: "Mocks.M2",
					Mock: &MockGoMocksNEW{},
				},
			}),
		},
		{
			Name: "identifies mocks when valid mocks are provided in an embedded struct",

			MocksInput: &mocks.All{},
			Entry: struct {
				Name  string
				Mocks *TableWithEmbeddedMocks
			}{},

			ExpectedMocks: testhelper.BuildMocks([]*testhelper.MockData{
				{
					Path: "Mocks.AllMocks.M1",
					Mock: &MockNoInsNEW{},
				},
				{
					Path: "Mocks.AllMocks.M2",
					Mock: &MockGoMocksNEW{},
				},
				{
					Path: "Mocks.M1",
					Mock: &MockNoInsNEW{},
				},
				{
					Path: "Mocks.M2",
					Mock: &MockGoMocksNEW{},
				},
			}),
		},
		{
			Name: "identifies mocks when valid mocks are provided in a pointer to an embedded struct",

			MocksInput: &mocks.All{},
			Entry: struct {
				Name  string
				Mocks *TableWithEmbeddedPointerMocks
			}{},

			ExpectedMocks: testhelper.BuildMocks([]*testhelper.MockData{
				{
					Path: "Mocks.AllMocks.M1",
					Mock: &MockNoInsNEW{},
				},
				{
					Path: "Mocks.AllMocks.M2",
					Mock: &MockGoMocksNEW{},
				},
				{
					Path: "Mocks.M1",
					Mock: &MockNoInsNEW{},
				},
				{
					Path: "Mocks.M2",
					Mock: &MockGoMocksNEW{},
				},
			}),
		},
		{
			Name: "identifies mocks when some optional mocks are provided",

			MocksInput: &mocks.All{},
			Entry: struct {
				Name  string
				Mocks *TableWithOptionalMocks
			}{},

			ExpectedMocks: testhelper.BuildMocks([]*testhelper.MockData{
				{
					Path: "Mocks.M1",
					Mock: &MockNoInsNEW{},
				},
				{
					Path:     "Mocks.M2",
					Mock:     &MockGoMocksNEW{},
					Optional: true,
				},
				{
					Path: "Mocks.M3",
					Mock: &MockNoInsNEW{},
				},
			}),
		},
		{
			Name: "identifies mocks when some mocks are ignored",

			MocksInput: &mocks.All{},
			Entry: struct {
				Name  string
				Mocks *TableWithIgnoredMocks
			}{},

			ExpectedMocks: testhelper.BuildMocks([]*testhelper.MockData{
				{
					Path: "Mocks.M1",
					Mock: &MockNoInsNEW{},
				},
				{
					Path: "Mocks.M3",
					Mock: &MockGoMocksNEW{},
				},
			}),
		},
		{
			Name: "returns error when tag is empty",

			MocksInput: &mocks.All{},
			Entry: struct {
				Name  string
				Mocks *struct {
					M1 *MockNoInsNEW `ensure:""`
					M2 *MockGoMocksNEW
				}
			}{},

			ExpectedError: stringerr.Newf("Unable to build Mocks field:\n - Mocks.M1: Only `ensure:\"-\"` or `ensure:\"ignoreunused\"` tags are supported, got: `ensure:\"\"`"),

			ExpectedMocks: testhelper.BuildMocks([]*testhelper.MockData{
				{
					Path: "Mocks.M2",
					Mock: &MockGoMocksNEW{},
				},
			}),
		},
		{
			Name: "returns error when tag is invalid",

			MocksInput: &mocks.All{},
			Entry: struct {
				Name  string
				Mocks *struct {
					M1 *MockNoInsNEW `ensure:"ignoreunused "`
					M2 *MockGoMocksNEW
				}
			}{},

			ExpectedError: stringerr.Newf("Unable to build Mocks field:\n - Mocks.M1: Only `ensure:\"-\"` or `ensure:\"ignoreunused\"` tags are supported, got: `ensure:\"ignoreunused \"`"),

			ExpectedMocks: testhelper.BuildMocks([]*testhelper.MockData{
				{
					Path: "Mocks.M2",
					Mock: &MockGoMocksNEW{},
				},
			}),
		},
		{
			Name: "returns error when Mocks has an exported field without a NEW method",

			MocksInput: &mocks.All{},
			Entry: struct {
				Name  string
				Mocks *struct {
					M1 *MockNoInsNEW
					M2 *struct{} // Doesn't have a NEW method
					M3 *MockGoMocksNEW
				}
			}{},

			ExpectedError: stringerr.Newf(
				"Unable to build Mocks field:\n - Mocks.M2 (*struct {}) is missing a NEW method, which is required for automatically initializing the mock. " +
					"To ignore Mocks.M2, add the `ensure:\"-\"` tag.",
			),

			ExpectedMocks: testhelper.BuildMocks([]*testhelper.MockData{
				{
					Path: "Mocks.M1",
					Mock: &MockNoInsNEW{},
				},
				{
					Path: "Mocks.M3",
					Mock: &MockGoMocksNEW{},
				},
			}),
		},
		{
			Name: "returns error when Mocks has an exported field with a NEW method with 2 ins",

			MocksInput: &mocks.All{},
			Entry: struct {
				Name  string
				Mocks *struct {
					M1 *MockNoInsNEW
					M2 *MockInvalidTwoInsNEW
					M3 *MockGoMocksNEW
				}
			}{},

			ExpectedError: stringerr.Newf("Unable to build Mocks field:\n" +
				" - Mocks.M2 (*mocks_test.MockInvalidTwoInsNEW) must have a NEW method matching one of the following signatures:\n" +
				"    - func (m *mocks_test.MockInvalidTwoInsNEW) NEW() *mocks_test.MockInvalidTwoInsNEW { ... }\n" +
				"    - func (m *mocks_test.MockInvalidTwoInsNEW) NEW(ctrl *gomock.Controller) *mocks_test.MockInvalidTwoInsNEW { ... }\n" +
				"   To ignore Mocks.M2, add the `ensure:\"-\"` tag.",
			),

			ExpectedMocks: testhelper.BuildMocks([]*testhelper.MockData{
				{
					Path: "Mocks.M1",
					Mock: &MockNoInsNEW{},
				},
				{
					Path: "Mocks.M3",
					Mock: &MockGoMocksNEW{},
				},
			}),
		},
		{
			Name: "returns error when Mocks has an exported field with a NEW method with invalid single in",

			MocksInput: &mocks.All{},
			Entry: struct {
				Name  string
				Mocks *struct {
					M1 *MockNoInsNEW
					M2 *MockInvalidOneInNEW
					M3 *MockGoMocksNEW
				}
			}{},

			ExpectedError: stringerr.Newf("Unable to build Mocks field:\n" +
				" - Mocks.M2 (*mocks_test.MockInvalidOneInNEW) must have a NEW method matching one of the following signatures:\n" +
				"    - func (m *mocks_test.MockInvalidOneInNEW) NEW() *mocks_test.MockInvalidOneInNEW { ... }\n" +
				"    - func (m *mocks_test.MockInvalidOneInNEW) NEW(ctrl *gomock.Controller) *mocks_test.MockInvalidOneInNEW { ... }\n" +
				"   To ignore Mocks.M2, add the `ensure:\"-\"` tag.",
			),

			ExpectedMocks: testhelper.BuildMocks([]*testhelper.MockData{
				{
					Path: "Mocks.M1",
					Mock: &MockNoInsNEW{},
				},
				{
					Path: "Mocks.M3",
					Mock: &MockGoMocksNEW{},
				},
			}),
		},
		{
			Name: "returns error when Mocks has an exported field with a NEW method with invalid single out",

			MocksInput: &mocks.All{},
			Entry: struct {
				Name  string
				Mocks *struct {
					M1 *MockNoInsNEW
					M2 *MockInvalidOneOutNEW
					M3 *MockGoMocksNEW
				}
			}{},

			ExpectedError: stringerr.Newf("Unable to build Mocks field:\n" +
				" - Mocks.M2 (*mocks_test.MockInvalidOneOutNEW) must have a NEW method matching one of the following signatures:\n" +
				"    - func (m *mocks_test.MockInvalidOneOutNEW) NEW() *mocks_test.MockInvalidOneOutNEW { ... }\n" +
				"    - func (m *mocks_test.MockInvalidOneOutNEW) NEW(ctrl *gomock.Controller) *mocks_test.MockInvalidOneOutNEW { ... }\n" +
				"   To ignore Mocks.M2, add the `ensure:\"-\"` tag.",
			),

			ExpectedMocks: testhelper.BuildMocks([]*testhelper.MockData{
				{
					Path: "Mocks.M1",
					Mock: &MockNoInsNEW{},
				},
				{
					Path: "Mocks.M3",
					Mock: &MockGoMocksNEW{},
				},
			}),
		},
		{
			Name: "returns error when Mocks has an exported field with a NEW method with invalid two outs",

			MocksInput: &mocks.All{},
			Entry: struct {
				Name  string
				Mocks *struct {
					M1 *MockNoInsNEW
					M2 *MockInvalidTwoOutsNEW
					M3 *MockGoMocksNEW
				}
			}{},

			ExpectedError: stringerr.Newf("Unable to build Mocks field:\n" +
				" - Mocks.M2 (*mocks_test.MockInvalidTwoOutsNEW) must have a NEW method matching one of the following signatures:\n" +
				"    - func (m *mocks_test.MockInvalidTwoOutsNEW) NEW() *mocks_test.MockInvalidTwoOutsNEW { ... }\n" +
				"    - func (m *mocks_test.MockInvalidTwoOutsNEW) NEW(ctrl *gomock.Controller) *mocks_test.MockInvalidTwoOutsNEW { ... }\n" +
				"   To ignore Mocks.M2, add the `ensure:\"-\"` tag.",
			),

			ExpectedMocks: testhelper.BuildMocks([]*testhelper.MockData{
				{
					Path: "Mocks.M1",
					Mock: &MockNoInsNEW{},
				},
				{
					Path: "Mocks.M3",
					Mock: &MockGoMocksNEW{},
				},
			}),
		},
		{
			Name: "returns error when Mocks has an exported field with a NEW method with invalid zero outs",

			MocksInput: &mocks.All{},
			Entry: struct {
				Name  string
				Mocks *struct {
					M1 *MockNoInsNEW
					M2 *MockInvalidZeroOutsNEW
					M3 *MockGoMocksNEW
				}
			}{},

			ExpectedError: stringerr.Newf("Unable to build Mocks field:\n" +
				" - Mocks.M2 (*mocks_test.MockInvalidZeroOutsNEW) must have a NEW method matching one of the following signatures:\n" +
				"    - func (m *mocks_test.MockInvalidZeroOutsNEW) NEW() *mocks_test.MockInvalidZeroOutsNEW { ... }\n" +
				"    - func (m *mocks_test.MockInvalidZeroOutsNEW) NEW(ctrl *gomock.Controller) *mocks_test.MockInvalidZeroOutsNEW { ... }\n" +
				"   To ignore Mocks.M2, add the `ensure:\"-\"` tag.",
			),

			ExpectedMocks: testhelper.BuildMocks([]*testhelper.MockData{
				{
					Path: "Mocks.M1",
					Mock: &MockNoInsNEW{},
				},
				{
					Path: "Mocks.M3",
					Mock: &MockGoMocksNEW{},
				},
			}),
		},
	}

	ensure.RunTableByIndex(table, func(ensure ensurepkg.Ensure, i int) {
		entry := table[i]

		plugin := mocksplugin.New(entry.MocksInput)
		res, err := plugin.ParseEntryType(reflect.TypeOf(entry.Entry))
		ensure(err).IsError(entry.ExpectedError)
		ensure(res == nil).Equals(err != nil) // res tested elsewhere

		if entry.ExpectedMocks != nil {
			ensure(entry.MocksInput.Slice()).Equals(entry.ExpectedMocks.Slice())
		} else {
			ensure(entry.MocksInput.Slice()).IsEmpty()
		}
	})
}

func TestParseEntryValue(t *testing.T) {
	ensure := ensure.New(t)

	table := []struct {
		Name string

		MocksInput *mocks.All
		Table      interface{}

		ExpectedTable interface{}
		ExpectedMocks *mocks.All
	}{
		{
			Name: "is a no-op when Mocks is not provided",

			MocksInput: &mocks.All{},
			Table: []struct{ Name string }{
				{Name: "first"},
				{Name: "second"},
			},

			ExpectedTable: []struct{ Name string }{
				{Name: "first"},
				{Name: "second"},
			},

			ExpectedMocks: &mocks.All{},
		},
		{
			Name: "is a no-op when Mocks has no fields",

			MocksInput: &mocks.All{},
			Table: []struct {
				Name  string
				Mocks *struct{}
			}{
				{Name: "first"},
				{Name: "second"},
			},

			ExpectedTable: []struct {
				Name  string
				Mocks *struct{}
			}{
				{Name: "first", Mocks: &struct{}{}},
				{Name: "second", Mocks: &struct{}{}},
			},

			ExpectedMocks: &mocks.All{},
		},
		{
			Name: "is a no-op when Mocks has only unexported fields",

			MocksInput: &mocks.All{},
			Table: []struct {
				Name  string
				Mocks *TableWithUnexportedMock
			}{
				{Name: "first"},
				{Name: "second"},
			},

			ExpectedTable: []struct {
				Name  string
				Mocks *TableWithUnexportedMock
			}{
				{
					Name:  "first",
					Mocks: &TableWithUnexportedMock{},
				},
				{
					Name:  "second",
					Mocks: &TableWithUnexportedMock{},
				},
			},

			ExpectedMocks: &mocks.All{},
		},
		{
			Name: "identifies mocks when valid mocks are provided",

			MocksInput: &mocks.All{},
			Table: []struct {
				Name  string
				Mocks *TableWithAllMocks
			}{
				{Name: "first"},
				{Name: "second"},
			},

			ExpectedTable: []struct {
				Name  string
				Mocks *TableWithAllMocks
			}{
				{
					Name: "first",
					Mocks: &TableWithAllMocks{
						M1: expectedMockNoInsNEW(),
						M2: expectedMockGoMocksNEW(0),
					},
				},
				{
					Name: "second",
					Mocks: &TableWithAllMocks{
						M1: expectedMockNoInsNEW(),
						M2: expectedMockGoMocksNEW(1),
					},
				},
			},

			ExpectedMocks: testhelper.BuildMocks([]*testhelper.MockData{
				{
					Path: "Mocks.M1",
					Mock: &MockNoInsNEW{},
					Values: []interface{}{
						expectedMockNoInsNEW(),
						expectedMockNoInsNEW(),
					},
				},
				{
					Path: "Mocks.M2",
					Mock: &MockGoMocksNEW{},
					Values: []interface{}{
						expectedMockGoMocksNEW(0),
						expectedMockGoMocksNEW(1),
					},
				},
			}),
		},
		{
			Name: "identifies mocks when valid mocks are provided in an embedded struct",

			MocksInput: &mocks.All{},
			Table: []struct {
				Name  string
				Mocks *TableWithEmbeddedMocks
			}{
				{Name: "first"},
				{Name: "second"},
			},

			ExpectedTable: []struct {
				Name  string
				Mocks *TableWithEmbeddedMocks
			}{
				{
					Name: "first",
					Mocks: &TableWithEmbeddedMocks{
						AllMocks: AllMocks{
							M1: expectedMockNoInsNEW(),
							M2: expectedMockGoMocksNEW(0),
						},

						M1: expectedMockNoInsNEW(),
						M2: expectedMockGoMocksNEW(0),
					},
				},
				{
					Name: "second",
					Mocks: &TableWithEmbeddedMocks{
						AllMocks: AllMocks{
							M1: expectedMockNoInsNEW(),
							M2: expectedMockGoMocksNEW(1),
						},

						M1: expectedMockNoInsNEW(),
						M2: expectedMockGoMocksNEW(1),
					},
				},
			},

			ExpectedMocks: testhelper.BuildMocks([]*testhelper.MockData{
				{
					Path: "Mocks.AllMocks.M1",
					Mock: &MockNoInsNEW{},
					Values: []interface{}{
						expectedMockNoInsNEW(),
						expectedMockNoInsNEW(),
					},
				},
				{
					Path: "Mocks.AllMocks.M2",
					Mock: &MockGoMocksNEW{},
					Values: []interface{}{
						expectedMockGoMocksNEW(0),
						expectedMockGoMocksNEW(1),
					},
				},
				{
					Path: "Mocks.M1",
					Mock: &MockNoInsNEW{},
					Values: []interface{}{
						expectedMockNoInsNEW(),
						expectedMockNoInsNEW(),
					},
				},
				{
					Path: "Mocks.M2",
					Mock: &MockGoMocksNEW{},
					Values: []interface{}{
						expectedMockGoMocksNEW(0),
						expectedMockGoMocksNEW(1),
					},
				},
			}),
		},
		{
			Name: "identifies mocks when some mocks are optional",

			MocksInput: &mocks.All{},
			Table: []struct {
				Name  string
				Mocks *TableWithOptionalMocks
			}{
				{Name: "first"},
				{Name: "second"},
			},

			ExpectedTable: []struct {
				Name  string
				Mocks *TableWithOptionalMocks
			}{
				{
					Name: "first",
					Mocks: &TableWithOptionalMocks{
						M1: expectedMockNoInsNEW(),
						M2: expectedMockGoMocksNEW(0),
						M3: expectedMockNoInsNEW(),
					},
				},
				{
					Name: "second",
					Mocks: &TableWithOptionalMocks{
						M1: expectedMockNoInsNEW(),
						M2: expectedMockGoMocksNEW(1),
						M3: expectedMockNoInsNEW(),
					},
				},
			},

			ExpectedMocks: testhelper.BuildMocks([]*testhelper.MockData{
				{
					Path: "Mocks.M1",
					Mock: &MockNoInsNEW{},
					Values: []interface{}{
						expectedMockNoInsNEW(),
						expectedMockNoInsNEW(),
					},
				},
				{
					Path:     "Mocks.M2",
					Mock:     &MockGoMocksNEW{},
					Optional: true,
					Values: []interface{}{
						expectedMockGoMocksNEW(0),
						expectedMockGoMocksNEW(1),
					},
				},
				{
					Path: "Mocks.M3",
					Mock: &MockNoInsNEW{},
					Values: []interface{}{
						expectedMockNoInsNEW(),
						expectedMockNoInsNEW(),
					},
				},
			}),
		},
		{
			Name: "identifies mocks when some mocks are ignored",

			MocksInput: &mocks.All{},
			Table: []struct {
				Name  string
				Mocks *TableWithIgnoredMocks
			}{
				{Name: "first"},
				{Name: "second"},
			},

			ExpectedTable: []struct {
				Name  string
				Mocks *TableWithIgnoredMocks
			}{
				{
					Name: "first",
					Mocks: &TableWithIgnoredMocks{
						M1: expectedMockNoInsNEW(),
						M3: expectedMockGoMocksNEW(0),
					},
				},
				{
					Name: "second",
					Mocks: &TableWithIgnoredMocks{
						M1: expectedMockNoInsNEW(),
						M3: expectedMockGoMocksNEW(1),
					},
				},
			},

			ExpectedMocks: testhelper.BuildMocks([]*testhelper.MockData{
				{
					Path: "Mocks.M1",
					Mock: &MockNoInsNEW{},
					Values: []interface{}{
						expectedMockNoInsNEW(),
						expectedMockNoInsNEW(),
					},
				},
				{
					Path: "Mocks.M3",
					Mock: &MockGoMocksNEW{},
					Values: []interface{}{
						expectedMockGoMocksNEW(0),
						expectedMockGoMocksNEW(1),
					},
				},
			}),
		},
	}

	ensure.RunTableByIndex(table, func(ensure ensurepkg.Ensure, i int) {
		entry := table[i]

		plugin := mocksplugin.New(entry.MocksInput)
		tableEntryPlugin, err := plugin.ParseEntryType(reflect.TypeOf(entry.Table).Elem())
		ensure(err).IsNotError()

		tableVal := reflect.ValueOf(entry.Table)
		for i := 0; i < tableVal.Len(); i++ {
			entryVal := tableVal.Index(i)

			tableEntryHooks, err := tableEntryPlugin.ParseEntryValue(entryVal, i)
			ensure(err).IsNotError()

			ctx := &testctx.Context{
				GoMockController: func() *gomock.Controller { return goMockController(i) },
			}

			ensure(tableEntryHooks.BeforeEntry(ctx)).IsNotError()
			ensure(tableEntryHooks.AfterEntry(ctx)).IsNotError()
		}

		ensure(entry.Table).Equals(entry.ExpectedTable)

		if entry.ExpectedMocks != nil {
			ensure(entry.MocksInput.Slice()).Equals(entry.ExpectedMocks.Slice())
		} else {
			ensure(entry.MocksInput.Slice()).IsEmpty()
		}
	})
}

func goMockController(i int) *gomock.Controller {
	return gomock.NewController(struct {
		gomock.TestReporter
		unique int
	}{unique: i})
}

type (
	//lint:ignore U1000 not used for test purposes
	TableWithUnexportedMock struct {
		unexported *struct{}
	}

	//lint:ignore U1000 not used for test purposes
	TableWithAllMocks struct {
		M1         *MockNoInsNEW
		unexported *struct{}
		M2         *MockGoMocksNEW
	}

	//lint:ignore U1000 not used for test purposes
	TableWithEmbeddedMocks struct {
		AllMocks
		M1         *MockNoInsNEW
		unexported *struct{}
		M2         *MockGoMocksNEW
	}

	//lint:ignore U1000 not used for test purposes
	TableWithEmbeddedPointerMocks struct {
		*AllMocks
		M1         *MockNoInsNEW
		unexported *struct{}
		M2         *MockGoMocksNEW
	}

	//lint:ignore U1000 not used for test purposes
	TableWithOptionalMocks struct {
		M1          *MockNoInsNEW
		notExported *struct{}
		M2          *MockGoMocksNEW `ensure:"ignoreunused"`
		M3          *MockNoInsNEW
	}

	TableWithIgnoredMocks struct {
		M1 *MockNoInsNEW
		M2 *struct{} `ensure:"-"`
		M3 *MockGoMocksNEW
		M4 *MockNoInsNEW `ensure:"-"`
	}
)

type AllMocks struct {
	M1 *MockNoInsNEW
	M2 *MockGoMocksNEW
}

type MockNoInsNEW struct {
	called bool
}

func (m *MockNoInsNEW) NEW() *MockNoInsNEW {
	return &MockNoInsNEW{
		called: true,
	}
}

func expectedMockNoInsNEW() *MockNoInsNEW {
	return &MockNoInsNEW{
		called: true,
	}
}

type MockGoMocksNEW struct {
	called bool
	ctrl   *gomock.Controller
}

func (m *MockGoMocksNEW) NEW(ctrl *gomock.Controller) *MockGoMocksNEW {
	return &MockGoMocksNEW{
		called: true,
		ctrl:   ctrl,
	}
}

func expectedMockGoMocksNEW(i int) *MockGoMocksNEW {
	return &MockGoMocksNEW{
		called: true,
		ctrl:   goMockController(i),
	}
}

type MockInvalidTwoInsNEW struct{}

func (m *MockInvalidTwoInsNEW) NEW(ctrl *gomock.Controller, ctrl2 *gomock.Controller) *MockInvalidTwoInsNEW {
	return &MockInvalidTwoInsNEW{}
}

type MockInvalidOneInNEW struct{}

//nolint:govet // Only used for testing an invalid case.
func (m *MockInvalidOneInNEW) NEW(ctrl gomock.Controller) *MockInvalidOneInNEW {
	return &MockInvalidOneInNEW{}
}

type MockInvalidOneOutNEW struct{}

func (m *MockInvalidOneOutNEW) NEW() MockInvalidOneOutNEW {
	return MockInvalidOneOutNEW{}
}

type MockInvalidTwoOutsNEW struct{}

func (m *MockInvalidTwoOutsNEW) NEW() (*MockInvalidTwoOutsNEW, *MockInvalidTwoOutsNEW) {
	return &MockInvalidTwoOutsNEW{}, &MockInvalidTwoOutsNEW{}
}

type MockInvalidZeroOutsNEW struct{}

func (m *MockInvalidZeroOutsNEW) NEW() {}
