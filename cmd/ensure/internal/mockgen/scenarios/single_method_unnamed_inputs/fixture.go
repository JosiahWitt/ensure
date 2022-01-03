package single_method_unnamed_inputs

import "github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader"

var Package = &ifacereader.Package{
	Name: "pkg1",
	Path: "pkgs/pkg1",
	Interfaces: []*ifacereader.Interface{
		{
			Name: "Transformable",
			Methods: []*ifacereader.Method{
				{
					Name: "TransformString",
					Inputs: []*ifacereader.Tuple{
						{VariableName: "", PackagePaths: []string{}, Type: "string"},
						{VariableName: "", PackagePaths: []string{}, Type: "string"},
					},
					Outputs: []*ifacereader.Tuple{
						{VariableName: "", PackagePaths: []string{}, Type: "string"},
						{VariableName: "", PackagePaths: []string{}, Type: "error"},
					},
				},
			},
		},
	},
}
