// Package entable provides helpers for constructing table driven tests.
// It requires Go 1.18+ due to the use of generics.
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

// Run runs the table, executing fn for each entry in the table.
//
// It uses [ensuring.E.RunTableByIndex], so the same functionality is supported,
// including built-in support for mocks.
func (t T[E]) Run(ensure ensuring.E, fn func(ensure ensuring.E, entry *E)) {
	ensure.T().Helper()

	t.RunWithIndex(ensure, func(ensure ensuring.E, _ int, entry *E) {
		fn(ensure, entry)
	})
}

// RunWithIndex runs the table, executing fn for each entry in the table.
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
