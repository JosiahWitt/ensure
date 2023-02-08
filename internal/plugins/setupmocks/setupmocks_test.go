package setupmocks_test

import (
	"reflect"
	"testing"

	"github.com/JosiahWitt/ensure"
	"github.com/JosiahWitt/ensure/ensuring"
	"github.com/JosiahWitt/ensure/internal/plugins/setupmocks"
	"github.com/JosiahWitt/ensure/internal/stringerr"
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

			ExpectedError: stringerr.Newf("expected SetupMocks field to be a func(*setupmocks_test.Mocks), got: *func(*setupmocks_test.Mocks)"),
		},
		{
			Name: "returns error when SetupMocks has no inputs",

			Entry: struct {
				Name       string
				Mocks      *Mocks
				SetupMocks func()
			}{},

			ExpectedError: stringerr.Newf("expected SetupMocks field to be a func(*setupmocks_test.Mocks), got: func()"),
		},
		{
			Name: "returns error when SetupMocks has two inputs",

			Entry: struct {
				Name       string
				Mocks      *Mocks
				SetupMocks func(*Mocks, *Mocks)
			}{},

			ExpectedError: stringerr.Newf("expected SetupMocks field to be a func(*setupmocks_test.Mocks), got: func(*setupmocks_test.Mocks, *setupmocks_test.Mocks)"),
		},
		{
			Name: "returns error when SetupMocks has an invalid input",

			Entry: struct {
				Name       string
				Mocks      *Mocks
				SetupMocks func(Mocks)
			}{},

			ExpectedError: stringerr.Newf("expected SetupMocks field to be a func(*setupmocks_test.Mocks), got: func(setupmocks_test.Mocks)"),
		},
		{
			Name: "returns error when SetupMocks returns values",

			Entry: struct {
				Name       string
				Mocks      *Mocks
				SetupMocks func(*Mocks) *Mocks
			}{},

			ExpectedError: stringerr.Newf("expected SetupMocks field to be a func(*setupmocks_test.Mocks), got: func(*setupmocks_test.Mocks) *setupmocks_test.Mocks"),
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

		Table interface{}

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
			Name: "executes SetupMocks when SetupMocks and Mocks are provided",

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

			ensure(tableEntryHooks.BeforeEntry(nil, entryVal, i)).IsNotError()
			ensure(tableEntryHooks.AfterEntry(nil, entryVal, i)).IsNotError()
		}

		ensure(entry.Table).Equals(entry.ExpectedTable)
	})
}

type Mocks struct {
	A string
}
