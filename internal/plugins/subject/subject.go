// Package subject provides a plugin that initializes and populates test entry Subject fields using the provided mocks.
package subject

import (
	"fmt"
	"reflect"

	"github.com/JosiahWitt/ensure/internal/plugins"
	"github.com/JosiahWitt/ensure/internal/plugins/internal/id"
	"github.com/JosiahWitt/ensure/internal/plugins/internal/iterate"
	"github.com/JosiahWitt/ensure/internal/plugins/internal/mocks"
	"github.com/JosiahWitt/ensure/internal/stringerr"
	"github.com/JosiahWitt/ensure/internal/testctx"
)

// New uses the collection of mocks to initialize and populate the Subject field in a test entry.
// The subject plugin should come after any steps that setup mocks.
func New(m *mocks.All) *TablePlugin {
	return &TablePlugin{mocks: m}
}

// TablePlugin uses the collection of mocks to initialize and populate the Subject field in a test entry.
// The subject plugin should come after any steps that setup mocks.
type TablePlugin struct {
	mocks *mocks.All
}

var _ plugins.TablePlugin = &TablePlugin{}

// ParseEntryType is called during the first pass of plugin initialization.
// It is responsible for making sure the type is as expected.
func (t *TablePlugin) ParseEntryType(entryType reflect.Type) (plugins.TableEntryHooks, error) {
	h := &TableEntryHooks{}

	subjectStruct, ok := entryType.FieldByName(id.Subject)
	if ok {
		if err := validateSubjectFieldType(&subjectStruct); err != nil {
			return nil, err
		}

		subjectMocks, structFieldsResult, err := t.mocksForSubject(&subjectStruct)
		if err != nil {
			return nil, err
		}

		h.hasSubject = true
		h.subjectMocks = subjectMocks
		h.structFields = structFieldsResult
	}

	return h, nil
}

func validateSubjectFieldType(subjectStruct *reflect.StructField) error {
	t := subjectStruct.Type

	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		return stringerr.Newf("expected %s field to be a pointer to a struct, got %s", id.Subject, t)
	}

	return nil
}

//nolint:cyclop,nestif // Seems clearer to keep it in one method
func (t *TablePlugin) mocksForSubject(subjectStruct *reflect.StructField) (map[string]*mocks.Mock, *iterate.StructFieldsResult, error) {
	mockPaths := t.mocks.PathSet()
	subjectMocks := map[string]*mocks.Mock{}

	structFieldsResult, errs := iterate.StructFields(id.Subject, subjectStruct.Type, func(subjectFieldPath string, subjectField *reflect.StructField) []error {
		if subjectField.Type.Kind() != reflect.Interface {
			return nil
		}

		var errs []error

		for _, mock := range t.mocks.Slice() {
			if mock.Implements(subjectField.Type) {
				delete(mockPaths, mock.Path) // Even if there's an error, we don't want this mock to appear in a later error message

				if prevMatch, ok := subjectMocks[subjectFieldPath]; ok {
					if !mocks.OnlyOneRequired(mock, prevMatch) {
						err := stringerr.Newf("%s is satisfied by more than one mock: %s, %s. Exactly one required mock must match. To mark a mock optional, add the %s tag.",
							subjectFieldPath,
							prevMatch.Path,
							mock.Path,
							id.ExampleIgnoreUnused,
						)
						errs = append(errs, err)
						continue
					}

					// If the previous mock was the required one, we should skip this one
					if !prevMatch.Optional {
						continue
					}
				}

				subjectMocks[subjectFieldPath] = mock
			}
		}

		return errs
	})

	// Figure out if any required mocks were not used
	for _, mock := range t.mocks.Slice() {
		if _, ok := mockPaths[mock.Path]; !ok {
			continue
		}

		if !mock.Optional {
			err := stringerr.Newf("%s was required but not matched by any interfaces in %s. To mark a mock optional, add the %s tag.",
				mock.Path,
				id.Subject,
				id.ExampleIgnoreUnused,
			)
			errs = append(errs, err)
		}
	}

	if len(errs) != 0 {
		return nil, nil, stringerr.NewGroup(fmt.Sprintf("Unable to select mocks for %s", id.Subject), errs)
	}

	return subjectMocks, structFieldsResult, nil
}

// TableEntryHooks exposes the before and after hooks for each entry in the table.
type TableEntryHooks struct {
	plugins.NoopAfterEntry

	hasSubject   bool
	subjectMocks map[string]*mocks.Mock
	structFields *iterate.StructFieldsResult
}

var _ plugins.TableEntryHooks = &TableEntryHooks{}

// BeforeEntry is called before the test is run for the table entry.
// It initializes the Subject and fills in any matching mocks.
func (h *TableEntryHooks) BeforeEntry(ctx testctx.Context, entryValue reflect.Value, i int) error {
	if !h.hasSubject {
		return nil
	}

	subjectField := entryValue.FieldByName(id.Subject)
	h.structFields.InitializeStruct(subjectField, func(fieldPath string, field reflect.Value) {
		mock, ok := h.subjectMocks[fieldPath]
		if !ok {
			return
		}

		field.Set(mock.ValueByEntryIndex(i))
	})

	return nil
}
