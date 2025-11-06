package iterate

import (
	"fmt"
	"reflect"

	"github.com/JosiahWitt/ensure/internal/stringerr"
)

// StructFieldsIterator is called by [StructFields] for every exported field.
// The fieldPath is also provided to [InitializeStructIterator].
type StructFieldsIterator func(fieldPath string, field *reflect.StructField) []error

// StructFields supports looping through all the exported fields in a struct, including embedded structs. All errors, including those
// from [iterator], are collected into a slice of errors.
//
// If keeping track of fields iterated over, use the path provided as the fieldPath parameter to [StructFieldsIterator], since that's
// safe for shadowed fields in embedded structs. The path names will be identically provided when iterating in [InitializeStruct].
// The iteration order is NOT guaranteed to be the same.
//
// Panics if the provided type is not a struct or a pointer to a struct.
func StructFields(prefix string, t reflect.Type, iterator StructFieldsIterator) (*StructFieldsResult, []error) {
	i := &structIterator{
		iterator: iterator,

		visitedAnonymousTypes: make(map[reflect.Type]struct{}),
	}

	return i.process(prefix, t)
}

type structIterator struct {
	iterator StructFieldsIterator

	visitedAnonymousTypes map[reflect.Type]struct{}
}

//nolint:funlen,cyclop // It seems clearer as one larger method
func (si *structIterator) process(prefix string, t reflect.Type) (*StructFieldsResult, []error) {
	t = indirectType(t)
	if t.Kind() != reflect.Struct {
		panicf("StructFields must be provided a struct, got: %v", t)
	}

	var allErrs []error
	result := &StructFieldsResult{t: t}

	for i := range t.NumField() {
		field := t.Field(i)
		fieldType := field.Type
		fieldKind := fieldType.Kind()
		fieldName := field.Name
		fieldPath := prefix + "." + fieldName

		// Skip unexported fields
		// TODO: Swap to use the IsExported method once Go 1.17 is the minimum supported version
		if field.PkgPath != "" {
			continue
		}

		// Support embedded structs
		if field.Anonymous {
			// To prevent infinitely looping through recursive anonymous fields, we only iterate through the first level.
			if _, ok := si.visitedAnonymousTypes[fieldType]; ok {
				continue
			}

			isPointer := fieldKind == reflect.Ptr
			isPointerToStruct := isPointer && fieldType.Elem().Kind() == reflect.Struct

			if fieldKind != reflect.Struct && !isPointerToStruct {
				err := stringerr.Newf("expected %s to be an embedded struct, got: %v", fieldPath, fieldType)
				allErrs = append(allErrs, err)
				continue
			}

			si.visitedAnonymousTypes[fieldType] = struct{}{}

			nestedResult, errs := si.process(fieldPath, fieldType)
			if len(errs) != 0 {
				allErrs = append(allErrs, errs...)
				continue
			}

			// If the anonymous field type appears separately, we want to include it again, thus it's safe to delete here.
			delete(si.visitedAnonymousTypes, fieldType)

			nestedResult.fieldName = fieldName
			nestedResult.isPointer = isPointer
			result.anonymousFields = append(result.anonymousFields, nestedResult)

			continue
		}

		result.fields = append(result.fields, &structField{
			path: fieldPath,
			name: fieldName,
		})

		if errs := si.iterator(fieldPath, &field); len(errs) != 0 {
			allErrs = append(allErrs, errs...)
			continue
		}
	}

	return result, allErrs
}

// StructFieldsResult is the result returned from [StructFields], which allows populating struct fields.
type StructFieldsResult struct {
	fieldName string
	t         reflect.Type
	isPointer bool

	fields          []*structField
	anonymousFields []*StructFieldsResult
}

type structField struct {
	path string
	name string
}

// InitializeStructIterator is called by [InitializeStruct] for every exported field.
// The fieldPath matches those provided to [StructFieldsIterator].
type InitializeStructIterator func(fieldPath string, field reflect.Value)

// InitializeStruct sets the struct to the zero value and recursively initializes all the embedded structs.
//
// Panics if the provided type is not a pointer to a struct, if the value is not the same type as the one
// provided to StructFields, or if the value is not addressable.
func (r *StructFieldsResult) InitializeStruct(v reflect.Value, iterator InitializeStructIterator) {
	r.initialize(v)
	r.initializeFields(v)
	r.iterate(v, iterator)
}

func (r *StructFieldsResult) initialize(v reflect.Value) {
	t := v.Type()

	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		panicf("InitializeStruct must be provided a pointer to a struct, got: %v", t)
	}

	if t != r.t && t.Elem() != r.t {
		panicf("InitializeStruct must be provided the type (%v) that was provided to StructFields, got: %v", r.t, t)
	}

	if !v.CanAddr() {
		panicf("InitializeStruct must be provided an addressable value, such as a field inside a pointer to a struct or an element in a slice, got: %v", t)
	}

	v.Set(reflect.New(t.Elem()))
}

func (r *StructFieldsResult) initializeFields(v reflect.Value) {
	v = reflect.Indirect(v)

	for _, metadata := range r.anonymousFields {
		field := v.FieldByName(metadata.fieldName)

		if metadata.isPointer {
			metadata.initialize(field)
		}

		metadata.initializeFields(field)
	}
}

func (r *StructFieldsResult) iterate(v reflect.Value, iterator func(fieldPath string, field reflect.Value)) {
	v = reflect.Indirect(v)

	for _, metadata := range r.fields {
		field := v.FieldByName(metadata.name)
		iterator(metadata.path, field)
	}

	for _, metadata := range r.anonymousFields {
		field := v.FieldByName(metadata.fieldName)
		metadata.iterate(field, iterator)
	}
}

func indirectType(t reflect.Type) reflect.Type {
	if t.Kind() != reflect.Ptr {
		return t
	}

	return indirectType(t.Elem())
}

func panicf(format string, a ...interface{}) {
	panic(fmt.Sprintf(format, a...))
}
