package single_method_external_imports

import (
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/uniqpkg"
)

var Package = &ifacereader.Package{
	Name: "pkg1",
	Path: "pkgs/pkg1",
	Interfaces: []*ifacereader.Interface{
		{
			Name: "Transformable",
			Methods: []*ifacereader.Method{
				{
					Name: "Transform",
					Inputs: []*ifacereader.Tuple{
						{VariableName: "user", Type: "*external1.User"},
						{VariableName: "message", Type: "*external2.Message"},
					},
					Outputs: []*ifacereader.Tuple{
						{VariableName: "", Type: "external1.String"},
					},
				},
			},
		},
	},
}

func AddImports(imports *uniqpkg.UniquePackagePaths) *uniqpkg.UniquePackagePaths {
	pkg := imports.ForPackage(Package.Path)

	pkg.AddImport("pkgs/external1", "external1")
	pkg.AddImport("pkgs/external2", "external2")

	return imports
}
