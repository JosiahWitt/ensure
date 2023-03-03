package entable_test

import (
	"testing"

	"github.com/JosiahWitt/ensure"
	"github.com/JosiahWitt/ensure/ensuring"
	"github.com/JosiahWitt/ensure/exp/entable"
	"github.com/JosiahWitt/ensure/exp/entable/internal/mocks/github.com/JosiahWitt/ensure/mock_ensuring"
	"github.com/golang/mock/gomock"
)

func TestNew(t *testing.T) {
	ensure := ensure.New(t)

	type Entry struct{} //nolint:unused

	ensure(entable.New[Entry]()).Equals(&entable.Table[Entry]{})
}

func TestFrom(t *testing.T) {
	ensure := ensure.New(t)

	type Entry struct{ Name string }

	table := entable.From([]*Entry{
		{Name: "one"},
		{Name: "two"},
	})

	expectedTable := &entable.Table[Entry]{}
	expectedTable.Append(
		&Entry{Name: "one"},
		&Entry{Name: "two"},
	)

	ensure(table).Equals(expectedTable)
}

func TestWithName(t *testing.T) {
	ensure := ensure.New(t)

	type Entry struct{} //nolint:unused

	table := entable.New[Entry]()
	namedTable := table.WithName("hello")

	ensure(table).Equals(entable.New[Entry]())
	ensure(namedTable).Equals(entable.New[Entry]().WithName("hello"))
	ensure(namedTable != table).IsTrue()
}

func TestAppend(t *testing.T) {
	ensure := ensure.New(t)

	type Entry struct{ Name string }

	table := entable.New[Entry]()
	table.Append(
		&Entry{Name: "one"},
		&Entry{Name: "two"},
	)
	table.Append(
		&Entry{Name: "three"},
	)

	ensure(table).Equals(entable.From([]*Entry{
		{Name: "one"},
		{Name: "two"},
		{Name: "three"},
	}))
}

func TestAppendTable(t *testing.T) {
	ensure := ensure.New(t)

	type Entry struct{}

	ensure.Run("when all subtables have names", func(ensure ensuring.E) {
		subtable1 := entable.From([]*Entry{{}}).WithName("sub1")
		subtable2 := entable.From([]*Entry{{}, {}}).WithName("sub2")

		table1 := &entable.Table[Entry]{}
		table1.AppendTable(subtable1, subtable2)

		table2 := &entable.Table[Entry]{}
		table2.AppendTable(subtable1)
		table2.AppendTable(subtable2)

		// Show it doesn't matter if they are appended together or separately
		ensure(table1).Equals(table2)
	})

	ensure.Run("when subtable is missing name", func(ensure ensuring.E) {
		table := &entable.Table[Entry]{}

		subtable1 := entable.From([]*Entry{{}}).WithName("sub1")
		subtable2 := entable.From([]*Entry{{}, {}})

		table.AppendTable(subtable1, subtable2)
	})
}

func TestAppendTableEntries(t *testing.T) {
	ensure := ensure.New(t)

	type Entry struct{ Name string }

	table := entable.New[Entry]()
	table.Append(
		&Entry{Name: "one"},
		&Entry{Name: "two"},
	)

	otherTable := entable.From([]*Entry{
		{Name: "five"},
		{Name: "six"},
	})

	table.AppendTableEntries(
		entable.From([]*Entry{
			{Name: "three"},
			{Name: "four"},
		}),
		otherTable,
	)

	// Show entries appended after AppendTableEntries are not appended to the parent
	otherTable.Append(&Entry{Name: "thrown away"})

	table.Append(
		&Entry{Name: "seven"},
		&Entry{Name: "eight"},
	)

	table.AppendTableEntries(
		entable.From([]*Entry{
			{Name: "nine"},
			{Name: "ten"},
		}),
	)

	ensure(table).Equals(entable.From([]*Entry{
		{Name: "one"},
		{Name: "two"},
		{Name: "three"},
		{Name: "four"},
		{Name: "five"},
		{Name: "six"},
		{Name: "seven"},
		{Name: "eight"},
		{Name: "nine"},
		{Name: "ten"},
	}))
}

func TestIterate(t *testing.T) {
	ensure := ensure.New(t)

	type Entry struct{ Name string }

	table := entable.From([]*Entry{
		{Name: "one"},
		{Name: "two"},
	})

	subtable1 := entable.From([]*Entry{{Name: "three"}}).WithName("sub1")
	subtable2 := entable.From([]*Entry{{Name: "four"}, {Name: "five"}}).WithName("sub2")

	table.AppendTable(subtable1, subtable2)

	table.Append(&Entry{Name: "six"})

	order := []string{}
	table.Iterate(func(entry *Entry) {
		order = append(order, entry.Name)
		entry.Name += " mississippi"
	})

	expectedTable := entable.From([]*Entry{
		{Name: "one mississippi"},
		{Name: "two mississippi"},
		{Name: "six mississippi"},
	})
	expectedTable.AppendTable(
		entable.From([]*Entry{{Name: "three mississippi"}}).WithName("sub1"),
		entable.From([]*Entry{{Name: "four mississippi"}, {Name: "five mississippi"}}).WithName("sub2"),
	)

	ensure(table).Equals(expectedTable)
	ensure(order).Equals([]string{"one", "two", "six", "three", "four", "five"})
}

func TestRun(t *testing.T) {
	ensure := ensure.New(t)

	type Entry struct {
		Name  string
		Value string
	}

	ensure.Run("when table has no entries and no subtables", func(ensure ensuring.E) {
		table := entable.New[Entry]()

		calls := 0
		table.Run(ensure, func(ensure ensuring.E, entry *Entry) {
			calls++
		})

		ensure(calls).Equals(0)
	})

	ensure.Run("when table has entries and no subtables", func(ensure ensuring.E) {
		table := entable.From([]*Entry{
			{
				Name:  "first",
				Value: "hello",
			},
			{
				Name:  "second",
				Value: "world",
			},
		})

		names := []string{}
		calls := []*Entry{}

		table.Run(ensure, func(ensure ensuring.E, entry *Entry) {
			names = append(names, ensure.T().Name())
			calls = append(calls, entry)
		})

		ensure(names).Equals([]string{
			ensure.T().Name() + "/first",
			ensure.T().Name() + "/second",
		})

		ensure(calls).Equals([]*Entry{
			{
				Name:  "first",
				Value: "hello",
			},
			{
				Name:  "second",
				Value: "world",
			},
		})
	})

	ensure.Run("when table has a name prefix and entries but no subtables", func(ensure ensuring.E) {
		table := entable.From([]*Entry{
			{
				Name:  "first",
				Value: "hello",
			},
			{
				Name:  "second",
				Value: "world",
			},
		}).WithName("special")

		names := []string{}
		calls := []*Entry{}

		table.Run(ensure, func(ensure ensuring.E, entry *Entry) {
			names = append(names, ensure.T().Name())
			calls = append(calls, entry)
		})

		ensure(names).Equals([]string{
			ensure.T().Name() + "/special/first",
			ensure.T().Name() + "/special/second",
		})

		ensure(calls).Equals([]*Entry{
			{
				Name:  "first",
				Value: "hello",
			},
			{
				Name:  "second",
				Value: "world",
			},
		})
	})

	ensure.Run("when table has no entries but has subtables", func(ensure ensuring.E) {
		table := entable.New[Entry]()

		subtable1 := entable.From([]*Entry{
			{
				Name:  "first",
				Value: "hello",
			},
			{
				Name:  "second",
				Value: "world",
			},
		}).WithName("sub1")

		subtable2 := entable.From([]*Entry{
			{
				Name:  "third",
				Value: "!",
			},
		}).WithName("sub2")

		table.AppendTable(subtable1, subtable2)

		names := []string{}
		calls := []*Entry{}

		table.Run(ensure, func(ensure ensuring.E, entry *Entry) {
			names = append(names, ensure.T().Name())
			calls = append(calls, entry)
		})

		ensure(names).Equals([]string{
			ensure.T().Name() + "/sub1/first",
			ensure.T().Name() + "/sub1/second",
			ensure.T().Name() + "/sub2/third",
		})

		ensure(calls).Equals([]*Entry{
			{
				Name:  "first",
				Value: "hello",
			},
			{
				Name:  "second",
				Value: "world",
			},
			{
				Name:  "third",
				Value: "!",
			},
		})
	})

	ensure.Run("when table has entries and subtables", func(ensure ensuring.E) {
		table := entable.From([]*Entry{
			{
				Name:  "first",
				Value: "hello",
			},
		})

		table.AppendTable(
			entable.From([]*Entry{
				{
					Name:  "third",
					Value: "!",
				},
			}).WithName("sub1"),
		)

		table.Append(&Entry{
			Name:  "second",
			Value: "world",
		})

		names := []string{}
		calls := []*Entry{}

		table.Run(ensure, func(ensure ensuring.E, entry *Entry) {
			names = append(names, ensure.T().Name())
			calls = append(calls, entry)
		})

		ensure(names).Equals([]string{
			ensure.T().Name() + "/first",
			ensure.T().Name() + "/second",
			ensure.T().Name() + "/sub1/third",
		})

		ensure(calls).Equals([]*Entry{
			{
				Name:  "first",
				Value: "hello",
			},
			{
				Name:  "second",
				Value: "world",
			},
			{
				Name:  "third",
				Value: "!",
			},
		})
	})

	ensure.Run("when table has a name prefix, entries, and subtables", func(ensure ensuring.E) {
		table := entable.From([]*Entry{
			{
				Name:  "first",
				Value: "hello",
			},
		}).WithName("special")

		table.AppendTable(
			entable.From([]*Entry{
				{
					Name:  "third",
					Value: "!",
				},
			}).WithName("sub1"),
		)

		table.Append(&Entry{
			Name:  "second",
			Value: "world",
		})

		table.AppendTable(
			entable.From([]*Entry{
				{
					Name:  "fourth",
					Value: "🌎",
				},
			}).WithName("sub2"),
		)

		names := []string{}
		calls := []*Entry{}

		table.Run(ensure, func(ensure ensuring.E, entry *Entry) {
			names = append(names, ensure.T().Name())
			calls = append(calls, entry)
		})

		ensure(names).Equals([]string{
			ensure.T().Name() + "/special/first",
			ensure.T().Name() + "/special/second",
			ensure.T().Name() + "/special/sub1/third",
			ensure.T().Name() + "/special/sub2/fourth",
		})

		ensure(calls).Equals([]*Entry{
			{
				Name:  "first",
				Value: "hello",
			},
			{
				Name:  "second",
				Value: "world",
			},
			{
				Name:  "third",
				Value: "!",
			},
			{
				Name:  "fourth",
				Value: "🌎",
			},
		})
	})

	ensure.Run("when subtables share names", func(ensure ensuring.E) {
		mockT := mock_ensuring.NewMockT(ensure.GoMockController())
		mockT.EXPECT().Helper().MinTimes(2)
		mockT.EXPECT().Cleanup(gomock.Any()).AnyTimes()
		mockEnsure := ensure.New(mockT)

		type Entry struct{}

		table := entable.New[Entry]()
		table.AppendTable(
			entable.New[Entry]().WithName("sub1"),
			entable.New[Entry]().WithName("sub2"),
			entable.New[Entry]().WithName("sub1"),
		)

		mockT.EXPECT().Fatalf("Subtables names are required to be unique. Subtables[%d] shares a name with Subtables[%d]: %s", 2, 0, "sub1")

		numCalls := 0
		table.Run(mockEnsure, func(ensure ensuring.E, entry *Entry) { numCalls++ })
		ensure(numCalls).Equals(0)
	})

	ensure.Run("when subtable is missing a name", func(ensure ensuring.E) {
		mockT := mock_ensuring.NewMockT(ensure.GoMockController())
		mockT.EXPECT().Helper().MinTimes(2)
		mockT.EXPECT().Cleanup(gomock.Any()).AnyTimes()
		mockEnsure := ensure.New(mockT)

		type Entry struct{}

		table := entable.New[Entry]()
		table.AppendTable(
			entable.New[Entry]().WithName("sub1"),
			entable.New[Entry](),
			entable.New[Entry]().WithName("sub3"),
		)

		mockT.EXPECT().Fatalf(
			"Subtables are required to be named using WithName. Subtables[%d] was not named. "+
				"If you want to append the entries from the table directly instead of naming it, use AppendTableEntries.",
			1,
		)

		numCalls := 0
		table.Run(mockEnsure, func(ensure ensuring.E, entry *Entry) { numCalls++ })
		ensure(numCalls).Equals(0)
	})
}
