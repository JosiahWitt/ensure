package single_interface_multiple_methods

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
						{VariableName: "str", Type: "string"},
					},
					Outputs: []*ifacereader.Tuple{
						{VariableName: "", Type: "string"},
						{VariableName: "", Type: "error"},
					},
				},
				{
					Name: "TransformFloat64",
					Inputs: []*ifacereader.Tuple{
						{VariableName: "f", Type: "float64"},
					},
					Outputs: []*ifacereader.Tuple{
						{VariableName: "", Type: "float64"},
					},
				},
			},
		},
	},
}
