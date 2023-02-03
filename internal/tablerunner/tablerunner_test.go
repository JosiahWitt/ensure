package tablerunner_test

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/JosiahWitt/ensure"
	"github.com/JosiahWitt/ensure/ensurepkg"
	"github.com/JosiahWitt/ensure/internal/plugins"
	"github.com/JosiahWitt/ensure/internal/stringerr"
	"github.com/JosiahWitt/ensure/internal/tablerunner"
	"github.com/JosiahWitt/ensure/internal/testctx"
)

func TestBuildTable(t *testing.T) {
	ensure := ensure.New(t)

	successfulPlugins := []plugins.TablePlugin{
		mockTablePlugin(func(entryType reflect.Type) (plugins.TableEntryHooks, error) {
			return nil, nil
		}),
		mockTablePlugin(func(entryType reflect.Type) (plugins.TableEntryHooks, error) {
			return nil, nil
		}),
	}

	failingPlugins := []plugins.TablePlugin{
		mockTablePlugin(func(entryType reflect.Type) (plugins.TableEntryHooks, error) {
			return nil, stringerr.Newf("not good")
		}),
		mockTablePlugin(func(entryType reflect.Type) (plugins.TableEntryHooks, error) {
			return nil, stringerr.Newf("nope")
		}),
	}

	buildPluginErrors := func(msgs ...string) error {
		return errors.New("Errors parsing table:\n - " + strings.Join(msgs, "\n - "))
	}

	table := []*struct {
		Name string

		Table   interface{}
		Plugins []plugins.TablePlugin

		ReturnsBuiltTable bool
		ExpectedError     error
	}{
		{
			Name:              "when provided a valid slice of structs",
			Table:             []struct{ Name string }{{Name: "First"}, {Name: "Second"}},
			ReturnsBuiltTable: true,
		},
		{
			Name:              "when provided a valid array of structs",
			Table:             [2]struct{ Name string }{{Name: "First"}, {Name: "Second"}},
			ReturnsBuiltTable: true,
		},
		{
			Name:              "when provided a valid slice of pointers to structs",
			Table:             []*struct{ Name string }{{Name: "First"}, {Name: "Second"}},
			ReturnsBuiltTable: true,
		},
		{
			Name:              "when provided a valid array of pointers to structs",
			Table:             [2]*struct{ Name string }{{Name: "First"}, {Name: "Second"}},
			ReturnsBuiltTable: true,
		},

		{
			Name:              "when provided a slice of structs with successful plugins",
			Table:             []struct{ Name string }{{Name: "First"}, {Name: "Second"}},
			Plugins:           successfulPlugins,
			ReturnsBuiltTable: true,
		},
		{
			Name:              "when provided a valid array of structs with successful plugins",
			Table:             [2]struct{ Name string }{{Name: "First"}, {Name: "Second"}},
			Plugins:           successfulPlugins,
			ReturnsBuiltTable: true,
		},
		{
			Name:              "when provided a valid slice of pointers to structs with successful plugins",
			Table:             []*struct{ Name string }{{Name: "First"}, {Name: "Second"}},
			Plugins:           successfulPlugins,
			ReturnsBuiltTable: true,
		},
		{
			Name:              "when provided a valid array of pointers to structs with successful plugins",
			Table:             [2]*struct{ Name string }{{Name: "First"}, {Name: "Second"}},
			Plugins:           successfulPlugins,
			ReturnsBuiltTable: true,
		},

		{
			Name:          "when provided table is nil",
			Table:         nil,
			ExpectedError: errors.New("Expected a slice or array for the table, got <nil>"),
		},
		{
			Name:          "when provided table is not an array or slice",
			Table:         "not an array or slice",
			ExpectedError: errors.New("Expected a slice or array for the table, got string"),
		},

		{
			Name:          "when provided a slice of non-struct items",
			Table:         []string{"Hello", "World"},
			ExpectedError: errors.New("Expected entry in table to be a struct or a pointer to a struct, got string"),
		},
		{
			Name:          "when provided an array of non-struct items",
			Table:         [2]string{"Hello", "World"},
			ExpectedError: errors.New("Expected entry in table to be a struct or a pointer to a struct, got string"),
		},

		{
			Name:          "when provided a slice of interfaces",
			Table:         []interface{}{"Hello", "World"},
			ExpectedError: errors.New("Expected entry in table to be a struct or a pointer to a struct, got interface {}"),
		},
		{
			Name:          "when provided an array of interfaces",
			Table:         [2]interface{}{"Hello", "World"},
			ExpectedError: errors.New("Expected entry in table to be a struct or a pointer to a struct, got interface {}"),
		},

		{
			Name:          "when provided a slice of structs with missing Name field",
			Table:         []struct{ NotName string }{{NotName: "nope"}, {NotName: "not it"}},
			ExpectedError: buildPluginErrors("Required Name field does not exist on struct in table"),
		},
		{
			Name:          "when provided an array of structs with missing Name field",
			Table:         [2]struct{ NotName string }{{NotName: "nope"}, {NotName: "not it"}},
			ExpectedError: buildPluginErrors("Required Name field does not exist on struct in table"),
		},
		{
			Name:          "when provided a slice of pointers to structs with missing Name field",
			Table:         []*struct{ NotName string }{{NotName: "nope"}, {NotName: "not it"}},
			ExpectedError: buildPluginErrors("Required Name field does not exist on struct in table"),
		},
		{
			Name:          "when provided an array of pointers to structs with missing Name field",
			Table:         [2]*struct{ NotName string }{{NotName: "nope"}, {NotName: "not it"}},
			ExpectedError: buildPluginErrors("Required Name field does not exist on struct in table"),
		},

		{
			Name:          "when provided a slice of structs with non-string Name field",
			Table:         []struct{ Name int }{{Name: 123}, {Name: 456}},
			ExpectedError: buildPluginErrors("Required Name field in struct in table is not a string"),
		},
		{
			Name:          "when provided an array of structs with non-string Name field",
			Table:         [2]struct{ Name int }{{Name: 123}, {Name: 456}},
			ExpectedError: buildPluginErrors("Required Name field in struct in table is not a string"),
		},
		{
			Name:          "when provided a slice of pointers to structs with non-string Name field",
			Table:         []*struct{ Name int }{{Name: 123}, {Name: 456}},
			ExpectedError: buildPluginErrors("Required Name field in struct in table is not a string"),
		},
		{
			Name:          "when provided an array of pointers to structs with non-string Name field",
			Table:         [2]*struct{ Name int }{{Name: 123}, {Name: 456}},
			ExpectedError: buildPluginErrors("Required Name field in struct in table is not a string"),
		},

		{
			Name:          "when provided a slice of structs with failing plugins",
			Table:         []struct{ Name string }{{Name: "First"}, {Name: "Second"}},
			Plugins:       failingPlugins,
			ExpectedError: buildPluginErrors("not good", "nope"),
		},
		{
			Name:          "when provided a array of structs with failing plugins",
			Table:         [2]struct{ Name string }{{Name: "First"}, {Name: "Second"}},
			Plugins:       failingPlugins,
			ExpectedError: buildPluginErrors("not good", "nope"),
		},
		{
			Name:          "when provided a slice of pointers to structs with failing plugins",
			Table:         []*struct{ Name string }{{Name: "First"}, {Name: "Second"}},
			Plugins:       failingPlugins,
			ExpectedError: buildPluginErrors("not good", "nope"),
		},
		{
			Name:          "when provided a array of pointers to structs with failing plugins",
			Table:         [2]*struct{ Name string }{{Name: "First"}, {Name: "Second"}},
			Plugins:       failingPlugins,
			ExpectedError: buildPluginErrors("not good", "nope"),
		},

		{
			Name:          "when provided a slice of structs with missing names and failing plugins",
			Table:         []struct{}{{}, {}},
			Plugins:       failingPlugins,
			ExpectedError: buildPluginErrors("Required Name field does not exist on struct in table", "not good", "nope"),
		},
		{
			Name:          "when provided a array of structs with missing names and failing plugins",
			Table:         [2]struct{}{{}, {}},
			Plugins:       failingPlugins,
			ExpectedError: buildPluginErrors("Required Name field does not exist on struct in table", "not good", "nope"),
		},
		{
			Name:          "when provided a slice of pointers to structs with missing names and failing plugins",
			Table:         []*struct{}{{}, {}},
			Plugins:       failingPlugins,
			ExpectedError: buildPluginErrors("Required Name field does not exist on struct in table", "not good", "nope"),
		},
		{
			Name:          "when provided a array of pointers to structs with missing names and failing plugins",
			Table:         [2]*struct{}{{}, {}},
			Plugins:       failingPlugins,
			ExpectedError: buildPluginErrors("Required Name field does not exist on struct in table", "not good", "nope"),
		},
	}

	for _, entry := range table {
		ensure.Run(entry.Name, func(ensure ensurepkg.Ensure) {
			builtTable, err := tablerunner.BuildTable(entry.Table, entry.Plugins)
			ensure(err).IsError(entry.ExpectedError)
			ensure(builtTable != nil).Equals(entry.ReturnsBuiltTable)
		})
	}
}

type mockTablePlugin func(entryType reflect.Type) (plugins.TableEntryHooks, error)

func (fn mockTablePlugin) ParseEntryType(entryType reflect.Type) (plugins.TableEntryHooks, error) {
	return fn(entryType)
}

type mockEntryHooks struct {
	before func(testctx.Context, reflect.Value, int) error
	after  func(testctx.Context, reflect.Value, int) error
}

func (m *mockEntryHooks) BeforeEntry(ctx testctx.Context, entryValue reflect.Value, i int) error {
	return m.before(ctx, entryValue, i)
}

func (m *mockEntryHooks) AfterEntry(ctx testctx.Context, entryValue reflect.Value, i int) error {
	return m.after(ctx, entryValue, i)
}
