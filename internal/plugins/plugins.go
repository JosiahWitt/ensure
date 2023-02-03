// Package plugins contains interfaces and helpers related to internal plugins for ensure.
package plugins

import (
	"reflect"

	"github.com/JosiahWitt/ensure/internal/testctx"
)

// TablePlugin is a plugin for table-driven tests run using ensure.
type TablePlugin interface {
	ParseEntryType(entryType reflect.Type) (TableEntryHooks, error)
}

// TableEntryHooks are hooks that run for a particular entry in table-driven tests run using ensure.
// It is exposed from [TablePlugin].
type TableEntryHooks interface {
	BeforeEntry(ctx testctx.Context, entryValue reflect.Value, i int) error
	AfterEntry(ctx testctx.Context, entryValue reflect.Value, i int) error
}

// NoopAfterEntry exposes [AfterEntry] which is a no-op.
type NoopAfterEntry struct{}

// AfterEntry is called after the test is run for the table entry.
func (NoopAfterEntry) AfterEntry(ctx testctx.Context, entryValue reflect.Value, i int) error {
	return nil
}
