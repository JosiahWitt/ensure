// Package plugins contains interfaces and helpers related to internal plugins for ensure.
package plugins

import "reflect"

// TablePlugin is a plugin for table-driven tests run using ensure.
type TablePlugin interface {
	ParseEntryType(entryType reflect.Type) (TableEntryPlugin, error)
}

// TableEntryPlugin is a plugin for entries in table-driven tests run using ensure.
// It is exposed from [TablePlugin].
type TableEntryPlugin interface {
	ParseEntryValue(entryValue reflect.Value, i int) (TableEntryHooks, error)
}

// TableEntryHooks are hooks that run for a particular entry in table-driven tests run using ensure.
// It is exposed from [TableEntryPlugin].
type TableEntryHooks interface {
	BeforeEntry() error
	AfterEntry() error
}
