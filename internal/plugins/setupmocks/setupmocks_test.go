package setupmocks_test

import (
	"reflect"
	"testing"

	"github.com/JosiahWitt/ensure"
	"github.com/JosiahWitt/ensure/ensuring"
	"github.com/JosiahWitt/ensure/internal/mocks/mock_testctx"
	"github.com/JosiahWitt/ensure/internal/plugins/setupmocks"
	"github.com/JosiahWitt/ensure/internal/stringerr"
	"go.uber.org/mock/gomock"
)

func TestParseEntryType(t *testing.T) {
	ensure := ensure.New(t)

	table := []struct {
		Name string

		Entry interface{}

		ExpectedError error
	}{
		{
			Name: "returns no errors when SetupMocks is not provided",

			Entry: struct{ Name string }{},
		},
		{
			Name: "returns no errors when SetupMocks is not provided, but Mocks is provided",

			Entry: struct {
				Name  string
				Mocks *Mocks
			}{},
		},
		{
			Name: "returns no errors when SetupMocks and Mocks are provided",

			Entry: struct {
				Name       string
				Mocks      *Mocks
				SetupMocks func(*Mocks)
			}{},
		},
		{
			Name: "returns no errors when SetupMocks is provided with an optional ensuring.E parameter",

			Entry: struct {
				Name       string
				Mocks      *Mocks
				SetupMocks func(*Mocks, ensuring.E)
			}{},
		},
		{
			Name: "returns error when SetupMocks is provided, but Mocks is not provided",

			Entry: struct {
				Name       string
				SetupMocks func(*Mocks)
			}{},

			ExpectedError: stringerr.Newf("Mocks field must be set on the table to use SetupMocks"),
		},
		{
			Name: "returns error when SetupMocks is not a function",

			Entry: struct {
				Name       string
				Mocks      *Mocks
				SetupMocks *func(*Mocks)
			}{},

			ExpectedError: stringerr.Newf(
				"expected SetupMocks field to be one of the following:\n" +
					" - func(m *setupmocks_test.Mocks)\n" +
					" - func(m *setupmocks_test.Mocks, ensure ensuring.E)\n" +
					"Got: *func(*setupmocks_test.Mocks)",
			),
		},
		{
			Name: "returns error when SetupMocks has no inputs",

			Entry: struct {
				Name       string
				Mocks      *Mocks
				SetupMocks func()
			}{},

			ExpectedError: stringerr.Newf(
				"expected SetupMocks field to be one of the following:\n" +
					" - func(m *setupmocks_test.Mocks)\n" +
					" - func(m *setupmocks_test.Mocks, ensure ensuring.E)\n" +
					"Got: func()",
			),
		},
		{
			Name: "returns error when SetupMocks has two equal inputs",

			Entry: struct {
				Name       string
				Mocks      *Mocks
				SetupMocks func(*Mocks, *Mocks)
			}{},

			ExpectedError: stringerr.Newf(
				"expected SetupMocks field to be one of the following:\n" +
					" - func(m *setupmocks_test.Mocks)\n" +
					" - func(m *setupmocks_test.Mocks, ensure ensuring.E)\n" +
					"Got: func(*setupmocks_test.Mocks, *setupmocks_test.Mocks)",
			),
		},
		{
			Name: "returns error when SetupMocks has an invalid input",

			Entry: struct {
				Name       string
				Mocks      *Mocks
				SetupMocks func(Mocks)
			}{},

			ExpectedError: stringerr.Newf(
				"expected SetupMocks field to be one of the following:\n" +
					" - func(m *setupmocks_test.Mocks)\n" +
					" - func(m *setupmocks_test.Mocks, ensure ensuring.E)\n" +
					"Got: func(setupmocks_test.Mocks)",
			),
		},
		{
			Name: "returns error when SetupMocks has two inputs, and the first is invalid",

			Entry: struct {
				Name       string
				Mocks      *Mocks
				SetupMocks func(Mocks, ensuring.E)
			}{},

			ExpectedError: stringerr.Newf(
				"expected SetupMocks field to be one of the following:\n" +
					" - func(m *setupmocks_test.Mocks)\n" +
					" - func(m *setupmocks_test.Mocks, ensure ensuring.E)\n" +
					"Got: func(setupmocks_test.Mocks, ensuring.E)",
			),
		},
		{
			Name: "returns error when SetupMocks has two inputs, and the second is invalid",

			Entry: struct {
				Name       string
				Mocks      *Mocks
				SetupMocks func(*Mocks, *ensuring.E)
			}{},

			ExpectedError: stringerr.Newf(
				"expected SetupMocks field to be one of the following:\n" +
					" - func(m *setupmocks_test.Mocks)\n" +
					" - func(m *setupmocks_test.Mocks, ensure ensuring.E)\n" +
					"Got: func(*setupmocks_test.Mocks, *ensuring.E)",
			),
		},
		{
			Name: "returns error when SetupMocks has three inputs",

			Entry: struct {
				Name       string
				Mocks      *Mocks
				SetupMocks func(*Mocks, ensuring.E, ensuring.E)
			}{},

			ExpectedError: stringerr.Newf(
				"expected SetupMocks field to be one of the following:\n" +
					" - func(m *setupmocks_test.Mocks)\n" +
					" - func(m *setupmocks_test.Mocks, ensure ensuring.E)\n" +
					"Got: func(*setupmocks_test.Mocks, ensuring.E, ensuring.E)",
			),
		},
		{
			Name: "returns error when SetupMocks returns values",

			Entry: struct {
				Name       string
				Mocks      *Mocks
				SetupMocks func(*Mocks) *Mocks
			}{},

			ExpectedError: stringerr.Newf(
				"expected SetupMocks field to be one of the following:\n" +
					" - func(m *setupmocks_test.Mocks)\n" +
					" - func(m *setupmocks_test.Mocks, ensure ensuring.E)\n" +
					"Got: func(*setupmocks_test.Mocks) *setupmocks_test.Mocks",
			),
		},
	}

	ensure.RunTableByIndex(table, func(ensure ensuring.E, i int) {
		entry := table[i]

		plugin := setupmocks.New()
		res, err := plugin.ParseEntryType(reflect.TypeOf(entry.Entry))
		ensure(err).IsError(entry.ExpectedError)
		ensure(res == nil).Equals(err != nil) // res tested elsewhere
	})
}

func TestParseEntryValue(t *testing.T) {
	ensure := ensure.New(t)

	table := []struct {
		Name string

		Table      interface{}
		SetupMockT func(m *mock_testctx.MockT, i int)

		ExpectedTable interface{}
	}{
		{
			Name: "is a no-op when SetupMocks is not provided",

			Table: []struct{ Name string }{
				{Name: "first"},
				{Name: "second"},
			},

			ExpectedTable: []struct{ Name string }{
				{Name: "first"},
				{Name: "second"},
			},
		},
		{
			Name: "is a no-op when SetupMocks is not provided, but Mocks is provided",

			Table: []struct {
				Name  string
				Mocks *Mocks
			}{
				{Name: "first"},
				{Name: "second"},
			},

			ExpectedTable: []struct {
				Name  string
				Mocks *Mocks
			}{
				{
					Name:  "first",
					Mocks: &Mocks{},
				},
				{
					Name:  "second",
					Mocks: &Mocks{},
				},
			},
		},
		{
			Name: "executes SetupMocks when SetupMocks(*Mocks) and Mocks are provided",

			Table: []struct {
				Name       string
				Mocks      *Mocks
				SetupMocks func(*Mocks)
			}{
				{
					Name: "first",
					SetupMocks: func(m *Mocks) {
						m.A = "first mocks"
					},
				},
				{
					Name: "second",
					SetupMocks: func(m *Mocks) {
						m.A = "second mocks"
					},
				},
			},

			ExpectedTable: []struct {
				Name       string
				Mocks      *Mocks
				SetupMocks func(*Mocks)
			}{
				{
					Name: "first",
					Mocks: &Mocks{
						A: "first mocks",
					},
				},
				{
					Name: "second",
					Mocks: &Mocks{
						A: "second mocks",
					},
				},
			},
		},
		{
			Name: "allows SetupMocks to be missing for some entries",

			Table: []struct {
				Name       string
				Mocks      *Mocks
				SetupMocks func(*Mocks)
			}{
				{
					Name: "first",
					SetupMocks: func(m *Mocks) {
						m.A = "first mocks"
					},
				},
				{
					Name: "second", // No SetupMocks function
				},
				{
					Name: "third",
					SetupMocks: func(m *Mocks) {
						m.A = "third mocks"
					},
				},
			},

			ExpectedTable: []struct {
				Name       string
				Mocks      *Mocks
				SetupMocks func(*Mocks)
			}{
				{
					Name: "first",
					Mocks: &Mocks{
						A: "first mocks",
					},
				},
				{
					Name:  "second",
					Mocks: &Mocks{},
				},
				{
					Name: "third",
					Mocks: &Mocks{
						A: "third mocks",
					},
				},
			},
		},
		{
			Name: "executes SetupMocks when SetupMocks(*Mocks, ensuring.E) and Mocks are provided",

			SetupMockT: func(m *mock_testctx.MockT, i int) {
				switch i {
				case 0:
					m.EXPECT().Fatalf("first fail")
				case 1:
					m.EXPECT().Fatalf("second fail")
				}
			},

			Table: []struct {
				Name       string
				Mocks      *Mocks
				SetupMocks func(*Mocks, ensuring.E)
			}{
				{
					Name: "first",
					SetupMocks: func(m *Mocks, ensure ensuring.E) {
						m.A = "first mocks"
						ensure.Failf("first fail") // Show ensure is connected correctly
					},
				},
				{
					Name: "second",
					SetupMocks: func(m *Mocks, ensure ensuring.E) {
						m.A = "second mocks"
						ensure.Failf("second fail") // Show ensure is connected correctly
					},
				},
			},

			ExpectedTable: []struct {
				Name       string
				Mocks      *Mocks
				SetupMocks func(*Mocks, ensuring.E)
			}{
				{
					Name: "first",
					Mocks: &Mocks{
						A: "first mocks",
					},
				},
				{
					Name: "second",
					Mocks: &Mocks{
						A: "second mocks",
					},
				},
			},
		},
	}

	ensure.RunTableByIndex(table, func(ensure ensuring.E, i int) {
		entry := table[i]

		plugin := setupmocks.New()
		tableEntryHooks, err := plugin.ParseEntryType(reflect.TypeOf(entry.Table).Elem())
		ensure(err).IsNotError()

		tableVal := reflect.ValueOf(entry.Table)
		for i := 0; i < tableVal.Len(); i++ {
			entryVal := tableVal.Index(i)

			if mocksField := entryVal.FieldByName("Mocks"); mocksField.IsValid() {
				mocksField.Set(reflect.New(reflect.TypeOf(Mocks{})))
			}

			mockT := mock_testctx.NewMockT(ensure.GoMockController())
			mockT.EXPECT().Helper().AnyTimes()
			mockT.EXPECT().Cleanup(gomock.Any()).AnyTimes()

			if entry.SetupMockT != nil {
				entry.SetupMockT(mockT, i)
			}

			mockCtx := mock_testctx.NewMockContext(ensure.GoMockController())
			mockCtx.EXPECT().Ensure().Return(ensure.New(mockT)).AnyTimes()

			ensure(tableEntryHooks.BeforeEntry(mockCtx, entryVal, i)).IsNotError()
			ensure(tableEntryHooks.AfterEntry(mockCtx, entryVal, i)).IsNotError()
		}

		ensure(entry.Table).Equals(entry.ExpectedTable)
	})
}

type Mocks struct {
	A string
}
