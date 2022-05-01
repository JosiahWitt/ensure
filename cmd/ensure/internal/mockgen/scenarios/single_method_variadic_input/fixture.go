package single_method_variadic_input

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
						{VariableName: "prefix", Type: "string"},
						{VariableName: "strs", Type: "[]string", Variadic: true},
					},
					Outputs: []*ifacereader.Tuple{
						{VariableName: "", Type: "string"},
						{VariableName: "", Type: "error"},
					},
				},
			},
		},
	},
}
