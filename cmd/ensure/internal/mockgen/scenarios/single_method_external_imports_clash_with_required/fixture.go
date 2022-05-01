package single_method_external_imports_clash_with_required

import (
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/uniqpkg"
)

var Package = &ifacereader.Package{
	Name: "pkg1",
	Path: "pkgs/pkg1",
	Interfaces: []*ifacereader.Interface{
		{
			Name: "Doable",
			Methods: []*ifacereader.Method{
				{
					Name: "Do",
					Inputs: []*ifacereader.Tuple{
						{VariableName: "thing", Type: "*reflect.Thing"},
						{VariableName: "other", Type: "*gomock.Other"},
					},
					Outputs: []*ifacereader.Tuple{},
				},
			},
		},
	},
}

func AddImports(imports *uniqpkg.UniquePackagePaths) *uniqpkg.UniquePackagePaths {
	pkg := imports.ForPackage(Package.Path)

	// These package names will clash with the required imports "reflect" and "gomock"
	pkg.AddImport("pkgs/reflect", "reflect")
	pkg.AddImport("pkgs/gomock", "gomock")

	return imports
}
