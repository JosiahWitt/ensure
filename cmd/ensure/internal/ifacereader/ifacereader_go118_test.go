//go:build go1.18
// +build go1.18

package ifacereader_test

import (
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader"
)

func buildGenericTests() []entry {
	return []entry{
		{
			Name: "package with interface with one generic type with no external constraints",

			WorkingDir: "fixtures/generics",

			PackageDetails: []*ifacereader.PackageDetails{
				{
					Path:       pathPrefix + "/generics/singletype",
					Interfaces: []string{"Identifier"},
				},
			},

			ExpectedPackages: []*ifacereader.Package{
				{
					Name: "singletype",
					Path: pathPrefix + "/generics/singletype",
					Interfaces: []*ifacereader.Interface{
						{
							Name:       "Identifier",
							TypeParams: []*ifacereader.TypeParam{{Name: "T", Type: "any"}},
							Methods: []*ifacereader.Method{
								{
									Name: "Identity",
									Inputs: []*ifacereader.Tuple{
										{VariableName: "in", Type: "T"},
									},
									Outputs: []*ifacereader.Tuple{
										{VariableName: "out", Type: "T"},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			Name: "package with interface with multiple generic types with no external constraints",

			WorkingDir: "fixtures/generics",

			PackageDetails: []*ifacereader.PackageDetails{
				{
					Path:       pathPrefix + "/generics/multipletypes",
					Interfaces: []string{"Thingable"},
				},
			},

			ExpectedPackages: []*ifacereader.Package{
				{
					Name: "multipletypes",
					Path: pathPrefix + "/generics/multipletypes",
					Interfaces: []*ifacereader.Interface{
						{
							Name: "Thingable",
							TypeParams: []*ifacereader.TypeParam{
								{Name: "T", Type: "any"},
								{Name: "V", Type: "any"},
							},
							Methods: []*ifacereader.Method{
								{
									Name: "Identity",
									Inputs: []*ifacereader.Tuple{
										{VariableName: "in", Type: "T"},
									},
									Outputs: []*ifacereader.Tuple{
										{VariableName: "out", Type: "T"},
									},
								},
								{
									Name: "Transform",
									Inputs: []*ifacereader.Tuple{
										{VariableName: "in", Type: "T"},
									},
									Outputs: []*ifacereader.Tuple{
										{VariableName: "out", Type: "V"},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			Name: "package with interface with multiple generic types with external constraints",

			WorkingDir: "fixtures/generics",

			PackageDetails: []*ifacereader.PackageDetails{
				{
					Path:       pathPrefix + "/generics/externalconstraints",
					Interfaces: []string{"Thingable"},
				},
			},

			ExpectedPackagePaths: []string{
				"golang.org/x/exp/constraints",
				"github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader/fixtures/generics/externalconstraints/constraints", // TODO: reference this directly once Go 1.18 is the lowest supported version
			},

			ExpectedPackages: []*ifacereader.Package{
				{
					Name: "externalconstraints",
					Path: pathPrefix + "/generics/externalconstraints",
					Interfaces: []*ifacereader.Interface{
						{
							Name: "Thingable",
							TypeParams: []*ifacereader.TypeParam{
								{Name: "T", Type: "constraints.Ordered"},
								{Name: "V", Type: "constraints.Thing"},
							},
							Methods: []*ifacereader.Method{
								{
									Name: "Identity",
									Inputs: []*ifacereader.Tuple{
										{VariableName: "in", Type: "T"},
									},
									Outputs: []*ifacereader.Tuple{
										{VariableName: "out", Type: "T"},
									},
								},
								{
									Name: "Transform",
									Inputs: []*ifacereader.Tuple{
										{VariableName: "in", Type: "T"},
									},
									Outputs: []*ifacereader.Tuple{
										{VariableName: "out", Type: "V"},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			Name: "package with interface with multiple generic types with complex usage",

			WorkingDir: "fixtures/generics",

			PackageDetails: []*ifacereader.PackageDetails{
				{
					Path:       pathPrefix + "/generics/complexconstraints",
					Interfaces: []string{"Thingable"},
				},
			},

			ExpectedPackagePaths: []string{
				"golang.org/x/exp/constraints",
				"github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader/fixtures/generics/complexconstraints",              // TODO: reference this directly once Go 1.18 is the lowest supported version
				"github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader/fixtures/generics/complexconstraints/externaltype", // TODO: reference this directly once Go 1.18 is the lowest supported version
			},

			ExpectedPackages: []*ifacereader.Package{
				{
					Name: "complexconstraints",
					Path: pathPrefix + "/generics/complexconstraints",
					Interfaces: []*ifacereader.Interface{
						{
							Name: "Thingable",
							TypeParams: []*ifacereader.TypeParam{
								{Name: "T", Type: "complexconstraints.Constraint"},
								{Name: "V", Type: "interface{~string}"},
								{Name: "Composite", Type: "*complexconstraints.Thing[T, V]"},
								{Name: "Unused", Type: "constraints.Complex"},
							},
							Methods: []*ifacereader.Method{
								{
									Name: "Crazyness",
									Inputs: []*ifacereader.Tuple{
										{VariableName: "in1", Type: "*complexconstraints.Thing[T, externaltype.MyType]"},
										{VariableName: "in2", Type: "Composite"},
									},
									Outputs: []*ifacereader.Tuple{
										{VariableName: "", Type: "T"},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
