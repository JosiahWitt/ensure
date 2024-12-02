package ifacereader_test

import (
	"fmt"
	"go/types"
	"os"
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
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader/scenarios/variadic"
	"github.com/JosiahWitt/ensure/ensuring"
	"golang.org/x/tools/go/packages"
)

const pathPrefix = "github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader/fixtures"

type entry struct {
	Name string

	// TODO: Remove this once Go 1.18 is the lowest supported version
	WorkingDir string

	PackageDetails       []*ifacereader.PackageDetails
	PackageNameGenerator ifacereader.PackageNameGenerator

	ExpectedPackages     []*ifacereader.Package
	ExpectedPackagePaths []string
	ExpectedError        error

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
					Name: "emptyiface",
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
					Name: "nodepsiface",
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
					Name: "nodepsiface",
					Path: pathPrefix + "/nodepsiface",
					Interfaces: []*ifacereader.Interface{
						{
							Name: "MultipleMethods",
							Methods: []*ifacereader.Method{
								{
									Name: "Method1",
									Inputs: []*ifacereader.Tuple{
										{VariableName: "a", Type: "string"},
									},
									Outputs: []*ifacereader.Tuple{
										{VariableName: "", Type: "string"},
									},
								},
								{
									Name: "Method2",
									Inputs: []*ifacereader.Tuple{
										{VariableName: "a", Type: "string"},
										{VariableName: "b", Type: "float64"},
									},
									Outputs: []*ifacereader.Tuple{
										{VariableName: "", Type: "string"},
										{VariableName: "", Type: "error"},
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
					Name: "nodepsiface",
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
										{VariableName: "a", Type: "string"},
									},
									Outputs: []*ifacereader.Tuple{
										{VariableName: "", Type: "string"},
									},
								},
								{
									Name: "Method2",
									Inputs: []*ifacereader.Tuple{
										{VariableName: "a", Type: "string"},
										{VariableName: "b", Type: "float64"},
									},
									Outputs: []*ifacereader.Tuple{
										{VariableName: "", Type: "string"},
										{VariableName: "", Type: "error"},
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
					Name: "nodepsiface",
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
					Name: "emptyiface",
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
					Name: "namedoutputs",
					Path: pathPrefix + "/namedoutputs",
					Interfaces: []*ifacereader.Interface{
						{
							Name: "NamedOutputs",
							Methods: []*ifacereader.Method{
								{
									Name:   "NamedOut",
									Inputs: []*ifacereader.Tuple{},
									Outputs: []*ifacereader.Tuple{
										{VariableName: "a", Type: "string"},
										{VariableName: "b", Type: "error"},
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

			ExpectedPackagePaths: []string{example1.PackagePath, example2.PackagePath},

			ExpectedPackages: []*ifacereader.Package{
				{
					Name: "externaltypes",
					Path: pathPrefix + "/externaltypes",
					Interfaces: []*ifacereader.Interface{
						{
							Name: "ExternalTypes",
							Methods: []*ifacereader.Method{
								{
									Name: "ExternalIO",
									Inputs: []*ifacereader.Tuple{
										{VariableName: "a", Type: "map[example2.Float64]*example1.Message"},
									},
									Outputs: []*ifacereader.Tuple{
										{VariableName: "", Type: "map[example1.String]*example2.User"},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			Name: "leverages provided package name generator",

			PackageDetails: []*ifacereader.PackageDetails{
				{
					Path:       pathPrefix + "/externaltypes",
					Interfaces: []string{"ExternalTypes"},
				},
			},

			PackageNameGenerator: packageNameGenerator(func(scopePackage *packages.Package, importedPackage *types.Package) string {
				return scopePackage.Name + "_" + importedPackage.Name() + "!"
			}),

			ExpectedPackages: []*ifacereader.Package{
				{
					Name: "externaltypes",
					Path: pathPrefix + "/externaltypes",
					Interfaces: []*ifacereader.Interface{
						{
							Name: "ExternalTypes",
							Methods: []*ifacereader.Method{
								{
									Name: "ExternalIO",
									Inputs: []*ifacereader.Tuple{
										{VariableName: "a", Type: "map[externaltypes_example2!.Float64]*externaltypes_example1!.Message"},
									},
									Outputs: []*ifacereader.Tuple{
										{VariableName: "", Type: "map[externaltypes_example1!.String]*externaltypes_example2!.User"},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			Name: "package with different name than the path suffix",

			PackageDetails: []*ifacereader.PackageDetails{
				{
					Path:       pathPrefix + "/packagewithdifferentname",
					Interfaces: []string{"Interface"},
				},
			},

			ExpectedPackages: []*ifacereader.Package{
				{
					Name: "notthesamename",
					Path: pathPrefix + "/packagewithdifferentname",
					Interfaces: []*ifacereader.Interface{
						{
							Name: "Interface",
							Methods: []*ifacereader.Method{
								{
									Name:    "Method",
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

	table = append(table, buildTypeTests()...)
	table = append(table, buildGenericTests()...)

	ensure.RunTableByIndex(table, func(ensure ensuring.E, i int) {
		entry := table[i]

		withinDirectory(entry.WorkingDir, func() {
			visitedPackages := map[string]bool{}

			pkgNameGen := entry.PackageNameGenerator
			if pkgNameGen == nil {
				pkgNameGen = packageNameGenerator(func(scopePackage *packages.Package, importedPackage *types.Package) string {
					visitedPackages[importedPackage.Path()] = true
					return importedPackage.Name()
				})
			}

			pkgs, err := entry.Subject.ReadPackages(entry.PackageDetails, pkgNameGen)
			ensure(err).IsError(entry.ExpectedError)
			ensure(pkgs).Equals(entry.ExpectedPackages)
			ensure(visitedPackages).Equals(buildExpectedPackagePathsMap(entry.ExpectedPackagePaths))
		})
	})
}

func buildTypeTests() []entry {
	allFixtures := []*base.ScenarioDetails{
		builtin.FixtureDetails,
		variadic.FixtureDetails,
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

			ExpectedPackagePaths: fixture.ExpectedPackagePaths,

			ExpectedPackages: []*ifacereader.Package{
				{
					Name: filepath.Base(pkgPath),
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

func withinDirectory(workingDir string, fn func()) {
	if workingDir == "" {
		fn()
		return
	}

	pwd, err := os.Getwd()
	checkErr(err)

	defer func() { checkErr(os.Chdir(pwd)) }()

	checkErr(os.Chdir(workingDir))
	fn()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

type packageNameGenerator func(scopePackage *packages.Package, importedPackage *types.Package) string

func (pkgNameGen packageNameGenerator) GeneratePackageName(scopePackage *packages.Package, importedPackage *types.Package) string {
	return pkgNameGen(scopePackage, importedPackage)
}

func buildExpectedPackagePathsMap(expectedPackagePaths []string) map[string]bool {
	expectedPackagePathsMap := make(map[string]bool, len(expectedPackagePaths))

	for _, expectedPackagePath := range expectedPackagePaths {
		expectedPackagePathsMap[expectedPackagePath] = true
	}

	return expectedPackagePathsMap
}

func TestInterfaceNames(t *testing.T) {
	ensure := ensure.New(t)

	pkg := ifacereader.Package{
		Name: "pkg1",
		Path: "pkgs/pkg1",
		Interfaces: []*ifacereader.Interface{
			{Name: "Iface1", Methods: []*ifacereader.Method{{Name: "Method"}}},
			{Name: "Iface2", Methods: []*ifacereader.Method{{Name: "Method"}}},
		},
	}

	ensure(pkg.InterfaceNames()).Equals([]string{"Iface1", "Iface2"})
}
