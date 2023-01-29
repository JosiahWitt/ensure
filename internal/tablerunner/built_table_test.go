package tablerunner_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/JosiahWitt/ensure"
	"github.com/JosiahWitt/ensure/ensurepkg"
	"github.com/JosiahWitt/ensure/internal/plugins"
	"github.com/JosiahWitt/ensure/internal/tablerunner"
	"github.com/JosiahWitt/ensure/internal/testctx"
	"github.com/golang/mock/gomock"
)

type ExampleEntry struct {
	Name string
}

func (e ExampleEntry) ptr() *ExampleEntry { return &e }

type RunEntry struct {
	Name string

	Table   []ExampleEntry
	Plugins []plugins.TablePlugin

	ExpectedNames  []string
	ExpectedFatals map[int]string
	ExpectedRuns   []int
	ExpectedState  []string
}

func TestBuiltTableRun(t *testing.T) {
	ensure := ensure.New(t)

	var state []string

	type hookFunc func(i int, name string) error

	noopHook := func(i int, name string) error { return nil }

	buildPlugins := func(plugin1Before, plugin1After, plugin2Before, plugin2After hookFunc) []plugins.TablePlugin {
		return []plugins.TablePlugin{
			mockTablePlugin(func(entryType reflect.Type) (plugins.TableEntryPlugin, error) {
				return mockEntryPlugin(func(entryValue reflect.Value, i int) (plugins.TableEntryHooks, error) {
					name := entryValue.FieldByName("Name").String()
					return &mockEntryHooks{
						before: func(t *testctx.Context) error {
							assertTestContext(t, i)
							state = append(state, fmt.Sprintf("plugin1_before_%d_%s", i, name))
							return plugin1Before(i, name)
						},
						after: func(t *testctx.Context) error {
							assertTestContext(t, i)
							state = append(state, fmt.Sprintf("plugin1_after_%d_%s", i, name))
							return plugin1After(i, name)
						},
					}, nil
				}), nil
			}),
			mockTablePlugin(func(entryType reflect.Type) (plugins.TableEntryPlugin, error) {
				return mockEntryPlugin(func(entryValue reflect.Value, i int) (plugins.TableEntryHooks, error) {
					name := entryValue.FieldByName("Name").String()
					return &mockEntryHooks{
						before: func(t *testctx.Context) error {
							assertTestContext(t, i)
							state = append(state, fmt.Sprintf("plugin2_before_%d_%s", i, name))
							return plugin2Before(i, name)
						},
						after: func(t *testctx.Context) error {
							assertTestContext(t, i)
							state = append(state, fmt.Sprintf("plugin2_after_%d_%s", i, name))
							return plugin2After(i, name)
						},
					}, nil
				}), nil
			}),
		}
	}

	table := []*RunEntry{
		{
			Name:  "when running table is successful with no plugins",
			Table: []ExampleEntry{{Name: "First"}, {Name: "Middle"}, {Name: "Last"}},

			ExpectedNames:  []string{"First", "Middle", "Last"},
			ExpectedFatals: map[int]string{},
			ExpectedRuns:   []int{0, 1, 2},
		},
		{
			Name:  "when running table is successful with plugins",
			Table: []ExampleEntry{{Name: "First"}, {Name: "Middle"}, {Name: "Last"}},

			Plugins: buildPlugins(
				noopHook,
				noopHook,

				noopHook,
				noopHook,
			),

			ExpectedNames:  []string{"First", "Middle", "Last"},
			ExpectedFatals: map[int]string{},
			ExpectedRuns:   []int{0, 1, 2},

			ExpectedState: []string{
				"plugin1_before_0_First", "plugin2_before_0_First",
				"plugin1_after_0_First", "plugin2_after_0_First",

				"plugin1_before_1_Middle", "plugin2_before_1_Middle",
				"plugin1_after_1_Middle", "plugin2_after_1_Middle",

				"plugin1_before_2_Last", "plugin2_before_2_Last",
				"plugin1_after_2_Last", "plugin2_after_2_Last",
			},
		},

		{
			Name:  "when name is missing for a field",
			Table: []ExampleEntry{{Name: "First"}, {Name: ""}, {Name: "Last"}},

			ExpectedNames:  []string{"First", "", "Last"},
			ExpectedFatals: map[int]string{1: "Errors parsing table entry:\n - table[1].Name is empty"},
			ExpectedRuns:   []int{0, 2},
		},
		{
			Name:  "when name is duplicated for a field",
			Table: []ExampleEntry{{Name: "First"}, {Name: "First"}, {Name: "Last"}},

			ExpectedNames:  []string{"First", "First", "Last"},
			ExpectedFatals: map[int]string{1: "Errors parsing table entry:\n - table[1].Name duplicates table[0].Name: First"},
			ExpectedRuns:   []int{0, 2},
		},

		// ***** Before Hook Failure Tests ***** //
		{
			Name:  "when before hook in one plugin fails for one item",
			Table: []ExampleEntry{{Name: "First"}, {Name: "Middle"}, {Name: "Last"}},

			Plugins: buildPlugins(
				func(i int, name string) error {
					if i == 1 {
						return fmt.Errorf("plugin1_before_failed_%d_%s", i, name)
					}

					return nil
				},
				noopHook,

				noopHook,
				noopHook,
			),

			ExpectedNames: []string{"First", "Middle", "Last"},
			ExpectedFatals: map[int]string{
				1: "Errors running plugins:\n - plugin1_before_failed_1_Middle",
			},
			ExpectedRuns: []int{0, 2},

			ExpectedState: []string{
				"plugin1_before_0_First", "plugin2_before_0_First",
				"plugin1_after_0_First", "plugin2_after_0_First",

				"plugin1_before_1_Middle", "plugin2_before_1_Middle",

				"plugin1_before_2_Last", "plugin2_before_2_Last",
				"plugin1_after_2_Last", "plugin2_after_2_Last",
			},
		},
		{
			Name:  "when before hook in multiple plugins fails for one item",
			Table: []ExampleEntry{{Name: "First"}, {Name: "Middle"}, {Name: "Last"}},

			Plugins: buildPlugins(
				func(i int, name string) error {
					if i == 1 {
						return fmt.Errorf("plugin1_before_failed_%d_%s", i, name)
					}

					return nil
				},
				noopHook,

				func(i int, name string) error {
					if i == 1 {
						return fmt.Errorf("plugin2_before_failed_%d_%s", i, name)
					}

					return nil
				},
				noopHook,
			),

			ExpectedNames: []string{"First", "Middle", "Last"},
			ExpectedFatals: map[int]string{
				1: "Errors running plugins:\n - plugin1_before_failed_1_Middle\n - plugin2_before_failed_1_Middle",
			},
			ExpectedRuns: []int{0, 2},

			ExpectedState: []string{
				"plugin1_before_0_First", "plugin2_before_0_First",
				"plugin1_after_0_First", "plugin2_after_0_First",

				"plugin1_before_1_Middle", "plugin2_before_1_Middle",

				"plugin1_before_2_Last", "plugin2_before_2_Last",
				"plugin1_after_2_Last", "plugin2_after_2_Last",
			},
		},
		{
			Name:  "when before hook in one plugin fails for all items",
			Table: []ExampleEntry{{Name: "First"}, {Name: "Middle"}, {Name: "Last"}},

			Plugins: buildPlugins(
				func(i int, name string) error {
					return fmt.Errorf("plugin1_before_failed_%d_%s", i, name)
				},
				noopHook,

				noopHook,
				noopHook,
			),

			ExpectedNames: []string{"First", "Middle", "Last"},
			ExpectedFatals: map[int]string{
				0: "Errors running plugins:\n - plugin1_before_failed_0_First",
				1: "Errors running plugins:\n - plugin1_before_failed_1_Middle",
				2: "Errors running plugins:\n - plugin1_before_failed_2_Last",
			},
			ExpectedRuns: []int{},

			ExpectedState: []string{
				"plugin1_before_0_First", "plugin2_before_0_First",

				"plugin1_before_1_Middle", "plugin2_before_1_Middle",

				"plugin1_before_2_Last", "plugin2_before_2_Last",
			},
		},
		{
			Name:  "when before hook in multiple plugins fails for all items",
			Table: []ExampleEntry{{Name: "First"}, {Name: "Middle"}, {Name: "Last"}},

			Plugins: buildPlugins(
				func(i int, name string) error {
					return fmt.Errorf("plugin1_before_failed_%d_%s", i, name)
				},
				noopHook,

				func(i int, name string) error {
					return fmt.Errorf("plugin2_before_failed_%d_%s", i, name)
				},
				noopHook,
			),

			ExpectedNames: []string{"First", "Middle", "Last"},
			ExpectedFatals: map[int]string{
				0: "Errors running plugins:\n - plugin1_before_failed_0_First\n - plugin2_before_failed_0_First",
				1: "Errors running plugins:\n - plugin1_before_failed_1_Middle\n - plugin2_before_failed_1_Middle",
				2: "Errors running plugins:\n - plugin1_before_failed_2_Last\n - plugin2_before_failed_2_Last",
			},
			ExpectedRuns: []int{},

			ExpectedState: []string{
				"plugin1_before_0_First", "plugin2_before_0_First",

				"plugin1_before_1_Middle", "plugin2_before_1_Middle",

				"plugin1_before_2_Last", "plugin2_before_2_Last",
			},
		},
		{
			Name:  "when before hook in multiple plugins fails for different items",
			Table: []ExampleEntry{{Name: "First"}, {Name: "Middle"}, {Name: "Last"}},

			Plugins: buildPlugins(
				func(i int, name string) error {
					if i == 1 {
						return fmt.Errorf("plugin1_before_failed_%d_%s", i, name)
					}

					return nil
				},
				noopHook,

				func(i int, name string) error {
					if i == 2 {
						return fmt.Errorf("plugin2_before_failed_%d_%s", i, name)
					}

					return nil
				},
				noopHook,
			),

			ExpectedNames: []string{"First", "Middle", "Last"},
			ExpectedFatals: map[int]string{
				1: "Errors running plugins:\n - plugin1_before_failed_1_Middle",
				2: "Errors running plugins:\n - plugin2_before_failed_2_Last",
			},
			ExpectedRuns: []int{0},

			ExpectedState: []string{
				"plugin1_before_0_First", "plugin2_before_0_First",
				"plugin1_after_0_First", "plugin2_after_0_First",

				"plugin1_before_1_Middle", "plugin2_before_1_Middle",

				"plugin1_before_2_Last", "plugin2_before_2_Last",
			},
		},

		// ***** After Hook Failure Tests ***** //
		{
			Name:  "when after hook in one plugin fails for one item",
			Table: []ExampleEntry{{Name: "First"}, {Name: "Middle"}, {Name: "Last"}},

			Plugins: buildPlugins(
				noopHook,
				func(i int, name string) error {
					if i == 1 {
						return fmt.Errorf("plugin1_after_failed_%d_%s", i, name)
					}

					return nil
				},

				noopHook,
				noopHook,
			),

			ExpectedNames: []string{"First", "Middle", "Last"},
			ExpectedFatals: map[int]string{
				1: "Errors running plugins:\n - plugin1_after_failed_1_Middle",
			},
			ExpectedRuns: []int{0, 1, 2},

			ExpectedState: []string{
				"plugin1_before_0_First", "plugin2_before_0_First",
				"plugin1_after_0_First", "plugin2_after_0_First",

				"plugin1_before_1_Middle", "plugin2_before_1_Middle",
				"plugin1_after_1_Middle", "plugin2_after_1_Middle",

				"plugin1_before_2_Last", "plugin2_before_2_Last",
				"plugin1_after_2_Last", "plugin2_after_2_Last",
			},
		},
		{
			Name:  "when after hook in multiple plugins fails for one item",
			Table: []ExampleEntry{{Name: "First"}, {Name: "Middle"}, {Name: "Last"}},

			Plugins: buildPlugins(
				noopHook,
				func(i int, name string) error {
					if i == 1 {
						return fmt.Errorf("plugin1_after_failed_%d_%s", i, name)
					}

					return nil
				},

				noopHook,
				func(i int, name string) error {
					if i == 1 {
						return fmt.Errorf("plugin2_after_failed_%d_%s", i, name)
					}

					return nil
				},
			),

			ExpectedNames: []string{"First", "Middle", "Last"},
			ExpectedFatals: map[int]string{
				1: "Errors running plugins:\n - plugin1_after_failed_1_Middle\n - plugin2_after_failed_1_Middle",
			},
			ExpectedRuns: []int{0, 1, 2},

			ExpectedState: []string{
				"plugin1_before_0_First", "plugin2_before_0_First",
				"plugin1_after_0_First", "plugin2_after_0_First",

				"plugin1_before_1_Middle", "plugin2_before_1_Middle",
				"plugin1_after_1_Middle", "plugin2_after_1_Middle",

				"plugin1_before_2_Last", "plugin2_before_2_Last",
				"plugin1_after_2_Last", "plugin2_after_2_Last",
			},
		},
		{
			Name:  "when after hook in one plugin fails for all items",
			Table: []ExampleEntry{{Name: "First"}, {Name: "Middle"}, {Name: "Last"}},

			Plugins: buildPlugins(
				noopHook,
				func(i int, name string) error {
					return fmt.Errorf("plugin1_after_failed_%d_%s", i, name)
				},

				noopHook,
				noopHook,
			),

			ExpectedNames: []string{"First", "Middle", "Last"},
			ExpectedFatals: map[int]string{
				0: "Errors running plugins:\n - plugin1_after_failed_0_First",
				1: "Errors running plugins:\n - plugin1_after_failed_1_Middle",
				2: "Errors running plugins:\n - plugin1_after_failed_2_Last",
			},
			ExpectedRuns: []int{0, 1, 2},

			ExpectedState: []string{
				"plugin1_before_0_First", "plugin2_before_0_First",
				"plugin1_after_0_First", "plugin2_after_0_First",

				"plugin1_before_1_Middle", "plugin2_before_1_Middle",
				"plugin1_after_1_Middle", "plugin2_after_1_Middle",

				"plugin1_before_2_Last", "plugin2_before_2_Last",
				"plugin1_after_2_Last", "plugin2_after_2_Last",
			},
		},
		{
			Name:  "when after hook in multiple plugins fails for all items",
			Table: []ExampleEntry{{Name: "First"}, {Name: "Middle"}, {Name: "Last"}},

			Plugins: buildPlugins(
				noopHook,
				func(i int, name string) error {
					return fmt.Errorf("plugin1_after_failed_%d_%s", i, name)
				},

				noopHook,
				func(i int, name string) error {
					return fmt.Errorf("plugin2_after_failed_%d_%s", i, name)
				},
			),

			ExpectedNames: []string{"First", "Middle", "Last"},
			ExpectedFatals: map[int]string{
				0: "Errors running plugins:\n - plugin1_after_failed_0_First\n - plugin2_after_failed_0_First",
				1: "Errors running plugins:\n - plugin1_after_failed_1_Middle\n - plugin2_after_failed_1_Middle",
				2: "Errors running plugins:\n - plugin1_after_failed_2_Last\n - plugin2_after_failed_2_Last",
			},
			ExpectedRuns: []int{0, 1, 2},

			ExpectedState: []string{
				"plugin1_before_0_First", "plugin2_before_0_First",
				"plugin1_after_0_First", "plugin2_after_0_First",

				"plugin1_before_1_Middle", "plugin2_before_1_Middle",
				"plugin1_after_1_Middle", "plugin2_after_1_Middle",

				"plugin1_before_2_Last", "plugin2_before_2_Last",
				"plugin1_after_2_Last", "plugin2_after_2_Last",
			},
		},
		{
			Name:  "when after hook in multiple plugins fails for different items",
			Table: []ExampleEntry{{Name: "First"}, {Name: "Middle"}, {Name: "Last"}},

			Plugins: buildPlugins(
				noopHook,
				func(i int, name string) error {
					if i == 1 {
						return fmt.Errorf("plugin1_after_failed_%d_%s", i, name)
					}

					return nil
				},

				noopHook,
				func(i int, name string) error {
					if i == 2 {
						return fmt.Errorf("plugin2_after_failed_%d_%s", i, name)
					}

					return nil
				},
			),

			ExpectedNames: []string{"First", "Middle", "Last"},
			ExpectedFatals: map[int]string{
				1: "Errors running plugins:\n - plugin1_after_failed_1_Middle",
				2: "Errors running plugins:\n - plugin2_after_failed_2_Last",
			},
			ExpectedRuns: []int{0, 1, 2},

			ExpectedState: []string{
				"plugin1_before_0_First", "plugin2_before_0_First",
				"plugin1_after_0_First", "plugin2_after_0_First",

				"plugin1_before_1_Middle", "plugin2_before_1_Middle",
				"plugin1_after_1_Middle", "plugin2_after_1_Middle",

				"plugin1_before_2_Last", "plugin2_before_2_Last",
				"plugin1_after_2_Last", "plugin2_after_2_Last",
			},
		},

		// ***** Before and After Hook Failure Tests ***** //
		{
			Name:  "when after and before hooks in multiple plugins fail for different items",
			Table: []ExampleEntry{{Name: "First"}, {Name: "Middle"}, {Name: "Last"}},

			Plugins: buildPlugins(
				func(i int, name string) error {
					if i == 1 {
						return fmt.Errorf("plugin1_before_failed_%d_%s", i, name)
					}

					return nil
				},
				func(i int, name string) error {
					if i == 0 {
						return fmt.Errorf("plugin1_after_failed_%d_%s", i, name)
					}

					return nil
				},

				noopHook,
				func(i int, name string) error {
					if i == 2 {
						return fmt.Errorf("plugin2_after_failed_%d_%s", i, name)
					}

					return nil
				},
			),

			ExpectedNames: []string{"First", "Middle", "Last"},
			ExpectedFatals: map[int]string{
				0: "Errors running plugins:\n - plugin1_after_failed_0_First",
				1: "Errors running plugins:\n - plugin1_before_failed_1_Middle",
				2: "Errors running plugins:\n - plugin2_after_failed_2_Last",
			},
			ExpectedRuns: []int{0, 2},

			ExpectedState: []string{
				"plugin1_before_0_First", "plugin2_before_0_First",
				"plugin1_after_0_First", "plugin2_after_0_First",

				"plugin1_before_1_Middle", "plugin2_before_1_Middle",

				"plugin1_before_2_Last", "plugin2_before_2_Last",
				"plugin1_after_2_Last", "plugin2_after_2_Last",
			},
		},
	}

	ensure.Run("when table is a slice without pointers", func(ensure ensurepkg.Ensure) {
		for _, entry := range table {
			ensure.Run(entry.Name, func(ensure ensurepkg.Ensure) {
				entry.runTable(ensure, &state, entry.Table)
			})
		}
	})

	ensure.Run("when table is a slice with pointers", func(ensure ensurepkg.Ensure) {
		for _, entry := range table {
			ensure.Run(entry.Name, func(ensure ensurepkg.Ensure) {
				pointerTable := []*ExampleEntry{}

				for _, entry := range entry.Table {
					pointerTable = append(pointerTable, entry.ptr())
				}

				entry.runTable(ensure, &state, pointerTable)
			})
		}
	})

	ensure.Run("when table is an array without pointers", func(ensure ensurepkg.Ensure) {
		for _, entry := range table {
			ensure.Run(entry.Name, func(ensure ensurepkg.Ensure) {
				arrayTable := [3]ExampleEntry{}
				copy(arrayTable[:], entry.Table)
				entry.runTable(ensure, &state, arrayTable)
			})
		}
	})

	ensure.Run("when table is an array with pointers", func(ensure ensurepkg.Ensure) {
		for _, entry := range table {
			ensure.Run(entry.Name, func(ensure ensurepkg.Ensure) {
				pointerTable := [3]*ExampleEntry{}

				for i, entry := range entry.Table {
					pointerTable[i] = entry.ptr()
				}

				entry.runTable(ensure, &state, pointerTable)
			})
		}
	})
}

func (entry *RunEntry) runTable(ensure ensurepkg.Ensure, state *[]string, table interface{}) {
	*state = []string{}

	builtTable, err := tablerunner.BuildTable(table, entry.Plugins)
	ensure(err).IsNotError()

	outerT := &mockT{unique: -1}
	names := []string{}
	fatals := map[int]string{}
	runs := []int{}
	i := 0

	builtTable.Run(outerT, func(name string, callback func(*testctx.Context, func(int))) {
		names = append(names, name)

		ctx, innerT := buildTestContext(i)
		callback(ctx, func(i int) {
			runs = append(runs, i)
		})

		ensure(innerT.helperCount).Equals(1)

		if innerT.fatalMessage != "" {
			fatals[i] = innerT.fatalMessage
		}

		i++
	})

	ensure(outerT.helperCount).Equals(1)
	ensure(outerT.fatalMessage).IsEmpty()

	ensure(names).Equals(entry.ExpectedNames)
	ensure(runs).Equals(entry.ExpectedRuns)
	ensure(fatals).Equals(entry.ExpectedFatals)

	if entry.ExpectedState != nil {
		ensure(*state).Equals(entry.ExpectedState)
	} else {
		ensure(*state).IsEmpty()
	}
}

type mockT struct {
	helperCount  int
	fatalMessage string
	unique       int
}

func (t *mockT) Helper() { t.helperCount++ }

func (t *mockT) Fatalf(format string, args ...interface{}) {
	if t.fatalMessage != "" {
		panic("double call to Fatalf")
	}

	t.fatalMessage = fmt.Sprintf(format, args...)
}

func buildTestContext(i int) (*testctx.Context, *mockT) {
	t := &mockT{unique: i}

	ctx := &testctx.Context{
		T: t,
		GoMockController: func() *gomock.Controller {
			return gomock.NewController(&goMockTestHelper{unique: i})
		},
	}

	return ctx, t
}

type goMockTestHelper struct {
	gomock.TestHelper
	unique int
}

func assertTestContext(actualCtx *testctx.Context, i int) {
	expected := []int{i, i}

	actualMockT := actualCtx.T.(*mockT)
	actualGoMockTestHelper := actualCtx.GoMockController().T.(*goMockTestHelper)
	actual := []int{actualMockT.unique, actualGoMockTestHelper.unique}

	if !reflect.DeepEqual(actual, expected) {
		panic(fmt.Sprintf("testctx.Context does not equal the expected (GOT %v, EXPECTED %v)", actual, expected))
	}
}
