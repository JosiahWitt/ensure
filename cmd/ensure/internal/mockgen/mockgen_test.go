package mockgen_test

import (
	"cmp"
	"os"
	"path/filepath"
	"testing"

	"github.com/JosiahWitt/ensure"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ensurefile"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/mockgen"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/mockgen/scenarios/enhanced_matcher_failures_disabled"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/mockgen/scenarios/generics_multiple_type_params"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/mockgen/scenarios/generics_single_type_param"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/mockgen/scenarios/multiple_interfaces"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/mockgen/scenarios/single_interface_multiple_methods"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/mockgen/scenarios/single_method_external_imports"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/mockgen/scenarios/single_method_external_imports_clash_with_required"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/mockgen/scenarios/single_method_external_imports_with_aliases"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/mockgen/scenarios/single_method_multiple_params"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/mockgen/scenarios/single_method_named_outputs"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/mockgen/scenarios/single_method_no_imports"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/mockgen/scenarios/single_method_no_params"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/mockgen/scenarios/single_method_unnamed_inputs"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/mockgen/scenarios/single_method_variadic_input"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/uniqpkg"
	"github.com/JosiahWitt/ensure/ensuring"
)

func TestGenerateMocks(t *testing.T) {
	ensure := ensure.New(t)

	table := []struct {
		Name string

		InputPackages []*ifacereader.Package
		Imports       *uniqpkg.UniquePackagePaths
		Config        *ensurefile.MockConfig

		ExpectedPackageMocks []*mockgen.PackageMock
	}{
		{
			Name: "with a single method with no imports",

			InputPackages: []*ifacereader.Package{
				single_method_no_imports.Package,
			},

			ExpectedPackageMocks: []*mockgen.PackageMock{
				{
					Package: single_method_no_imports.Package,

					FileContents: readExpectationFile("single_method_no_imports", "pkg1"),
				},
			},
		},
		{
			Name: "with a single method with no inputs or outputs",

			InputPackages: []*ifacereader.Package{
				single_method_no_params.Package,
			},

			ExpectedPackageMocks: []*mockgen.PackageMock{
				{
					Package: single_method_no_params.Package,

					FileContents: readExpectationFile("single_method_no_params", "noop"),
				},
			},
		},
		{
			Name: "with a single method with multiple inputs and outputs",

			InputPackages: []*ifacereader.Package{
				single_method_multiple_params.Package,
			},

			ExpectedPackageMocks: []*mockgen.PackageMock{
				{
					Package: single_method_multiple_params.Package,

					FileContents: readExpectationFile("single_method_multiple_params", "pkg1"),
				},
			},
		},
		{
			Name: "with a single method with variadic input",

			InputPackages: []*ifacereader.Package{
				single_method_variadic_input.Package,
			},

			ExpectedPackageMocks: []*mockgen.PackageMock{
				{
					Package: single_method_variadic_input.Package,

					FileContents: readExpectationFile("single_method_variadic_input", "pkg1"),
				},
			},
		},
		{
			Name: "with multiple methods within an interface",

			InputPackages: []*ifacereader.Package{
				single_interface_multiple_methods.Package,
			},

			ExpectedPackageMocks: []*mockgen.PackageMock{
				{
					Package: single_interface_multiple_methods.Package,

					FileContents: readExpectationFile("single_interface_multiple_methods", "pkg1"),
				},
			},
		},
		{
			Name: "with a single package with multiple interfaces",

			InputPackages: []*ifacereader.Package{
				multiple_interfaces.Package,
			},

			ExpectedPackageMocks: []*mockgen.PackageMock{
				{
					Package: multiple_interfaces.Package,

					FileContents: readExpectationFile("multiple_interfaces", "pkg1"),
				},
			},
		},
		{
			Name: "with a single method with unnamed inputs",

			InputPackages: []*ifacereader.Package{
				single_method_unnamed_inputs.Package,
			},

			ExpectedPackageMocks: []*mockgen.PackageMock{
				{
					Package: single_method_unnamed_inputs.Package,

					FileContents: readExpectationFile("single_method_unnamed_inputs", "pkg1"),
				},
			},
		},
		{
			Name: "with a single method with named outputs",

			InputPackages: []*ifacereader.Package{
				single_method_named_outputs.Package,
			},

			ExpectedPackageMocks: []*mockgen.PackageMock{
				{
					Package: single_method_named_outputs.Package,

					FileContents: readExpectationFile("single_method_named_outputs", "pkg1"),
				},
			},
		},
		{
			Name: "with multiple packages",

			InputPackages: []*ifacereader.Package{
				single_method_no_imports.Package,
				single_method_multiple_params.Package,
			},

			ExpectedPackageMocks: []*mockgen.PackageMock{
				{
					Package: single_method_no_imports.Package,

					FileContents: readExpectationFile("single_method_no_imports", "pkg1"),
				},
				{
					Package: single_method_multiple_params.Package,

					FileContents: readExpectationFile("single_method_multiple_params", "pkg1"),
				},
			},
		},
		{
			Name: "with a single method with imports",

			InputPackages: []*ifacereader.Package{
				single_method_external_imports.Package,
			},

			Imports: single_method_external_imports.AddImports(&uniqpkg.UniquePackagePaths{}),

			ExpectedPackageMocks: []*mockgen.PackageMock{
				{
					Package: single_method_external_imports.Package,

					FileContents: readExpectationFile("single_method_external_imports", "pkg1"),
				},
			},
		},
		{
			Name: "with a single method with imports with aliases",

			InputPackages: []*ifacereader.Package{
				single_method_external_imports_with_aliases.Package,
			},

			Imports: single_method_external_imports_with_aliases.AddImports(&uniqpkg.UniquePackagePaths{}),

			ExpectedPackageMocks: []*mockgen.PackageMock{
				{
					Package: single_method_external_imports_with_aliases.Package,

					FileContents: readExpectationFile("single_method_external_imports_with_aliases", "pkg1"),
				},
			},
		},
		{
			Name: "with a single method with imports that clash with the required imports",

			InputPackages: []*ifacereader.Package{
				single_method_external_imports_clash_with_required.Package,
			},

			Imports: single_method_external_imports_clash_with_required.AddImports(&uniqpkg.UniquePackagePaths{}),

			ExpectedPackageMocks: []*mockgen.PackageMock{
				{
					Package: single_method_external_imports_clash_with_required.Package,

					FileContents: readExpectationFile("single_method_external_imports_clash_with_required", "pkg1"),
				},
			},
		},
		{
			Name: "with a single generic type param with no imports",

			InputPackages: []*ifacereader.Package{
				generics_single_type_param.Package,
			},

			ExpectedPackageMocks: []*mockgen.PackageMock{
				{
					Package: generics_single_type_param.Package,

					FileContents: readExpectationFile("generics_single_type_param", "pkg1"),
				},
			},
		},
		{
			Name: "with multiple generic type params with imports",

			InputPackages: []*ifacereader.Package{
				generics_multiple_type_params.Package,
			},

			Imports: generics_multiple_type_params.AddImports(&uniqpkg.UniquePackagePaths{}),

			ExpectedPackageMocks: []*mockgen.PackageMock{
				{
					Package: generics_multiple_type_params.Package,

					FileContents: readExpectationFile("generics_multiple_type_params", "pkg1"),
				},
			},
		},
		{
			Name: "with enhanced matcher failures disabled",

			InputPackages: []*ifacereader.Package{
				enhanced_matcher_failures_disabled.Package,
			},
			Config: &ensurefile.MockConfig{
				DisableEnhancedMatcherFailures: true,
			},

			ExpectedPackageMocks: []*mockgen.PackageMock{
				{
					Package: enhanced_matcher_failures_disabled.Package,

					FileContents: readExpectationFile("enhanced_matcher_failures_disabled", "pkg1"),
				},
			},
		},
	}

	ensure.RunTableByIndex(table, func(ensure ensuring.E, i int) {
		entry := table[i]

		g, err := mockgen.New()
		ensure(err).IsNotError()

		imports := cmp.Or(entry.Imports, &uniqpkg.UniquePackagePaths{})
		config := cmp.Or(entry.Config, &ensurefile.MockConfig{})

		mocks, err := g.GenerateMocks(entry.InputPackages, imports, config)
		ensure(err).IsNotError()
		ensure(mocks).Equals(entry.ExpectedPackageMocks)
	})
}

func readExpectationFile(scenarioName, pkgName string) string {
	data, err := os.ReadFile(filepath.Join("scenarios", scenarioName, pkgName+".expected"))
	if err != nil {
		panic(err)
	}

	return string(data)
}
