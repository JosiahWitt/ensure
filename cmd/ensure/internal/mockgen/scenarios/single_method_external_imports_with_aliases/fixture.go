package single_method_external_imports_with_aliases

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
						{VariableName: "user", PackagePaths: []string{}, Type: "*models.User"},
						{VariableName: "message", PackagePaths: []string{}, Type: "*models2.Message"},
					},
					Outputs: []*ifacereader.Tuple{
						{VariableName: "", PackagePaths: []string{}, Type: "models.String"},
					},
				},
			},
		},
	},
}

func AddImports(imports *uniqpkg.UniquePackagePaths) *uniqpkg.UniquePackagePaths {
	pkg := imports.ForPackage(Package.Path)

	pkg.AddImport("pkgs/1/models", "models")
	pkg.AddImport("pkgs/2/models", "models")

	return imports
}
