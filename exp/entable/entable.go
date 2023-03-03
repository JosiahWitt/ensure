// Package entable provides helpers for constructing table driven tests.
// It requires Go 1.18+ due to the use of generics.
//
// Use [Table] to construct a table.
package entable

import (
	"github.com/JosiahWitt/ensure/ensuring"
)

// Table supports constructing and running tables for table-driven testing.
// It also supports nesting other named tables.
//
// For example:
//
//	type Entry struct {
//		Name    string
//		Input   string
//		IsEmpty bool
//	}
//
//	func TestIsEmpty(t *testing.T) {
//		ensure := ensure.New(t)
//
//		table := entable.From([]*Entry{
//			{
//				Name:    "with non empty string",
//				Input:   "my string",
//				IsEmpty: false,
//			},
//			{
//				Name:    "with empty string",
//				Input:   "",
//				IsEmpty: true,
//			},
//		})
//
//		table.AppendTable(isEmptyJSONTests())
//
//		table.Run(ensure, func(ensure ensuring.E, entry *Entry) {
//			isEmpty := strs.IsEmpty(entry.Input)
//			ensure(isEmpty).Equals(entry.IsEmpty)
//		})
//	}
//
//	func isEmptyJSONTests() *entable.Table[Entry] {
//		return entable.From([]*Entry{
//			{
//				Name:    "with empty object",
//				Input:   "{}",
//				IsEmpty: true,
//			},
//			{
//				Name:    "with null",
//				Input:   "null",
//				IsEmpty: true,
//			},
//		}).WithName("JSON")
//	}
type Table[E any] struct {
	name      string
	entries   []*E
	subtables []*Table[E]
}

// New creates an empty [Table].
func New[E any]() *Table[E] {
	return From[E](nil)
}

// From creates a [Table] from the provided slice of table entries.
func From[E any](entries []*E) *Table[E] {
	return &Table[E]{
		entries: entries,
	}
}

// WithName sets a name on the returned copy of [Table]. When the
// Table is run, it executes under a [ensuring.E.Run] block with
// the provided name. WithName creates a shallow copy of the Table,
// so the original Table is unaffected.
func (t Table[E]) WithName(name string) *Table[E] {
	t.name = name
	return &t
}

// Append adds the provided entries to the end of the [Table].
// It modifies the Table in place.
func (t *Table[E]) Append(entries ...*E) {
	t.entries = append(t.entries, entries...)
}

// AppendTable adds the subtables to the end of the [Table]. They
// must be named using [entable.Table.WithName]. Subtables are
// executed after all entries added using [entable.Table.Append]
// or [entable.Table.AppendTableEntries]. AppendTable modifies the
// Table in place.
func (t *Table[E]) AppendTable(subtables ...*Table[E]) {
	t.subtables = append(t.subtables, subtables...)
}

// AppendTableEntries appends the entries of the subtables to the
// end of the [Table]. They are run inline with the other entries
// in the order appended. If any entries are appended to the
// subtables after calling this method, they are NOT picked up.
// AppendTableEntries modifies the Table in place.
func (t *Table[E]) AppendTableEntries(subtables ...*Table[E]) {
	for _, subtable := range subtables {
		t.entries = append(t.entries, subtable.entries...)
	}
}

// Iterate recursively calls fn for every entry in the table and
// any subtables. It iterates through direct entries before proceeding
// to entries of subtables.
func (t *Table[E]) Iterate(fn func(entry *E)) {
	for _, entry := range t.entries {
		fn(entry)
	}

	for _, subTable := range t.subtables {
		subTable.Iterate(fn)
	}
}

// Run runs the Table, executing fn for each entry in the Table and any entries
// in subtables. All direct entries are run before iterating over subtables in
// the order they were appended. See docs for [Table] for an example.
//
// It fails if subtables do not have names or peer subtables share names.
//
// It uses [ensuring.E.RunTableByIndex], so the same functionality is supported,
// including built-in support for mocks.
func (t *Table[E]) Run(ensure ensuring.E, fn func(ensure ensuring.E, entry *E)) {
	ensure.InterfaceT().Helper()

	if t.name != "" {
		ensure.Run(t.name, func(ensure ensuring.E) {
			t.run(ensure, fn)
		})
	} else {
		t.run(ensure, fn)
	}
}

func (t *Table[E]) run(ensure ensuring.E, fn func(ensure ensuring.E, entry *E)) {
	ensure.InterfaceT().Helper()

	subtableNames := make(map[string]int, len(t.subtables))
	for i, subtable := range t.subtables {
		if subtable.name == "" {
			ensure.Failf(
				"Subtables are required to be named using WithName. Subtables[%d] was not named. "+
					"If you want to append the entries from the table directly instead of naming it, use AppendTableEntries.",
				i,
			)
			return
		}

		if lastIdx, ok := subtableNames[subtable.name]; ok {
			ensure.Failf("Subtables names are required to be unique. Subtables[%d] shares a name with Subtables[%d]: %s", i, lastIdx, subtable.name)
			return
		}

		subtableNames[subtable.name] = i
	}

	ensure.RunTableByIndex(t.entries, func(ensure ensuring.E, i int) {
		entry := t.entries[i]
		fn(ensure, entry)
	})

	for _, subtable := range t.subtables {
		subtable.Run(ensure, fn)
	}
}
