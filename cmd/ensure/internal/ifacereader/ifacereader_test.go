package ifacereader_test

import (
	"fmt"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"github.com/JosiahWitt/ensure"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader/fixtures/externaltypes/example1"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader/fixtures/externaltypes/example2"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader/scenarios/base"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader/scenarios/builtin"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader/scenarios/complexexternaltypes"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader/scenarios/externaltypes"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader/scenarios/inlineexternaltypes"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader/scenarios/iterableexternaltypes"
	"github.com/JosiahWitt/ensure/ensurepkg"
)

const pathPrefix = "github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader/fixtures"

type entry struct {
	Name string

	PackageDetails []*ifacereader.PackageDetails

	ExpectedPackages []*ifacereader.Package
	ExpectedError    error

	Subject *ifacereader.InterfaceReader
}

func TestReadPackages(t *testing.T) {
	ensure := ensure.New(t)

	table := []entry{
		{
			Name: "package with empty interface",

			PackageDetails: []*ifacereader.PackageDetails{
				{
					Path:       pathPrefix + "/emptyiface",
					Interfaces: []string{"Empty"},
				},
			},

			ExpectedPackages: []*ifacereader.Package{
				{
					Path: pathPrefix + "/emptyiface",
					Interfaces: []*ifacereader.Interface{
						{
							Name:    "Empty",
							Methods: []*ifacereader.Method{},
						},
					},
				},
			},
		},
		{
			Name: "package with interface with method with no inputs or outputs",

			PackageDetails: []*ifacereader.PackageDetails{
				{
					Path:       pathPrefix + "/nodepsiface",
					Interfaces: []string{"NoIO"},
				},
			},

			ExpectedPackages: []*ifacereader.Package{
				{
					Path: pathPrefix + "/nodepsiface",
					Interfaces: []*ifacereader.Interface{
						{
							Name: "NoIO",
							Methods: []*ifacereader.Method{
								{
									Name:    "Method1",
									Inputs:  []*ifacereader.Tuple{},
									Outputs: []*ifacereader.Tuple{},
								},
							},
						},
					},
				},
			},
		},
		{
			Name: "package with interface with methods with no external dependencies",

			PackageDetails: []*ifacereader.PackageDetails{
				{
					Path:       pathPrefix + "/nodepsiface",
					Interfaces: []string{"MultipleMethods"},
				},
			},

			ExpectedPackages: []*ifacereader.Package{
				{
					Path: pathPrefix + "/nodepsiface",
					Interfaces: []*ifacereader.Interface{
						{
							Name: "MultipleMethods",
							Methods: []*ifacereader.Method{
								{
									Name: "Method1",
									Inputs: []*ifacereader.Tuple{
										{VariableName: "a", PackagePaths: []string{}, Type: "string"},
									},
									Outputs: []*ifacereader.Tuple{
										{VariableName: "", PackagePaths: []string{}, Type: "string"},
									},
								},
								{
									Name: "Method2",
									Inputs: []*ifacereader.Tuple{
										{VariableName: "a", PackagePaths: []string{}, Type: "string"},
										{VariableName: "b", PackagePaths: []string{}, Type: "float64"},
									},
									Outputs: []*ifacereader.Tuple{
										{VariableName: "", PackagePaths: []string{}, Type: "string"},
										{VariableName: "", PackagePaths: []string{}, Type: "error"},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			Name: "package with multiple interfaces",

			PackageDetails: []*ifacereader.PackageDetails{
				{
					Path:       pathPrefix + "/nodepsiface",
					Interfaces: []string{"NoIO", "MultipleMethods"},
				},
			},

			ExpectedPackages: []*ifacereader.Package{
				{
					Path: pathPrefix + "/nodepsiface",
					Interfaces: []*ifacereader.Interface{
						{
							Name: "NoIO",
							Methods: []*ifacereader.Method{
								{
									Name:    "Method1",
									Inputs:  []*ifacereader.Tuple{},
									Outputs: []*ifacereader.Tuple{},
								},
							},
						},
						{
							Name: "MultipleMethods",
							Methods: []*ifacereader.Method{
								{
									Name: "Method1",
									Inputs: []*ifacereader.Tuple{
										{VariableName: "a", PackagePaths: []string{}, Type: "string"},
									},
									Outputs: []*ifacereader.Tuple{
										{VariableName: "", PackagePaths: []string{}, Type: "string"},
									},
								},
								{
									Name: "Method2",
									Inputs: []*ifacereader.Tuple{
										{VariableName: "a", PackagePaths: []string{}, Type: "string"},
										{VariableName: "b", PackagePaths: []string{}, Type: "float64"},
									},
									Outputs: []*ifacereader.Tuple{
										{VariableName: "", PackagePaths: []string{}, Type: "string"},
										{VariableName: "", PackagePaths: []string{}, Type: "error"},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			Name: "multiple packages",

			PackageDetails: []*ifacereader.PackageDetails{
				{
					Path:       pathPrefix + "/nodepsiface",
					Interfaces: []string{"NoIO"},
				},
				{
					Path:       pathPrefix + "/emptyiface",
					Interfaces: []string{"Empty"},
				},
			},

			ExpectedPackages: []*ifacereader.Package{
				{
					Path: pathPrefix + "/nodepsiface",
					Interfaces: []*ifacereader.Interface{
						{
							Name: "NoIO",
							Methods: []*ifacereader.Method{
								{
									Name:    "Method1",
									Inputs:  []*ifacereader.Tuple{},
									Outputs: []*ifacereader.Tuple{},
								},
							},
						},
					},
				},
				{
					Path: pathPrefix + "/emptyiface",
					Interfaces: []*ifacereader.Interface{
						{
							Name:    "Empty",
							Methods: []*ifacereader.Method{},
						},
					},
				},
			},
		},
		{
			Name: "interface with named outputs",

			PackageDetails: []*ifacereader.PackageDetails{
				{
					Path:       pathPrefix + "/namedoutputs",
					Interfaces: []string{"NamedOutputs"},
				},
			},

			ExpectedPackages: []*ifacereader.Package{
				{
					Path: pathPrefix + "/namedoutputs",
					Interfaces: []*ifacereader.Interface{
						{
							Name: "NamedOutputs",
							Methods: []*ifacereader.Method{
								{
									Name:   "NamedOut",
									Inputs: []*ifacereader.Tuple{},
									Outputs: []*ifacereader.Tuple{
										{VariableName: "a", PackagePaths: []string{}, Type: "string"},
										{VariableName: "b", PackagePaths: []string{}, Type: "error"},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			Name: "package with external dependencies",

			PackageDetails: []*ifacereader.PackageDetails{
				{
					Path:       pathPrefix + "/externaltypes",
					Interfaces: []string{"ExternalTypes"},
				},
			},

			ExpectedPackages: []*ifacereader.Package{
				{
					Path: pathPrefix + "/externaltypes",
					Interfaces: []*ifacereader.Interface{
						{
							Name: "ExternalTypes",
							Methods: []*ifacereader.Method{
								{
									Name: "ExternalIO",
									Inputs: []*ifacereader.Tuple{
										{VariableName: "a", PackagePaths: []string{example1.PackagePath}, Type: "*example1.Message"},
									},
									Outputs: []*ifacereader.Tuple{
										{VariableName: "", PackagePaths: []string{example2.PackagePath}, Type: "*example2.User"},
									},
								},
							},
						},
					},
				},
			},
		},

		{
			Name: "returns ErrNoInterfaces when no interfaces are provided for a path",

			PackageDetails: []*ifacereader.PackageDetails{
				{Path: "path1", Interfaces: []string{"Abc"}},
				{Path: "path2"}, // Missing interfaces
			},

			ExpectedError: ifacereader.ErrNoInterfaces,
		},
		{
			Name: "returns ErrDuplicatePath when a path is duplicated",

			PackageDetails: []*ifacereader.PackageDetails{
				{Path: "path1", Interfaces: []string{"Abc"}},
				{Path: "path2", Interfaces: []string{"Xyz"}},
				{Path: "path1", Interfaces: []string{"Qwerty"}}, // Duplicate path
			},

			ExpectedError: ifacereader.ErrDuplicatePath,
		},
		{
			Name: "returns ErrInvalidInterface when an interface doesn't exist in a package",

			PackageDetails: []*ifacereader.PackageDetails{
				{
					Path:       pathPrefix + "/emptyiface",
					Interfaces: []string{"DNE"}, // Doesn't exist
				},
			},

			ExpectedError: ifacereader.ErrInvalidInterface,
		},
		{
			Name: "returns ErrReadingPackage when package does not exist",

			PackageDetails: []*ifacereader.PackageDetails{
				{
					Path:       "dne/dne", // Doesn't exist
					Interfaces: []string{"DNE"},
				},
			},

			ExpectedError: ifacereader.ErrReadingPackage,
		},
		{
			Name: "returns ErrNotInterface when the type provided isn't an interface",

			PackageDetails: []*ifacereader.PackageDetails{
				{
					Path:       pathPrefix + "/notiface",
					Interfaces: []string{"NotInterface"},
				},
			},

			ExpectedError: ifacereader.ErrNotInterface,
		},
	}

	table = append(table, buildTypeTests(ensure)...)

	ensure.RunTableByIndex(table, func(ensure ensurepkg.Ensure, i int) {
		entry := table[i]

		pkgs, err := entry.Subject.ReadPackages(entry.PackageDetails)
		ensure(pkgs).Equals(entry.ExpectedPackages)
		ensure(err).IsError(entry.ExpectedError)
	})
}

func buildTypeTests(ensure ensurepkg.Ensure) []entry {
	allFixtures := []*base.ScenarioDetails{
		builtin.FixtureDetails,
		externaltypes.FixtureDetails,
		iterableexternaltypes.FixtureDetails,
		inlineexternaltypes.FixtureDetails,
		complexexternaltypes.FixtureDetails,
	}

	entries := make([]entry, 0, len(allFixtures))
	for _, fixture := range allFixtures {
		fixtureType := reflect.TypeOf(fixture.Fixture).Elem()
		pkgPath := fixtureType.PkgPath()
		fixtureName := fixtureType.Name()

		// Go parses interface methods sorted by name, so we sort them to match
		expectedMethods := make([]*ifacereader.Method, len(fixture.ExpectedMethods))
		copy(expectedMethods, fixture.ExpectedMethods)
		sort.Slice(expectedMethods, func(i, j int) bool {
			return expectedMethods[i].Name < expectedMethods[j].Name
		})

		entries = append(entries, entry{
			Name: fmt.Sprintf("type fixture '%s' is parsed correctly", filepath.Base(pkgPath)),

			PackageDetails: []*ifacereader.PackageDetails{
				{
					Path:       pkgPath,
					Interfaces: []string{fixtureName},
				},
			},

			ExpectedPackages: []*ifacereader.Package{
				{
					Path: pkgPath,
					Interfaces: []*ifacereader.Interface{
						{
							Name:    fixtureName,
							Methods: expectedMethods,
						},
					},
				},
			},
		})
	}

	return entries
}
