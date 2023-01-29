package subject_test

import (
	"reflect"
	"testing"

	"github.com/JosiahWitt/ensure"
	"github.com/JosiahWitt/ensure/ensurepkg"
	"github.com/JosiahWitt/ensure/internal/plugins/internal/mocks"
	"github.com/JosiahWitt/ensure/internal/plugins/internal/testhelper"
	"github.com/JosiahWitt/ensure/internal/plugins/subject"
	"github.com/JosiahWitt/ensure/internal/stringerr"
)

func TestParseEntryType(t *testing.T) {
	ensure := ensure.New(t)

	table := []struct {
		Name string

		MocksInput *mocks.All
		Entry      interface{}

		ExpectedError error
	}{
		{
			Name: "returns no errors when subject is not provided",

			MocksInput: &mocks.All{},
			Entry:      struct{ Name string }{},

			ExpectedError: nil,
		},
		{
			Name: "returns error when subject is not a pointer",

			MocksInput: &mocks.All{},
			Entry: struct {
				Name    string
				Subject struct{}
			}{},

			ExpectedError: stringerr.Newf("expected Subject field to be a pointer to a struct, got struct {}"),
		},
		{
			Name: "returns error when subject is not a pointer to a struct",

			MocksInput: &mocks.All{},
			Entry: struct {
				Name    string
				Subject *string
			}{},

			ExpectedError: stringerr.Newf("expected Subject field to be a pointer to a struct, got *string"),
		},
		{
			Name: "returns no errors when no mocks are provided and subject has no fields",

			MocksInput: &mocks.All{},
			Entry: struct {
				Name    string
				Subject *struct{}
			}{},

			ExpectedError: nil,
		},
		{
			Name: "returns no errors when no mocks are provided and subject has non-exported fields",

			MocksInput: &mocks.All{},
			Entry: struct {
				Name    string
				Subject *SubjectNotExported
			}{},

			ExpectedError: nil,
		},
		{
			Name: "returns no errors when no mocks are provided and subject has non-interface fields",

			MocksInput: &mocks.All{},
			Entry: struct {
				Name    string
				Subject *SubjectNotInterface
			}{},

			ExpectedError: nil,
		},
		{
			Name: "returns no errors when no mocks are provided and subject has interface fields",

			MocksInput: &mocks.All{},
			Entry: struct {
				Name    string
				Subject *SubjectWithInterfaces
			}{},

			ExpectedError: nil,
		},
		{
			Name: "returns no errors when a mock implements interface in subject",

			MocksInput: testhelper.BuildMocks([]*testhelper.MockData{
				{
					Path: "Mocks.Bingo",
					Mock: &ExampleBingo{},
				},
			}),

			Entry: struct {
				Name    string
				Subject *SubjectWithInterfaces
			}{},

			ExpectedError: nil,
		},
		{
			Name: "returns no errors when a mock implements interface in subject's embedded field",

			MocksInput: testhelper.BuildMocks([]*testhelper.MockData{
				{
					Path: "Mocks.Bingo",
					Mock: &ExampleBingo{},
				},
			}),

			Entry: struct {
				Name    string
				Subject *SubjectWithEmbeddedField
			}{},

			ExpectedError: nil,
		},
		{
			Name: "returns no errors when a mock implements interface in subject's embedded pointer field",

			MocksInput: testhelper.BuildMocks([]*testhelper.MockData{
				{
					Path: "Mocks.Bingo",
					Mock: &ExampleBingo{},
				},
			}),

			Entry: struct {
				Name    string
				Subject *SubjectWithEmbeddedPointerField
			}{},
		},
		{
			Name: "returns no errors when a mock implements interface in subject's shadowed embedded field",

			MocksInput: testhelper.BuildMocks([]*testhelper.MockData{
				{
					Path: "Mocks.Bingo",
					Mock: &ExampleBingo{},
				},
				{
					Path: "Mocks.Hello",
					Mock: &ExampleHello{},
				},
			}),

			Entry: struct {
				Name    string
				Subject *SubjectWithShadowedEmbeddedField
			}{},

			ExpectedError: nil,
		},
		{
			Name: "returns no errors when mocks implement all interfaces in subject",

			MocksInput: testhelper.BuildMocks([]*testhelper.MockData{
				{
					Path: "Mocks.Bingo",
					Mock: &ExampleBingo{},
				},
				{
					Path: "Mocks.Hello",
					Mock: &ExampleHello{},
				},
			}),

			Entry: struct {
				Name    string
				Subject *SubjectWithInterfaces
			}{},

			ExpectedError: nil,
		},
		{
			Name: "returns no errors when a mock implements multiple interfaces in subject",

			MocksInput: testhelper.BuildMocks([]*testhelper.MockData{
				{
					Path: "Mocks.Composite",
					Mock: &ExampleComposite{},
				},
			}),

			Entry: struct {
				Name    string
				Subject *SubjectWithInterfaces
			}{},

			ExpectedError: nil,
		},
		{
			Name: "returns no errors when two mocks implement an interface in subject but only last ones are required",

			MocksInput: testhelper.BuildMocks([]*testhelper.MockData{
				{
					Path:     "Mocks.Composite",
					Optional: true,
					Mock:     &ExampleComposite{},
				},
				{
					Path: "Mocks.Bingo",
					Mock: &ExampleBingo{},
				},
				{
					Path: "Mocks.Hello",
					Mock: &ExampleHello{},
				},
			}),

			Entry: struct {
				Name    string
				Subject *SubjectWithInterfaces
			}{},

			ExpectedError: nil,
		},
		{
			Name: "returns no errors when two mocks implement an interface in subject but only first one is required",

			MocksInput: testhelper.BuildMocks([]*testhelper.MockData{
				{
					Path: "Mocks.Composite",
					Mock: &ExampleComposite{},
				},
				{
					Path:     "Mocks.Bingo",
					Optional: true,
					Mock:     &ExampleBingo{},
				},
				{
					Path:     "Mocks.Hello",
					Optional: true,
					Mock:     &ExampleHello{},
				},
			}),

			Entry: struct {
				Name    string
				Subject *SubjectWithInterfaces
			}{},

			ExpectedError: nil,
		},
		{
			Name: "returns no errors when multiple mocks implement an interface in subject but only one is required",

			MocksInput: testhelper.BuildMocks([]*testhelper.MockData{
				{
					Path:     "Mocks.Composite",
					Optional: true,
					Mock:     &ExampleComposite{},
				},
				{
					Path: "Mocks.Bingo",
					Mock: &ExampleBingo{},
				},
				{
					Path: "Mocks.Hello",
					Mock: &ExampleHello{},
				},
				{
					Path:     "Mocks.Hello2",
					Optional: true,
					Mock:     &ExampleHello{},
				},
			}),

			Entry: struct {
				Name    string
				Subject *SubjectWithInterfaces
			}{},

			ExpectedError: nil,
		},
		{
			Name: "returns error when two mocks implement an interface in subject and all are required",

			MocksInput: testhelper.BuildMocks([]*testhelper.MockData{
				{
					Path: "Mocks.Composite",
					Mock: &ExampleComposite{},
				},
				{
					Path: "Mocks.Bingo",
					Mock: &ExampleBingo{},
				},
				{
					Path: "Mocks.Hello",
					Mock: &ExampleHello{},
				},
			}),

			Entry: struct {
				Name    string
				Subject *SubjectWithInterfaces
			}{},

			ExpectedError: stringerr.Newf("Unable to select mocks for Subject:\n" +
				" - Subject.Interface is satisfied by more than one mock: Mocks.Composite, Mocks.Hello. Exactly one required mock must match. To mark a mock optional, add the `ensure:\"ignoreunused\"` tag.\n" +
				" - Subject.Yep is satisfied by more than one mock: Mocks.Composite, Mocks.Bingo. Exactly one required mock must match. To mark a mock optional, add the `ensure:\"ignoreunused\"` tag.",
			),
		},
		{
			Name: "returns error when multiple mocks implement an interface in subject and all are required",

			MocksInput: testhelper.BuildMocks([]*testhelper.MockData{
				{
					Path: "Mocks.Composite",
					Mock: &ExampleComposite{},
				},
				{
					Path: "Mocks.Bingo",
					Mock: &ExampleBingo{},
				},
				{
					Path: "Mocks.Hello",
					Mock: &ExampleHello{},
				},
				{
					Path: "Mocks.Hello2",
					Mock: &ExampleHello{},
				},
			}),

			Entry: struct {
				Name    string
				Subject *SubjectWithInterfaces
			}{},

			ExpectedError: stringerr.Newf("Unable to select mocks for Subject:\n" +
				" - Subject.Interface is satisfied by more than one mock: Mocks.Composite, Mocks.Hello. Exactly one required mock must match. To mark a mock optional, add the `ensure:\"ignoreunused\"` tag.\n" +
				" - Subject.Interface is satisfied by more than one mock: Mocks.Composite, Mocks.Hello2. Exactly one required mock must match. To mark a mock optional, add the `ensure:\"ignoreunused\"` tag.\n" +
				" - Subject.Yep is satisfied by more than one mock: Mocks.Composite, Mocks.Bingo. Exactly one required mock must match. To mark a mock optional, add the `ensure:\"ignoreunused\"` tag.",
			),
		},
		{
			Name: "returns error when multiple mocks implement an interface in subject and too many are required",

			MocksInput: testhelper.BuildMocks([]*testhelper.MockData{
				{
					Path: "Mocks.Composite",
					Mock: &ExampleComposite{},
				},
				{
					Path: "Mocks.Bingo",
					Mock: &ExampleBingo{},
				},
				{
					Path: "Mocks.Hello",
					Mock: &ExampleHello{},
				},
				{
					Path:     "Mocks.Hello2",
					Optional: true,
					Mock:     &ExampleHello{},
				},
			}),

			Entry: struct {
				Name    string
				Subject *SubjectWithInterfaces
			}{},

			ExpectedError: stringerr.Newf("Unable to select mocks for Subject:\n" +
				" - Subject.Interface is satisfied by more than one mock: Mocks.Composite, Mocks.Hello. Exactly one required mock must match. To mark a mock optional, add the `ensure:\"ignoreunused\"` tag.\n" +
				" - Subject.Yep is satisfied by more than one mock: Mocks.Composite, Mocks.Bingo. Exactly one required mock must match. To mark a mock optional, add the `ensure:\"ignoreunused\"` tag.",
			),
		},
		{
			Name: "returns error when a required mock is not matched",

			MocksInput: testhelper.BuildMocks([]*testhelper.MockData{
				{
					Path: "Mocks.Bingo",
					Mock: &ExampleBingo{},
				},
				{
					Path: "Mocks.Hello",
					Mock: &ExampleHello{},
				},
			}),

			Entry: struct {
				Name    string
				Subject *struct {
					notExported  interface{ AMethod() string }
					NotInterface *struct{}
					Interface    interface{ Hello(string) string }
				}
			}{},

			ExpectedError: stringerr.Newf("Unable to select mocks for Subject:\n" +
				" - Mocks.Bingo was required but not matched by any interfaces in Subject. To mark a mock optional, add the `ensure:\"ignoreunused\"` tag.",
			),
		},
		{
			Name: "returns error when all mocks are not matched and they are all required",

			MocksInput: testhelper.BuildMocks([]*testhelper.MockData{
				{
					Path: "Mocks.Bingo",
					Mock: &ExampleBingo{},
				},
				{
					Path: "Mocks.Hello",
					Mock: &ExampleHello{},
				},
			}),

			Entry: struct {
				Name    string
				Subject *SubjectWithNotMatchedInterfaces
			}{},

			ExpectedError: stringerr.Newf("Unable to select mocks for Subject:\n" +
				" - Mocks.Bingo was required but not matched by any interfaces in Subject. To mark a mock optional, add the `ensure:\"ignoreunused\"` tag.\n" +
				" - Mocks.Hello was required but not matched by any interfaces in Subject. To mark a mock optional, add the `ensure:\"ignoreunused\"` tag.",
			),
		},
		{
			Name: "returns error when all mocks are not matched and one is required",

			MocksInput: testhelper.BuildMocks([]*testhelper.MockData{
				{
					Path: "Mocks.Bingo",
					Mock: &ExampleBingo{},
				},
				{
					Path:     "Mocks.Hello",
					Optional: true,
					Mock:     &ExampleHello{},
				},
			}),

			Entry: struct {
				Name    string
				Subject *SubjectWithNotMatchedInterfaces
			}{},

			ExpectedError: stringerr.Newf("Unable to select mocks for Subject:\n" +
				" - Mocks.Bingo was required but not matched by any interfaces in Subject. To mark a mock optional, add the `ensure:\"ignoreunused\"` tag.",
			),
		},
		{
			Name: "returns no errors when all mocks are not matched and they are not required",

			MocksInput: testhelper.BuildMocks([]*testhelper.MockData{
				{
					Path:     "Mocks.Bingo",
					Optional: true,
					Mock:     &ExampleBingo{},
				},
				{
					Path:     "Mocks.Hello",
					Optional: true,
					Mock:     &ExampleHello{},
				},
			}),

			Entry: struct {
				Name    string
				Subject *SubjectWithNotMatchedInterfaces
			}{},

			ExpectedError: nil,
		},
		{
			Name: "returns error when all overlapping mocks are matched and none of them are required",

			MocksInput: testhelper.BuildMocks([]*testhelper.MockData{
				{
					Path:     "Mocks.Bingo",
					Mock:     &ExampleBingo{},
					Optional: true,
				},
				{
					Path:     "Mocks.Composite",
					Mock:     &ExampleComposite{},
					Optional: true,
				},
				{
					Path:     "Mocks.Hello",
					Mock:     &ExampleHello{},
					Optional: true,
				},
				{
					Path:     "Mocks.Hello2",
					Mock:     &ExampleHello{},
					Optional: true,
				},
			}),

			Entry: struct {
				Name    string
				Subject *SubjectWithInterfaces
			}{},

			ExpectedError: stringerr.Newf("Unable to select mocks for Subject:\n" +
				" - Subject.Interface is satisfied by more than one mock: Mocks.Composite, Mocks.Hello. Exactly one required mock must match. To mark a mock optional, add the `ensure:\"ignoreunused\"` tag.\n" +
				" - Subject.Interface is satisfied by more than one mock: Mocks.Composite, Mocks.Hello2. Exactly one required mock must match. To mark a mock optional, add the `ensure:\"ignoreunused\"` tag.\n" +
				" - Subject.Yep is satisfied by more than one mock: Mocks.Bingo, Mocks.Composite. Exactly one required mock must match. To mark a mock optional, add the `ensure:\"ignoreunused\"` tag.",
			),
		},
	}

	ensure.RunTableByIndex(table, func(ensure ensurepkg.Ensure, i int) {
		entry := table[i]

		plugin := subject.New(entry.MocksInput)
		res, err := plugin.ParseEntryType(reflect.TypeOf(entry.Entry))
		ensure(err).IsError(entry.ExpectedError)
		ensure(res == nil).Equals(err != nil) // res tested elsewhere
	})
}

func TestParseEntryValue(t *testing.T) {
	ensure := ensure.New(t)

	table := []struct {
		Name string

		MocksInput *mocks.All
		Table      interface{}

		ExpectedTable interface{}
	}{
		{
			Name: "makes no changes when subject is not provided",

			MocksInput: &mocks.All{},
			Table: []struct{ Name string }{
				{Name: "first"},
				{Name: "second"},
			},

			ExpectedTable: []struct{ Name string }{
				{Name: "first"},
				{Name: "second"},
			},
		},
		{
			Name: "initializes subject when no mocks are provided and subject has no fields",

			MocksInput: &mocks.All{},
			Table: []struct {
				Name    string
				Subject *struct{}
			}{
				{Name: "first"},
				{Name: "second"},
			},

			ExpectedTable: []struct {
				Name    string
				Subject *struct{}
			}{
				{Name: "first", Subject: &struct{}{}},
				{Name: "second", Subject: &struct{}{}},
			},
		},
		{
			Name: "initializes subject when no mocks are provided and subject has non-exported fields",

			MocksInput: &mocks.All{},
			Table: []struct {
				Name    string
				Subject *SubjectNotExported
			}{
				{Name: "first"},
				{Name: "second"},
			},

			ExpectedTable: []struct {
				Name    string
				Subject *SubjectNotExported
			}{
				{Name: "first", Subject: &SubjectNotExported{}},
				{Name: "second", Subject: &SubjectNotExported{}},
			},
		},
		{
			Name: "initializes subject when no mocks are provided and subject has non-interface fields",

			MocksInput: &mocks.All{},
			Table: []struct {
				Name    string
				Subject *SubjectNotInterface
			}{
				{Name: "first"},
				{Name: "second"},
			},

			ExpectedTable: []struct {
				Name    string
				Subject *SubjectNotInterface
			}{
				{Name: "first", Subject: &SubjectNotInterface{}},
				{Name: "second", Subject: &SubjectNotInterface{}},
			},
		},
		{
			Name: "initializes subject when no mocks are provided and subject has interface fields",

			MocksInput: &mocks.All{},
			Table: []struct {
				Name    string
				Subject *SubjectWithInterfaces
			}{
				{Name: "first"},
				{Name: "second"},
			},

			ExpectedTable: []struct {
				Name    string
				Subject *SubjectWithInterfaces
			}{
				{Name: "first", Subject: &SubjectWithInterfaces{}},
				{Name: "second", Subject: &SubjectWithInterfaces{}},
			},
		},
		{
			Name: "initializes subject when a mock implements interface in subject",

			MocksInput: testhelper.BuildMocks([]*testhelper.MockData{
				{
					Path: "Mocks.Bingo",
					Mock: &ExampleBingo{},
					Values: []interface{}{
						&ExampleBingo{"bingo1"},
						&ExampleBingo{"bingo2"},
					},
				},
			}),

			Table: []struct {
				Name    string
				Subject *SubjectWithInterfaces
			}{
				{Name: "first"},
				{Name: "second"},
			},

			ExpectedTable: []struct {
				Name    string
				Subject *SubjectWithInterfaces
			}{
				{
					Name: "first",
					Subject: &SubjectWithInterfaces{
						Yep: &ExampleBingo{"bingo1"},
					},
				},
				{
					Name: "second",
					Subject: &SubjectWithInterfaces{
						Yep: &ExampleBingo{"bingo2"},
					},
				},
			},
		},
		{
			Name: "initializes subject when a mock implements interface in subject's embedded field",

			MocksInput: testhelper.BuildMocks([]*testhelper.MockData{
				{
					Path: "Mocks.Bingo",
					Mock: &ExampleBingo{},
					Values: []interface{}{
						&ExampleBingo{"bingo1"},
						&ExampleBingo{"bingo2"},
					},
				},
			}),

			Table: []struct {
				Name    string
				Subject *SubjectWithEmbeddedField
			}{
				{Name: "first"},
				{Name: "second"},
			},

			ExpectedTable: []struct {
				Name    string
				Subject *SubjectWithEmbeddedField
			}{
				{
					Name: "first",
					Subject: &SubjectWithEmbeddedField{
						SubjectWithInterfaces: SubjectWithInterfaces{
							Yep: &ExampleBingo{"bingo1"},
						},
					},
				},
				{
					Name: "second",
					Subject: &SubjectWithEmbeddedField{
						SubjectWithInterfaces: SubjectWithInterfaces{
							Yep: &ExampleBingo{"bingo2"},
						},
					},
				},
			},
		},
		{
			Name: "initializes subject when a mock implements interface in subject's embedded pointer field",

			MocksInput: testhelper.BuildMocks([]*testhelper.MockData{
				{
					Path: "Mocks.Bingo",
					Mock: &ExampleBingo{},
					Values: []interface{}{
						&ExampleBingo{"bingo1"},
						&ExampleBingo{"bingo2"},
					},
				},
			}),

			Table: []struct {
				Name    string
				Subject *SubjectWithEmbeddedPointerField
			}{
				{Name: "first"},
				{Name: "second"},
			},

			ExpectedTable: []struct {
				Name    string
				Subject *SubjectWithEmbeddedPointerField
			}{
				{
					Name: "first",
					Subject: &SubjectWithEmbeddedPointerField{
						SubjectWithInterfaces: &SubjectWithInterfaces{
							Yep: &ExampleBingo{"bingo1"},
						},
					},
				},
				{
					Name: "second",
					Subject: &SubjectWithEmbeddedPointerField{
						SubjectWithInterfaces: &SubjectWithInterfaces{
							Yep: &ExampleBingo{"bingo2"},
						},
					},
				},
			},
		},
		{
			Name: "initializes subject when a mock implements interface in subject's shadowed embedded field",

			MocksInput: testhelper.BuildMocks([]*testhelper.MockData{
				{
					Path: "Mocks.Bingo",
					Mock: &ExampleBingo{},
					Values: []interface{}{
						&ExampleBingo{"bingo1"},
						&ExampleBingo{"bingo2"},
					},
				},
				{
					Path: "Mocks.Hello",
					Mock: &ExampleHello{},
					Values: []interface{}{
						&ExampleHello{"hello1"},
						&ExampleHello{"hello2"},
					},
				},
			}),

			Table: []struct {
				Name    string
				Subject *SubjectWithShadowedEmbeddedField
			}{
				{Name: "first"},
				{Name: "second"},
			},

			ExpectedTable: []struct {
				Name    string
				Subject *SubjectWithShadowedEmbeddedField
			}{
				{
					Name: "first",
					Subject: &SubjectWithShadowedEmbeddedField{
						SubjectWithInterfaces: SubjectWithInterfaces{
							Interface: &ExampleHello{"hello1"},
							Yep:       &ExampleBingo{"bingo1"},
						},
						Interface: &ExampleHello{"hello1"},
					},
				},
				{
					Name: "second",
					Subject: &SubjectWithShadowedEmbeddedField{
						SubjectWithInterfaces: SubjectWithInterfaces{
							Interface: &ExampleHello{"hello2"},
							Yep:       &ExampleBingo{"bingo2"},
						},
						Interface: &ExampleHello{"hello2"},
					},
				},
			},
		},
		{
			Name: "initializes subject when mocks implement all interfaces in subject",

			MocksInput: testhelper.BuildMocks([]*testhelper.MockData{
				{
					Path: "Mocks.Bingo",
					Mock: &ExampleBingo{},
					Values: []interface{}{
						&ExampleBingo{"bingo1"},
						&ExampleBingo{"bingo2"},
					},
				},
				{
					Path: "Mocks.Hello",
					Mock: &ExampleHello{},
					Values: []interface{}{
						&ExampleHello{"hello1"},
						&ExampleHello{"hello2"},
					},
				},
			}),

			Table: []struct {
				Name    string
				Subject *SubjectWithInterfaces
			}{
				{Name: "first"},
				{Name: "second"},
			},

			ExpectedTable: []struct {
				Name    string
				Subject *SubjectWithInterfaces
			}{
				{
					Name: "first",
					Subject: &SubjectWithInterfaces{
						Interface: &ExampleHello{"hello1"},
						Yep:       &ExampleBingo{"bingo1"},
					},
				},
				{
					Name: "second",
					Subject: &SubjectWithInterfaces{
						Interface: &ExampleHello{"hello2"},
						Yep:       &ExampleBingo{"bingo2"},
					},
				},
			},
		},
		{
			Name: "initializes subject when a mock implements multiple interfaces in subject",

			MocksInput: testhelper.BuildMocks([]*testhelper.MockData{
				{
					Path: "Mocks.Composite",
					Mock: &ExampleComposite{},
					Values: []interface{}{
						&ExampleComposite{unique: "composite1"},
						&ExampleComposite{unique: "composite2"},
					},
				},
			}),

			Table: []struct {
				Name    string
				Subject *SubjectWithInterfaces
			}{
				{Name: "first"},
				{Name: "second"},
			},

			ExpectedTable: []struct {
				Name    string
				Subject *SubjectWithInterfaces
			}{
				{
					Name: "first",
					Subject: &SubjectWithInterfaces{
						Interface: &ExampleComposite{unique: "composite1"},
						Yep:       &ExampleComposite{unique: "composite1"},
					},
				},
				{
					Name: "second",
					Subject: &SubjectWithInterfaces{
						Interface: &ExampleComposite{unique: "composite2"},
						Yep:       &ExampleComposite{unique: "composite2"},
					},
				},
			},
		},
		{
			Name: "initializes subject when two mocks implement an interface in subject but only last ones are required",

			MocksInput: testhelper.BuildMocks([]*testhelper.MockData{
				{
					Path:     "Mocks.Composite",
					Optional: true,
					Mock:     &ExampleComposite{},
				},
				{
					Path: "Mocks.Bingo",
					Mock: &ExampleBingo{},
					Values: []interface{}{
						&ExampleBingo{"bingo1"},
						&ExampleBingo{"bingo2"},
					},
				},
				{
					Path: "Mocks.Hello",
					Mock: &ExampleHello{},
					Values: []interface{}{
						&ExampleHello{"hello1"},
						&ExampleHello{"hello2"},
					},
				},
			}),

			Table: []struct {
				Name    string
				Subject *SubjectWithInterfaces
			}{
				{Name: "first"},
				{Name: "second"},
			},

			ExpectedTable: []struct {
				Name    string
				Subject *SubjectWithInterfaces
			}{
				{
					Name: "first",
					Subject: &SubjectWithInterfaces{
						Interface: &ExampleHello{"hello1"},
						Yep:       &ExampleBingo{"bingo1"},
					},
				},
				{
					Name: "second",
					Subject: &SubjectWithInterfaces{
						Interface: &ExampleHello{"hello2"},
						Yep:       &ExampleBingo{"bingo2"},
					},
				},
			},
		},
		{
			Name: "initializes subject when two mocks implement an interface in subject but only first one is required",

			MocksInput: testhelper.BuildMocks([]*testhelper.MockData{
				{
					Path: "Mocks.Composite",
					Mock: &ExampleComposite{},
					Values: []interface{}{
						&ExampleComposite{unique: "composite1"},
						&ExampleComposite{unique: "composite2"},
					},
				},
				{
					Path:     "Mocks.Bingo",
					Optional: true,
					Mock:     &ExampleBingo{},
				},
				{
					Path:     "Mocks.Hello",
					Optional: true,
					Mock:     &ExampleHello{},
				},
			}),

			Table: []struct {
				Name    string
				Subject *SubjectWithInterfaces
			}{
				{Name: "first"},
				{Name: "second"},
			},

			ExpectedTable: []struct {
				Name    string
				Subject *SubjectWithInterfaces
			}{
				{
					Name: "first",
					Subject: &SubjectWithInterfaces{
						Interface: &ExampleComposite{unique: "composite1"},
						Yep:       &ExampleComposite{unique: "composite1"},
					},
				},
				{
					Name: "second",
					Subject: &SubjectWithInterfaces{
						Interface: &ExampleComposite{unique: "composite2"},
						Yep:       &ExampleComposite{unique: "composite2"},
					},
				},
			},
		},
		{
			Name: "initializes subject when multiple mocks implement an interface in subject but only one is required",

			MocksInput: testhelper.BuildMocks([]*testhelper.MockData{
				{
					Path:     "Mocks.Composite",
					Optional: true,
					Mock:     &ExampleComposite{},
				},
				{
					Path: "Mocks.Bingo",
					Mock: &ExampleBingo{},
					Values: []interface{}{
						&ExampleBingo{"bingo1"},
						&ExampleBingo{"bingo2"},
					},
				},
				{
					Path: "Mocks.Hello",
					Mock: &ExampleHello{},
					Values: []interface{}{
						&ExampleHello{"hello1"},
						&ExampleHello{"hello2"},
					},
				},
				{
					Path:     "Mocks.Hello2",
					Optional: true,
					Mock:     &ExampleHello{},
				},
			}),

			Table: []struct {
				Name    string
				Subject *SubjectWithInterfaces
			}{
				{Name: "first"},
				{Name: "second"},
			},

			ExpectedTable: []struct {
				Name    string
				Subject *SubjectWithInterfaces
			}{
				{
					Name: "first",
					Subject: &SubjectWithInterfaces{
						Interface: &ExampleHello{"hello1"},
						Yep:       &ExampleBingo{"bingo1"},
					},
				},
				{
					Name: "second",
					Subject: &SubjectWithInterfaces{
						Interface: &ExampleHello{"hello2"},
						Yep:       &ExampleBingo{"bingo2"},
					},
				},
			},
		},
		{
			Name: "initializes subject when all mocks are not matched and they are not required",

			MocksInput: testhelper.BuildMocks([]*testhelper.MockData{
				{
					Path:     "Mocks.Bingo",
					Optional: true,
					Mock:     &ExampleBingo{},
				},
				{
					Path:     "Mocks.Hello",
					Optional: true,
					Mock:     &ExampleHello{},
				},
			}),

			Table: []struct {
				Name    string
				Subject *SubjectWithNotMatchedInterfaces
			}{
				{Name: "first"},
				{Name: "second"},
			},

			ExpectedTable: []struct {
				Name    string
				Subject *SubjectWithNotMatchedInterfaces
			}{
				{Name: "first", Subject: &SubjectWithNotMatchedInterfaces{}},
				{Name: "second", Subject: &SubjectWithNotMatchedInterfaces{}},
			},
		},
	}

	ensure.RunTableByIndex(table, func(ensure ensurepkg.Ensure, i int) {
		entry := table[i]

		plugin := subject.New(entry.MocksInput)
		tableEntryHooks, err := plugin.ParseEntryType(reflect.TypeOf(entry.Table).Elem())
		ensure(err).IsNotError()

		tableVal := reflect.ValueOf(entry.Table)
		for i := 0; i < tableVal.Len(); i++ {
			entryVal := tableVal.Index(i)

			ensure(tableEntryHooks.BeforeEntry(nil, entryVal, i)).IsNotError()
			ensure(tableEntryHooks.AfterEntry(nil, entryVal, i)).IsNotError()
		}

		ensure(entry.Table).Equals(entry.ExpectedTable)
	})
}

type (
	//lint:ignore U1000 not used for test purposes
	SubjectNotExported struct {
		notExported string
		nope        *struct{}
		notEvenThis interface{ AMethod() string }
	}

	//lint:ignore U1000 not used for test purposes
	SubjectNotInterface struct {
		notExported  interface{ AMethod() string }
		NotInterface *struct{}
		Nope         *struct{}
	}

	//lint:ignore U1000 not used for test purposes
	SubjectWithInterfaces struct {
		notExported  interface{ AMethod() string }
		NotInterface *struct{}
		Interface    interface{ Hello(string) string }
		Yep          interface{ Bingo() bool }
	}

	//lint:ignore U1000 not used for test purposes
	SubjectWithNotMatchedInterfaces struct {
		notExported  interface{ AMethod() string }
		NotInterface *struct{}
		Interface    interface{ Something() }
	}

	SubjectWithEmbeddedField struct {
		SubjectWithInterfaces
	}

	SubjectWithEmbeddedPointerField struct {
		*SubjectWithInterfaces
	}

	SubjectWithShadowedEmbeddedField struct {
		SubjectWithInterfaces
		Interface interface{ Hello(string) string } // Name and interface match embedded field
		Yep       interface{ SomethingElse() bool } // Name matches embedded field, type is different
	}
)

type ExampleBingo struct{ unique string }

func (*ExampleBingo) Bingo() bool { return false }
func (*ExampleBingo) Other()      {}

type ExampleHello struct{ unique string }

func (*ExampleHello) Hello(s string) string { return "hello, " + s }
func (*ExampleHello) Other()                {}

type ExampleComposite struct {
	ExampleBingo
	ExampleHello

	unique string
}
