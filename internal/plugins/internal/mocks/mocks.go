// Package mocks provides an abstracted wrapper around mocks.
package mocks

import (
	"fmt"
	"reflect"
)

// All is a collection of [Mock]s.
type All struct {
	mocks  []*Mock
	byPath map[string]*Mock
}

// PathSet is a set of mock paths, used to keep track of which mocks are used.
type PathSet map[string]struct{}

// AddMock adds the mock to the [All] collection.
//
// It panics if a mock with that path was already added.
func (a *All) AddMock(path string, optional bool, t reflect.Type) *Mock {
	if prevMock, ok := a.byPath[path]; ok {
		panic(fmt.Sprintf("mock with path %q was already added: (PREVIOUS TYPE: %v, NEW TYPE: %v)", path, prevMock.t, t))
	}

	if a.byPath == nil {
		a.byPath = make(map[string]*Mock)
	}

	m := &Mock{
		Path:     path,
		Optional: optional,

		t:      t,
		values: make(map[int]reflect.Value),
	}

	a.byPath[path] = m
	a.mocks = append(a.mocks, m)
	return m
}

func (a *All) len() int {
	return len(a.mocks)
}

// Slice returns the underlying slice of mocks, and should not be modified.
func (a *All) Slice() []*Mock {
	return a.mocks
}

// PathSet constructs a new [PathSet] for the mocks, and is safe to modify.
func (a *All) PathSet() PathSet {
	m := make(PathSet, a.len())

	for k := range a.byPath {
		m[k] = struct{}{}
	}

	return m
}

// Mock holds metadata about a mock, and should not be modified directly.
type Mock struct {
	Path     string
	Optional bool

	t      reflect.Type
	values map[int]reflect.Value
}

// Implements returns true if the Mock implements the provided interface.
//
// It panics if the provided type is not an interface.
func (m *Mock) Implements(iface reflect.Type) bool {
	if iface.Kind() != reflect.Interface {
		panic(fmt.Sprintf("expected an interface to be provided to Implements, got: %v", iface))
	}

	return m.t.Implements(iface)
}

// SetValueByEntryIndex sets the value at the provided index.
//
// It panics if the type of v is does not match the expected type, or if a value has already been set at the index.
func (m *Mock) SetValueByEntryIndex(i int, v reflect.Value) {
	if actualType := v.Type(); actualType != m.t {
		panic(fmt.Sprintf("type of value for mock with path %q was not the expected type: (EXPECTED: %v, GOT: %v)", m.Path, m.t, actualType))
	}

	if _, ok := m.values[i]; ok {
		panic(fmt.Sprintf("value at index %d was already added for mock with path %q and type: %v", i, m.Path, m.t))
	}

	m.values[i] = v
}

// ValueByEntryIndex returns the value set for the provided index.
//
// It panics if the value was not set at that index.
func (m *Mock) ValueByEntryIndex(i int) reflect.Value {
	val, ok := m.values[i]
	if !ok {
		panic(fmt.Sprintf("value at index %d was not set for mock with path %q and type: %v", i, m.Path, m.t))
	}

	return val
}

// OnlyOneRequired returns true if and only if only one of the provided mocks is required.
func OnlyOneRequired(mocks ...*Mock) bool {
	foundRequired := false

	for _, mock := range mocks {
		if !mock.Optional {
			if foundRequired {
				return false // More than one is required
			}

			foundRequired = true
		}
	}

	return foundRequired
}
