package generics_multiple_type_params

import (
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/uniqpkg"
)

var Package = &ifacereader.Package{
	Name: "pkg1",
	Path: "pkgs/pkg1",
	Interfaces: []*ifacereader.Interface{
		{
			Name: "Thingable",
			TypeParams: []*ifacereader.TypeParam{
				{Name: "T", Type: "constraints.Complex"},
				{Name: "V", Type: "thingy.Constraint"},
			},
			Methods: []*ifacereader.Method{
				{
					Name: "Identity",
					Inputs: []*ifacereader.Tuple{
						{VariableName: "in", Type: "T"},
					},
					Outputs: []*ifacereader.Tuple{
						{VariableName: "", Type: "T"},
					},
				},
				{
					Name: "Transform",
					Inputs: []*ifacereader.Tuple{
						{VariableName: "in", Type: "T"},
					},
					Outputs: []*ifacereader.Tuple{
						{VariableName: "", Type: "V"},
					},
				},
			},
		},
	},
}

func AddImports(imports *uniqpkg.UniquePackagePaths) *uniqpkg.UniquePackagePaths {
	pkg := imports.ForPackage(Package.Path)

	pkg.AddImport("pkgs/constraints", "constraints")
	pkg.AddImport("pkgs/thingy", "thingy")

	return imports
}
