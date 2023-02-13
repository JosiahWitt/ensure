// Package entable provides helpers for constructing table driven tests.
// It requires Go 1.18+ due to the use of generics.
//
// Use [T] to construct a table. Use [G] to nest a table in a group.
package entable

import "github.com/JosiahWitt/ensure/ensuring"

// T allows constructing tables and running them.
//
// For example:
//
//	type Entry struct {
//	  Name    string
//	  Input   string
//	  IsEmpty bool
//	}
//
//	table := entable.T[Entry]{
//	  {
//	    Name:    "with non empty input",
//	    Input:   "my string",
//	    IsEmpty: false,
//	  },
//	  {
//	    Name:    "with empty input",
//	    Input:   "",
//	    IsEmpty: true,
//	  },
//	}
//
//	table.Run(ensure, func(ensure Ensure, entry *Entry) {
//	  isEmpty := strs.IsEmpty(entry.Input)
//	  ensure(isEmpty).Equals(entry.IsEmpty)
//	})
//
//	// Or
//
//	table.RunWithIndex(ensure, func(ensure Ensure, i int, entry *Entry) {
//	  isEmpty := strs.IsEmpty(entry.Input)
//	  ensure(isEmpty).Equals(entry.IsEmpty)
//	})
type T[E any] []*E

// Run runs the table, executing fn for each entry in the table. See docs for [T]
// for an example.
//
// It uses [ensuring.E.RunTableByIndex], so the same functionality is supported,
// including built-in support for mocks.
func (t T[E]) Run(ensure ensuring.E, fn func(ensure ensuring.E, entry *E)) {
	ensure.T().Helper()

	t.RunWithIndex(ensure, func(ensure ensuring.E, _ int, entry *E) {
		fn(ensure, entry)
	})
}

// RunWithIndex runs the table, executing fn for each entry in the table. See
// docs for [T] for an example.
//
// It uses [ensuring.E.RunTableByIndex], so the same functionality is supported,
// including built-in support for mocks.
func (t T[E]) RunWithIndex(ensure ensuring.E, fn func(ensure ensuring.E, i int, entry *E)) {
	ensure.T().Helper()

	ensure.RunTableByIndex(t, func(ensure ensuring.E, i int) {
		entry := t[i]
		fn(ensure, i, entry)
	})
}

// G allows constructing a table within a group. It executes all entries
// in a separate test scope with the provided name.
//
// For example:
//
//	type Entry struct {
//		Name    string
//		Input   string
//		IsEmpty bool
//	}
//
//	group := entable.G[Entry]{
//		Name: "grp",
//		Table: entable.T[Entry]{
//			{
//				Name:    "with non empty input",
//				Input:   "my string",
//				IsEmpty: false,
//			},
//			{
//				Name:    "with empty input",
//				Input:   "",
//				IsEmpty: true,
//			},
//		},
//	}
//
//	group.Run(ensure, func(ensure Ensure, entry *Entry) {
//		// ensure.T().Name() ends with `/grp/<entry name here>`
//		isEmpty := strs.IsEmpty(entry.Input)
//		ensure(isEmpty).Equals(entry.IsEmpty)
//	})
//
//	// Or
//
//	group.RunWithIndex(ensure, func(ensure Ensure, i int, entry *Entry) {
//		// ensure.T().Name() ends with `/grp/<entry name here>`
//		isEmpty := strs.IsEmpty(entry.Input)
//		ensure(isEmpty).Equals(entry.IsEmpty)
//	})
type G[E any] struct {
	Name  string
	Table T[E]
}

// Run runs the table, executing fn for each entry in the table. See docs for [G]
// for an example.
//
// It uses [ensuring.E.RunTableByIndex], so the same functionality is supported,
// including built-in support for mocks.
func (g *G[E]) Run(ensure ensuring.E, fn func(ensure ensuring.E, entry *E)) {
	ensure.T().Helper()

	g.RunWithIndex(ensure, func(ensure ensuring.E, _ int, entry *E) {
		fn(ensure, entry)
	})
}

// RunWithIndex runs the table, executing fn for each entry in the table. See
// docs for [G] for an example.
//
// It uses [ensuring.E.RunTableByIndex], so the same functionality is supported,
// including built-in support for mocks.
func (g *G[E]) RunWithIndex(ensure ensuring.E, fn func(ensure ensuring.E, i int, entry *E)) {
	ensure.T().Helper()

	ensure.Run(g.Name, func(ensure ensuring.E) {
		ensure.RunTableByIndex(g.Table, func(ensure ensuring.E, i int) {
			entry := g.Table[i]
			fn(ensure, i, entry)
		})
	})
}
