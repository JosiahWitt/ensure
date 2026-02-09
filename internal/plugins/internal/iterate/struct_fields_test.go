package iterate_test

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/JosiahWitt/ensure"
	"github.com/JosiahWitt/ensure/ensuring"
	"github.com/JosiahWitt/ensure/internal/plugins/internal/iterate"
	"github.com/JosiahWitt/ensure/internal/stringerr"
)

func TestStructFields(t *testing.T) {
	ensure := ensure.New(t)

	const prefix = "Hello"

	type Field struct {
		FieldPath string
		Type      reflect.Type
	}

	table := []struct {
		Name string

		Struct   any
		Iterator func(fieldPath string, field *reflect.StructField) []error

		ExpectedFields []*Field
		ExpectedErrors []error
	}{
		{
			Name: "returns no errors when provided a struct with no fields",

			Struct: struct{}{},
		},
		{
			Name: "returns no errors when provided a struct with unexported fields",

			Struct: struct {
				notExported string
				notVisible  interface{ A() int }
			}{},
		},
		{
			Name: "returns no errors when provided a struct with exported fields",

			Struct: struct {
				AStruct     *struct{}
				notExported string
				AString     string
				notVisible  interface{ A() int }
				AnInterface interface{ B() int }
			}{},

			Iterator: func(fieldPath string, field *reflect.StructField) []error {
				return nil
			},

			ExpectedFields: []*Field{
				{
					FieldPath: prefix + ".AStruct",
					Type:      reflect.TypeFor[*struct{}](),
				},
				{
					FieldPath: prefix + ".AString",
					Type:      reflect.TypeFor[string](),
				},
				{
					FieldPath: prefix + ".AnInterface",
					Type:      reflect.TypeFor[interface{ B() int }](),
				},
			},
		},
		{
			Name: "returns no errors when provided a pointer to a struct with exported fields",

			Struct: &struct {
				AStruct     *struct{}
				notExported string
				AString     string
				notVisible  interface{ A() int }
				AnInterface interface{ B() int }
			}{},

			Iterator: func(fieldPath string, field *reflect.StructField) []error {
				return nil
			},

			ExpectedFields: []*Field{
				{
					FieldPath: prefix + ".AStruct",
					Type:      reflect.TypeFor[*struct{}](),
				},
				{
					FieldPath: prefix + ".AString",
					Type:      reflect.TypeFor[string](),
				},
				{
					FieldPath: prefix + ".AnInterface",
					Type:      reflect.TypeFor[interface{ B() int }](),
				},
			},
		},
		{
			Name: "returns errors when provided a struct with exported fields and iterator returns errors",

			Struct: struct {
				AStruct     *struct{}
				notExported string
				AString     string
				notVisible  interface{ A() int }
				AnInterface interface{ B() int }
			}{},

			Iterator: func(fieldPath string, field *reflect.StructField) []error {
				if strings.Contains(fieldPath, ".AStr") {
					return []error{stringerr.Newf("something's wrong with %s", fieldPath)}
				}

				return nil
			},

			ExpectedFields: []*Field{
				{
					FieldPath: prefix + ".AStruct",
					Type:      reflect.TypeFor[*struct{}](),
				},
				{
					FieldPath: prefix + ".AString",
					Type:      reflect.TypeFor[string](),
				},
				{
					FieldPath: prefix + ".AnInterface",
					Type:      reflect.TypeFor[interface{ B() int }](),
				},
			},

			ExpectedErrors: []error{
				stringerr.Newf("something's wrong with %s.AStruct", prefix),
				stringerr.Newf("something's wrong with %s.AString", prefix),
			},
		},
		{
			Name: "returns no errors when provided a struct with embedded structs",

			Struct: struct {
				Anonymous1
				AStruct     *struct{}
				notExported string
				AString     string
				Anonymous2
				notVisible  interface{ A() int }
				AnInterface interface{ B() int }
			}{},

			Iterator: func(fieldPath string, field *reflect.StructField) []error {
				return nil
			},

			ExpectedFields: []*Field{
				{
					FieldPath: prefix + ".Anonymous1.AStruct",
					Type:      reflect.TypeFor[*struct{}](),
				},
				{
					FieldPath: prefix + ".Anonymous1.AString",
					Type:      reflect.TypeFor[string](),
				},
				{
					FieldPath: prefix + ".Anonymous1.AnInterface",
					Type:      reflect.TypeFor[interface{ B() int }](),
				},
				{
					FieldPath: prefix + ".AStruct",
					Type:      reflect.TypeFor[*struct{}](),
				},
				{
					FieldPath: prefix + ".AString",
					Type:      reflect.TypeFor[string](),
				},
				{
					FieldPath: prefix + ".Anonymous2.AnInt",
					Type:      reflect.TypeFor[int](),
				},
				{
					FieldPath: prefix + ".AnInterface",
					Type:      reflect.TypeFor[interface{ B() int }](),
				},
			},
		},
		{
			Name: "returns no errors when provided a struct with embedded pointers to structs",

			Struct: struct {
				*Anonymous1
				AStruct     *struct{}
				notExported string
				AString     string
				*Anonymous2
				notVisible  interface{ A() int }
				AnInterface interface{ B() int }
			}{},

			Iterator: func(fieldPath string, field *reflect.StructField) []error {
				return nil
			},

			ExpectedFields: []*Field{
				{
					FieldPath: prefix + ".Anonymous1.AStruct",
					Type:      reflect.TypeFor[*struct{}](),
				},
				{
					FieldPath: prefix + ".Anonymous1.AString",
					Type:      reflect.TypeFor[string](),
				},
				{
					FieldPath: prefix + ".Anonymous1.AnInterface",
					Type:      reflect.TypeFor[interface{ B() int }](),
				},
				{
					FieldPath: prefix + ".AStruct",
					Type:      reflect.TypeFor[*struct{}](),
				},
				{
					FieldPath: prefix + ".AString",
					Type:      reflect.TypeFor[string](),
				},
				{
					FieldPath: prefix + ".Anonymous2.AnInt",
					Type:      reflect.TypeFor[int](),
				},
				{
					FieldPath: prefix + ".AnInterface",
					Type:      reflect.TypeFor[interface{ B() int }](),
				},
			},
		},
		{
			Name: "returns no errors when provided a struct with double nested embedded structs",

			Struct: struct {
				*Anonymous1Nested
				AStruct     *struct{}
				notExported string
				AString     string
				*Anonymous2Nested
				notVisible  interface{ A() int }
				AnInterface interface{ B() int }
			}{},

			Iterator: func(fieldPath string, field *reflect.StructField) []error {
				return nil
			},

			ExpectedFields: []*Field{
				{
					FieldPath: prefix + ".Anonymous1Nested.AFloat",
					Type:      reflect.TypeFor[float64](),
				},
				{
					FieldPath: prefix + ".Anonymous1Nested.Anonymous1.AStruct",
					Type:      reflect.TypeFor[*struct{}](),
				},
				{
					FieldPath: prefix + ".Anonymous1Nested.Anonymous1.AString",
					Type:      reflect.TypeFor[string](),
				},
				{
					FieldPath: prefix + ".Anonymous1Nested.Anonymous1.AnInterface",
					Type:      reflect.TypeFor[interface{ B() int }](),
				},
				{
					FieldPath: prefix + ".AStruct",
					Type:      reflect.TypeFor[*struct{}](),
				},
				{
					FieldPath: prefix + ".AString",
					Type:      reflect.TypeFor[string](),
				},
				{
					FieldPath: prefix + ".Anonymous2Nested.ABool",
					Type:      reflect.TypeFor[bool](),
				},
				{
					FieldPath: prefix + ".Anonymous2Nested.Anonymous2.AnInt",
					Type:      reflect.TypeFor[int](),
				},
				{
					FieldPath: prefix + ".AnInterface",
					Type:      reflect.TypeFor[interface{ B() int }](),
				},
			},
		},
		{
			Name: "returns no errors when provided a struct with ambiguous field names",

			Struct: struct {
				AString string
				Anonymous1
				AStruct     *struct{ Different bool }
				notExported string
				*Anonymous2
				AnInt       int
				notVisible  interface{ A() int }
				AnInterface interface{ B() int }
			}{},

			Iterator: func(fieldPath string, field *reflect.StructField) []error {
				return nil
			},

			ExpectedFields: []*Field{
				{
					FieldPath: prefix + ".AString",
					Type:      reflect.TypeFor[string](),
				},
				{
					FieldPath: prefix + ".Anonymous1.AStruct",
					Type:      reflect.TypeFor[*struct{}](),
				},
				{
					FieldPath: prefix + ".Anonymous1.AString",
					Type:      reflect.TypeFor[string](),
				},
				{
					FieldPath: prefix + ".Anonymous1.AnInterface",
					Type:      reflect.TypeFor[interface{ B() int }](),
				},
				{
					FieldPath: prefix + ".AStruct",
					Type:      reflect.TypeFor[*struct{ Different bool }](),
				},
				{
					FieldPath: prefix + ".Anonymous2.AnInt",
					Type:      reflect.TypeFor[int](),
				},
				{
					FieldPath: prefix + ".AnInt",
					Type:      reflect.TypeFor[int](),
				},
				{
					FieldPath: prefix + ".AnInterface",
					Type:      reflect.TypeFor[interface{ B() int }](),
				},
			},
		},
		{
			Name: "returns no errors when provided a struct with recursively nested embedded structs",

			Struct: struct {
				*AnonymousRecursiveNested
				AStruct     *struct{}
				notExported string
				AString     string
				*AnonymousDoubleRecursiveNested
				notVisible  interface{ A() int }
				AnInterface interface{ B() int }
			}{},

			Iterator: func(fieldPath string, field *reflect.StructField) []error {
				return nil
			},

			ExpectedFields: []*Field{
				{
					FieldPath: prefix + ".AnonymousRecursiveNested.ABool",
					Type:      reflect.TypeFor[bool](),
				},
				{
					FieldPath: prefix + ".AnonymousRecursiveNested.ASlice",
					Type:      reflect.TypeFor[[]string](),
				},
				{
					FieldPath: prefix + ".AStruct",
					Type:      reflect.TypeFor[*struct{}](),
				},
				{
					FieldPath: prefix + ".AString",
					Type:      reflect.TypeFor[string](),
				},
				{
					FieldPath: prefix + ".AnonymousDoubleRecursiveNested.AnInt",
					Type:      reflect.TypeFor[int](),
				},
				{
					FieldPath: prefix + ".AnonymousDoubleRecursiveNested.AnonymousRecursiveNested.ABool",
					Type:      reflect.TypeFor[bool](),
				},
				{
					FieldPath: prefix + ".AnonymousDoubleRecursiveNested.AnonymousRecursiveNested.ASlice",
					Type:      reflect.TypeFor[[]string](),
				},
				{
					FieldPath: prefix + ".AnInterface",
					Type:      reflect.TypeFor[interface{ B() int }](),
				},
			},
		},
		{
			Name: "treats all embedded non-struct types as regular fields",

			Struct: struct {
				AnonymousInterface // Embedded interface
				AnonymousSlice     // Embedded slice type
				AnonymousMap       // Embedded map type
				AnonymousChan      // Embedded chan type
				AnonymousInt       // Embedded basic type
				AnonymousFunc      // Embedded func type
				*AnonymousString   // Embedded pointer to basic type alias
				AString            string
			}{},

			Iterator: func(fieldPath string, field *reflect.StructField) []error {
				return nil
			},

			ExpectedFields: []*Field{
				{
					FieldPath: prefix + ".AnonymousInterface",
					Type:      reflect.TypeFor[AnonymousInterface](),
				},
				{
					FieldPath: prefix + ".AnonymousSlice",
					Type:      reflect.TypeFor[AnonymousSlice](),
				},
				{
					FieldPath: prefix + ".AnonymousMap",
					Type:      reflect.TypeFor[AnonymousMap](),
				},
				{
					FieldPath: prefix + ".AnonymousChan",
					Type:      reflect.TypeFor[AnonymousChan](),
				},
				{
					FieldPath: prefix + ".AnonymousInt",
					Type:      reflect.TypeFor[AnonymousInt](),
				},
				{
					FieldPath: prefix + ".AnonymousFunc",
					Type:      reflect.TypeFor[AnonymousFunc](),
				},
				{
					FieldPath: prefix + ".AnonymousString",
					Type:      reflect.TypeFor[*AnonymousString](),
				},
				{
					FieldPath: prefix + ".AString",
					Type:      reflect.TypeFor[string](),
				},
			},
		},
		{
			Name: "handles mixed embedded types: struct, interface, and other types together",

			Struct: struct {
				Anonymous1         // Embedded struct - recursed into
				AnonymousInterface // Embedded interface - regular field
				AnonymousSlice     // Embedded slice - regular field
				AString            string
			}{},

			Iterator: func(fieldPath string, field *reflect.StructField) []error {
				return nil
			},

			ExpectedFields: []*Field{
				{
					FieldPath: prefix + ".Anonymous1.AStruct",
					Type:      reflect.TypeFor[*struct{}](),
				},
				{
					FieldPath: prefix + ".Anonymous1.AString",
					Type:      reflect.TypeFor[string](),
				},
				{
					FieldPath: prefix + ".Anonymous1.AnInterface",
					Type:      reflect.TypeFor[interface{ B() int }](),
				},
				{
					FieldPath: prefix + ".AnonymousInterface",
					Type:      reflect.TypeFor[AnonymousInterface](),
				},
				{
					FieldPath: prefix + ".AnonymousSlice",
					Type:      reflect.TypeFor[AnonymousSlice](),
				},
				{
					FieldPath: prefix + ".AString",
					Type:      reflect.TypeFor[string](),
				},
			},
		},
		{
			Name: "iterator receives embedded interface fields and can return errors",

			Struct: struct {
				AnonymousInterface // Embedded interface - treated as regular field
				AString            string
			}{},

			Iterator: func(fieldPath string, field *reflect.StructField) []error {
				if field.Type.Kind() == reflect.Interface {
					return []error{stringerr.Newf("iterator error for interface field: %s", fieldPath)}
				}
				return nil
			},

			ExpectedFields: []*Field{
				{
					FieldPath: prefix + ".AnonymousInterface",
					Type:      reflect.TypeFor[AnonymousInterface](),
				},
				{
					FieldPath: prefix + ".AString",
					Type:      reflect.TypeFor[string](),
				},
			},

			ExpectedErrors: []error{
				stringerr.Newf("iterator error for interface field: %s.AnonymousInterface", prefix),
			},
		},
		{
			Name: "returns errors when provided a struct with double nested embedded structs and iterator returns errors",

			Struct: struct {
				*Anonymous1Nested
				AStruct     *struct{}
				notExported string
				AString     string
				*Anonymous2Nested
				notVisible  interface{ A() int }
				AnInterface interface{ B() int }
			}{},

			Iterator: func(fieldPath string, field *reflect.StructField) []error {
				if strings.Contains(fieldPath, ".AStr") {
					return []error{stringerr.Newf("something's wrong with %s", fieldPath)}
				}

				return nil
			},

			ExpectedFields: []*Field{
				{
					FieldPath: prefix + ".Anonymous1Nested.AFloat",
					Type:      reflect.TypeFor[float64](),
				},
				{
					FieldPath: prefix + ".Anonymous1Nested.Anonymous1.AStruct",
					Type:      reflect.TypeFor[*struct{}](),
				},
				{
					FieldPath: prefix + ".Anonymous1Nested.Anonymous1.AString",
					Type:      reflect.TypeFor[string](),
				},
				{
					FieldPath: prefix + ".Anonymous1Nested.Anonymous1.AnInterface",
					Type:      reflect.TypeFor[interface{ B() int }](),
				},
				{
					FieldPath: prefix + ".AStruct",
					Type:      reflect.TypeFor[*struct{}](),
				},
				{
					FieldPath: prefix + ".AString",
					Type:      reflect.TypeFor[string](),
				},
				{
					FieldPath: prefix + ".Anonymous2Nested.ABool",
					Type:      reflect.TypeFor[bool](),
				},
				{
					FieldPath: prefix + ".Anonymous2Nested.Anonymous2.AnInt",
					Type:      reflect.TypeFor[int](),
				},
				{
					FieldPath: prefix + ".AnInterface",
					Type:      reflect.TypeFor[interface{ B() int }](),
				},
			},

			ExpectedErrors: []error{
				stringerr.Newf("something's wrong with %s.AStruct", prefix),
				stringerr.Newf("something's wrong with %s.AString", prefix),
			},
		},
	}

	ensure.RunTableByIndex(table, func(ensure ensuring.E, i int) {
		entry := table[i]

		var fields []*Field
		iterator := func(fieldPath string, field *reflect.StructField) []error {
			fields = append(fields, &Field{FieldPath: fieldPath, Type: field.Type})
			return entry.Iterator(fieldPath, field)
		}

		_, errs := iterate.StructFields(prefix, reflect.TypeOf(entry.Struct), iterator)
		ensure(errorContainer(errs)).MatchesAllErrors(entry.ExpectedErrors...)
		ensure(fields).Equals(entry.ExpectedFields)
	})

	ensure.Run("panics when provided a non-struct", func(ensure ensuring.E) {
		defer func() {
			ensure(recover()).Equals("StructFields must be provided a struct, got: string")
		}()

		iterate.StructFields("", reflect.TypeFor[string](), nil)
	})
}

func TestInitializeStruct(t *testing.T) {
	ensure := ensure.New(t)

	const prefix = "Hello"

	noopStructFieldsIterator := func(fieldPath string, field *reflect.StructField) []error { return nil }

	type Field struct {
		FieldPath string
		Type      reflect.Type
	}

	table := []struct {
		Name string

		NestedStruct any
		Iterator     iterate.InitializeStructIterator

		ExpectedInitializedStruct any
		ExpectedFields            []*Field
	}{
		{
			Name: "initializes struct when provided a struct with no fields",

			NestedStruct: &struct{ Struct *struct{} }{},
			Iterator:     func(fieldPath string, field reflect.Value) {},

			ExpectedInitializedStruct: &struct{}{},
		},
		{
			Name: "initializes struct when provided a struct with unexported fields",

			NestedStruct: &struct {
				Struct *struct {
					notExported string
					notVisible  interface{ A() int }
				}
			}{},
			Iterator: func(fieldPath string, field reflect.Value) {},

			ExpectedInitializedStruct: &struct {
				notExported string
				notVisible  interface{ A() int }
			}{},
		},
		{
			Name: "initializes struct when provided a struct with exported fields",

			NestedStruct: &struct {
				Struct *struct {
					AStruct     *struct{}
					notExported string
					AString     string
					notVisible  interface{ A() int }
					AnInterface interface{ B() int }
				}
			}{},
			Iterator: func(fieldPath string, field reflect.Value) {},

			ExpectedInitializedStruct: &struct {
				AStruct     *struct{}
				notExported string
				AString     string
				notVisible  interface{ A() int }
				AnInterface interface{ B() int }
			}{},

			ExpectedFields: []*Field{
				{
					FieldPath: prefix + ".AStruct",
					Type:      reflect.TypeFor[*struct{}](),
				},
				{
					FieldPath: prefix + ".AString",
					Type:      reflect.TypeFor[string](),
				},
				{
					FieldPath: prefix + ".AnInterface",
					Type:      reflect.TypeFor[interface{ B() int }](),
				},
			},
		},
		{
			Name: "initializes struct when provided a struct with embedded structs",

			NestedStruct: &struct {
				Struct *struct {
					Anonymous1
					AStruct     *struct{}
					notExported string
					AString     string
					Anonymous2
					notVisible  interface{ A() int }
					AnInterface interface{ B() int }
				}
			}{},
			Iterator: func(fieldPath string, field reflect.Value) {},

			ExpectedInitializedStruct: &struct {
				Anonymous1
				AStruct     *struct{}
				notExported string
				AString     string
				Anonymous2
				notVisible  interface{ A() int }
				AnInterface interface{ B() int }
			}{
				Anonymous1: Anonymous1{},
				Anonymous2: Anonymous2{},
			},

			ExpectedFields: []*Field{
				{
					FieldPath: prefix + ".AStruct",
					Type:      reflect.TypeFor[*struct{}](),
				},
				{
					FieldPath: prefix + ".AString",
					Type:      reflect.TypeFor[string](),
				},
				{
					FieldPath: prefix + ".AnInterface",
					Type:      reflect.TypeFor[interface{ B() int }](),
				},
				{
					FieldPath: prefix + ".Anonymous1.AStruct",
					Type:      reflect.TypeFor[*struct{}](),
				},
				{
					FieldPath: prefix + ".Anonymous1.AString",
					Type:      reflect.TypeFor[string](),
				},
				{
					FieldPath: prefix + ".Anonymous1.AnInterface",
					Type:      reflect.TypeFor[interface{ B() int }](),
				},
				{
					FieldPath: prefix + ".Anonymous2.AnInt",
					Type:      reflect.TypeFor[int](),
				},
			},
		},
		{
			Name: "initializes struct when provided a struct with embedded pointers to structs",

			NestedStruct: &struct {
				Struct *struct {
					*Anonymous1
					AStruct     *struct{}
					notExported string
					AString     string
					*Anonymous2
					notVisible  interface{ A() int }
					AnInterface interface{ B() int }
				}
			}{},
			Iterator: func(fieldPath string, field reflect.Value) {},

			ExpectedInitializedStruct: &struct {
				*Anonymous1
				AStruct     *struct{}
				notExported string
				AString     string
				*Anonymous2
				notVisible  interface{ A() int }
				AnInterface interface{ B() int }
			}{
				Anonymous1: &Anonymous1{},
				Anonymous2: &Anonymous2{},
			},

			ExpectedFields: []*Field{
				{
					FieldPath: prefix + ".AStruct",
					Type:      reflect.TypeFor[*struct{}](),
				},
				{
					FieldPath: prefix + ".AString",
					Type:      reflect.TypeFor[string](),
				},
				{
					FieldPath: prefix + ".AnInterface",
					Type:      reflect.TypeFor[interface{ B() int }](),
				},
				{
					FieldPath: prefix + ".Anonymous1.AStruct",
					Type:      reflect.TypeFor[*struct{}](),
				},
				{
					FieldPath: prefix + ".Anonymous1.AString",
					Type:      reflect.TypeFor[string](),
				},
				{
					FieldPath: prefix + ".Anonymous1.AnInterface",
					Type:      reflect.TypeFor[interface{ B() int }](),
				},
				{
					FieldPath: prefix + ".Anonymous2.AnInt",
					Type:      reflect.TypeFor[int](),
				},
			},
		},
		{
			Name: "initializes struct when provided a struct with double nested embedded structs",

			NestedStruct: &struct {
				Struct *struct {
					Anonymous1Nested                        // Non-pointer with a embedded pointer to show it recurses
					AStruct          *struct{ Name string } // Different type than the embedded one with the same field name to show it picks the right one
					notExported      string
					AString          string
					*Anonymous2Nested
					notVisible  interface{ A() int }
					AnInterface interface{ B() int }
				}
			}{},
			Iterator: func(fieldPath string, field reflect.Value) {},

			ExpectedInitializedStruct: &struct {
				Anonymous1Nested
				AStruct     *struct{ Name string }
				notExported string
				AString     string
				*Anonymous2Nested
				notVisible  interface{ A() int }
				AnInterface interface{ B() int }
			}{
				Anonymous1Nested: Anonymous1Nested{
					Anonymous1: &Anonymous1{},
				},
				Anonymous2Nested: &Anonymous2Nested{
					Anonymous2: Anonymous2{},
				},
			},

			ExpectedFields: []*Field{
				{
					FieldPath: prefix + ".AStruct",
					Type:      reflect.TypeFor[*struct{ Name string }](),
				},
				{
					FieldPath: prefix + ".AString",
					Type:      reflect.TypeFor[string](),
				},
				{
					FieldPath: prefix + ".AnInterface",
					Type:      reflect.TypeFor[interface{ B() int }](),
				},
				{
					FieldPath: prefix + ".Anonymous1Nested.AFloat",
					Type:      reflect.TypeFor[float64](),
				},
				{
					FieldPath: prefix + ".Anonymous1Nested.Anonymous1.AStruct",
					Type:      reflect.TypeFor[*struct{}](),
				},
				{
					FieldPath: prefix + ".Anonymous1Nested.Anonymous1.AString",
					Type:      reflect.TypeFor[string](),
				},
				{
					FieldPath: prefix + ".Anonymous1Nested.Anonymous1.AnInterface",
					Type:      reflect.TypeFor[interface{ B() int }](),
				},
				{
					FieldPath: prefix + ".Anonymous2Nested.ABool",
					Type:      reflect.TypeFor[bool](),
				},
				{
					FieldPath: prefix + ".Anonymous2Nested.Anonymous2.AnInt",
					Type:      reflect.TypeFor[int](),
				},
			},
		},
		{
			Name: "returns no errors when provided a struct with recursively nested embedded structs",

			NestedStruct: &struct {
				Struct *struct {
					*AnonymousRecursiveNested
					AStruct     *struct{}
					notExported string
					AString     string
					*AnonymousDoubleRecursiveNested
					notVisible  interface{ A() int }
					AnInterface interface{ B() int }
				}
			}{},

			Iterator: func(fieldPath string, field reflect.Value) {},

			ExpectedInitializedStruct: &struct {
				*AnonymousRecursiveNested
				AStruct     *struct{}
				notExported string
				AString     string
				*AnonymousDoubleRecursiveNested
				notVisible  interface{ A() int }
				AnInterface interface{ B() int }
			}{
				AnonymousRecursiveNested: &AnonymousRecursiveNested{},
				AnonymousDoubleRecursiveNested: &AnonymousDoubleRecursiveNested{
					AnonymousRecursiveNested: &AnonymousRecursiveNested{},
				},
			},

			ExpectedFields: []*Field{
				{
					FieldPath: prefix + ".AStruct",
					Type:      reflect.TypeFor[*struct{}](),
				},
				{
					FieldPath: prefix + ".AString",
					Type:      reflect.TypeFor[string](),
				},
				{
					FieldPath: prefix + ".AnInterface",
					Type:      reflect.TypeFor[interface{ B() int }](),
				},
				{
					FieldPath: prefix + ".AnonymousRecursiveNested.ABool",
					Type:      reflect.TypeFor[bool](),
				},
				{
					FieldPath: prefix + ".AnonymousRecursiveNested.ASlice",
					Type:      reflect.TypeFor[[]string](),
				},
				{
					FieldPath: prefix + ".AnonymousDoubleRecursiveNested.AnInt",
					Type:      reflect.TypeFor[int](),
				},
				{
					FieldPath: prefix + ".AnonymousDoubleRecursiveNested.AnonymousRecursiveNested.ABool",
					Type:      reflect.TypeFor[bool](),
				},
				{
					FieldPath: prefix + ".AnonymousDoubleRecursiveNested.AnonymousRecursiveNested.ASlice",
					Type:      reflect.TypeFor[[]string](),
				},
			},
		},
		{
			Name: "runs iterator over every item",

			NestedStruct: &struct {
				Struct *struct {
					Anonymous1Nested
					AStruct     *struct{ Name string }
					notExported string
					AString     string
					*Anonymous2Nested
					notVisible  interface{ A() int }
					AnInterface interface{ B() int }
					*AnonymousRecursiveNested
					*AnonymousDoubleRecursiveNested
					ASlice []string
				}
			}{},

			Iterator: func(fieldPath string, field reflect.Value) {
				switch fieldPath {
				// Show we can set a different value to a shadowed field
				case prefix + ".AString":
					field.Set(reflect.ValueOf("Hello"))
				case prefix + ".Anonymous1Nested.Anonymous1.AString":
					field.Set(reflect.ValueOf("World"))

					// Show recursive embedded fields are handled correctly
				case prefix + ".AnonymousRecursiveNested.ASlice":
					field.Set(reflect.ValueOf([]string{"a", "b", "c"}))
				case prefix + ".AnonymousDoubleRecursiveNested.AnonymousRecursiveNested.ASlice":
					field.Set(reflect.ValueOf([]string{"x", "y", "z"}))
				case prefix + ".ASlice":
					field.Set(reflect.ValueOf([]string{"q", "w", "e", "r", "t", "y"}))
				}
			},

			ExpectedInitializedStruct: &struct {
				Anonymous1Nested
				AStruct     *struct{ Name string }
				notExported string
				AString     string
				*Anonymous2Nested
				notVisible  interface{ A() int }
				AnInterface interface{ B() int }
				*AnonymousRecursiveNested
				*AnonymousDoubleRecursiveNested
				ASlice []string
			}{
				Anonymous1Nested: Anonymous1Nested{
					Anonymous1: &Anonymous1{
						AString: "World",
					},
				},
				AString: "Hello",
				Anonymous2Nested: &Anonymous2Nested{
					Anonymous2: Anonymous2{},
				},
				AnonymousRecursiveNested: &AnonymousRecursiveNested{
					ASlice: []string{"a", "b", "c"},
				},
				AnonymousDoubleRecursiveNested: &AnonymousDoubleRecursiveNested{
					AnonymousRecursiveNested: &AnonymousRecursiveNested{
						ASlice: []string{"x", "y", "z"},
					},
				},
				ASlice: []string{"q", "w", "e", "r", "t", "y"},
			},

			ExpectedFields: []*Field{
				{
					FieldPath: prefix + ".AStruct",
					Type:      reflect.TypeFor[*struct{ Name string }](),
				},
				{
					FieldPath: prefix + ".AString",
					Type:      reflect.TypeFor[string](),
				},
				{
					FieldPath: prefix + ".AnInterface",
					Type:      reflect.TypeFor[interface{ B() int }](),
				},
				{
					FieldPath: prefix + ".ASlice",
					Type:      reflect.TypeFor[[]string](),
				},
				{
					FieldPath: prefix + ".Anonymous1Nested.AFloat",
					Type:      reflect.TypeFor[float64](),
				},
				{
					FieldPath: prefix + ".Anonymous1Nested.Anonymous1.AStruct",
					Type:      reflect.TypeFor[*struct{}](),
				},
				{
					FieldPath: prefix + ".Anonymous1Nested.Anonymous1.AString",
					Type:      reflect.TypeFor[string](),
				},
				{
					FieldPath: prefix + ".Anonymous1Nested.Anonymous1.AnInterface",
					Type:      reflect.TypeFor[interface{ B() int }](),
				},
				{
					FieldPath: prefix + ".Anonymous2Nested.ABool",
					Type:      reflect.TypeFor[bool](),
				},
				{
					FieldPath: prefix + ".Anonymous2Nested.Anonymous2.AnInt",
					Type:      reflect.TypeFor[int](),
				},
				{
					FieldPath: prefix + ".AnonymousRecursiveNested.ABool",
					Type:      reflect.TypeFor[bool](),
				},
				{
					FieldPath: prefix + ".AnonymousRecursiveNested.ASlice",
					Type:      reflect.TypeFor[[]string](),
				},
				{
					FieldPath: prefix + ".AnonymousDoubleRecursiveNested.AnInt",
					Type:      reflect.TypeFor[int](),
				},
				{
					FieldPath: prefix + ".AnonymousDoubleRecursiveNested.AnonymousRecursiveNested.ABool",
					Type:      reflect.TypeFor[bool](),
				},
				{
					FieldPath: prefix + ".AnonymousDoubleRecursiveNested.AnonymousRecursiveNested.ASlice",
					Type:      reflect.TypeFor[[]string](),
				},
			},
		},
	}

	ensure.RunTableByIndex(table, func(ensure ensuring.E, i int) {
		entry := table[i]
		s := reflect.ValueOf(entry.NestedStruct).Elem().FieldByName("Struct")
		ensure(s == (reflect.Value{})).IsFalse()

		res, errs := iterate.StructFields(prefix, s.Type(), noopStructFieldsIterator)
		ensure(errs).IsNotError()

		var fields []*Field
		iterator := func(fieldPath string, v reflect.Value) {
			fields = append(fields, &Field{FieldPath: fieldPath, Type: v.Type()})
			entry.Iterator(fieldPath, v)
		}

		res.InitializeStruct(s, iterator)
		ensure(s.Interface()).Equals(entry.ExpectedInitializedStruct)
		ensure(fields).Equals(entry.ExpectedFields)
	})

	ensure.Run("panics when not provided a pointer", func(ensure ensuring.E) {
		defer func() {
			ensure(recover()).Equals("InitializeStruct must be provided a pointer to a struct, got: struct {}")
		}()

		res, errs := iterate.StructFields(prefix, reflect.TypeFor[struct{}](), noopStructFieldsIterator)
		ensure(errs).IsNotError()

		res.InitializeStruct(reflect.ValueOf(struct{}{}), nil)
	})

	ensure.Run("panics when not provided a struct pointer", func(ensure ensuring.E) {
		defer func() {
			ensure(recover()).Equals("InitializeStruct must be provided a pointer to a struct, got: *string")
		}()

		res, errs := iterate.StructFields(prefix, reflect.TypeFor[struct{}](), noopStructFieldsIterator)
		ensure(errs).IsNotError()

		str := "not a struct"
		ptr := &str
		res.InitializeStruct(reflect.ValueOf(ptr), nil)
	})

	ensure.Run("panics when type doesn't match", func(ensure ensuring.E) {
		defer func() {
			ensure(recover()).Equals("InitializeStruct must be provided the type (iterate_test.s1) that was provided to StructFields, got: *iterate_test.s2")
		}()

		type s1 struct{}
		type s2 struct{}

		res, errs := iterate.StructFields(prefix, reflect.TypeFor[s1](), noopStructFieldsIterator)
		ensure(errs).IsNotError()

		res.InitializeStruct(reflect.ValueOf(&s2{}), nil)
	})

	ensure.Run("panics when type is not addressable", func(ensure ensuring.E) {
		defer func() {
			ensure(recover()).Equals("InitializeStruct must be provided an addressable value, such as a field inside a pointer to a struct or an element in a slice, got: *iterate_test.s")
		}()

		type s struct{}

		res, errs := iterate.StructFields(prefix, reflect.TypeFor[s](), noopStructFieldsIterator)
		ensure(errs).IsNotError()

		res.InitializeStruct(reflect.ValueOf(&s{}), nil)
	})
}

type (
	//lint:ignore U1000 not used for test purposes
	Anonymous1 struct {
		AStruct     *struct{}
		notExported string
		AString     string
		notVisible  interface{ A() int }
		AnInterface interface{ B() int }
	}

	//lint:ignore U1000 not used for test purposes
	Anonymous2 struct {
		AnInt       int
		notExported int
	}

	//lint:ignore U1000 not used for test purposes
	Anonymous1Nested struct {
		AFloat float64
		*Anonymous1
		notExported int
	}

	//lint:ignore U1000 not used for test purposes
	Anonymous2Nested struct {
		ABool bool
		Anonymous2
		notExported int
	}

	//lint:ignore U1000 not used for test purposes
	AnonymousRecursiveNested struct {
		ABool bool
		*AnonymousRecursiveNested
		ASlice      []string
		notExported float64
	}

	//lint:ignore U1000 not used for test purposes
	AnonymousDoubleRecursiveNested struct {
		AnInt int
		*AnonymousRecursiveNested
		notExported float64
		*AnonymousDoubleRecursiveNested
	}

	AnonymousInterface interface{ C() int }
	AnonymousString    string
	AnonymousSlice     []string
	AnonymousMap       map[string]int
	AnonymousChan      chan int
	AnonymousInt       int
	AnonymousFunc      func() error
)

type errorContainer []error //nolint:errname

func (errs errorContainer) Is(target error) bool {
	if target == nil {
		return len(errs) == 0
	}

	for _, err := range errs {
		if errors.Is(err, target) {
			return true
		}
	}

	return false
}

func (errs errorContainer) Error() string {
	return fmt.Sprintf("%v", []error(errs))
}
