package single_method_named_outputs

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
						{VariableName: "prefix", PackagePaths: []string{}, Type: "string"},
						{VariableName: "str", PackagePaths: []string{}, Type: "string"},
					},
					Outputs: []*ifacereader.Tuple{
						{VariableName: "transformedStr", PackagePaths: []string{}, Type: "string"},
						{VariableName: "err", PackagePaths: []string{}, Type: "error"},
					},
				},
			},
		},
	},
}
