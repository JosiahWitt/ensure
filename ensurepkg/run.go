package ensurepkg

import (
	"fmt"
	"reflect"
	"testing"
)

// Run fn as a subtest called name.
func (e Ensure) Run(name string, fn func(ensure Ensure)) {
	c := e(nil)
	c.t.Helper()
	c.run(name, fn)
}

func (e Ensure) RunTableByIndex(table interface{}, fn func(ensure Ensure, i int)) {
	val := reflect.ValueOf(table)
	if val.Kind() != reflect.Array && val.Kind() != reflect.Slice {
		panic(fmt.Sprintf("Expected a slice or array for the table, got %T", table))
	}

	existingNames := make(map[string]int)

	for i := 0; i < val.Len(); i++ {
		i := i // Pin range variable

		entry := val.Index(i)
		if entry.Kind() != reflect.Struct {
			panic(fmt.Sprintf("Expected entry in table to be a struct, got %T", entry.Interface()))
		}

		entryName := entry.FieldByName("Name")
		zeroValue := reflect.Value{}
		if entryName == zeroValue {
			panic("Name field does not exist on struct in table")
		}

		if entryName.Kind() != reflect.String {
			panic("Name field in struct in table is not a string")
		}

		name := entryName.String()
		if name == "" {
			panic(fmt.Sprintf("Name not set for item with index %v", i))
		}

		if originalIndex, nameExists := existingNames[name]; nameExists {
			panic(fmt.Sprintf("Name duplicate found: index %v and %v have the same Name", originalIndex, i))
		}
		existingNames[name] = i

		c := e(nil)
		c.t.Helper()
		c.run(name, func(ensure Ensure) {
			fn(ensure, i)
		})
	}
}

func (c Chain) run(name string, fn func(ensure Ensure)) {
	c.t.Helper()
	c.t.Run(name, func(t *testing.T) {
		ensure := wrap(t)
		fn(ensure)
	})
}
