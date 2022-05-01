package multiple_interfaces

import "github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader"

var Package = &ifacereader.Package{
	Name: "pkg1",
	Path: "pkgs/pkg1",
	Interfaces: []*ifacereader.Interface{
		{
			Name: "StringTransformable",
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
			},
		},
		{
			Name: "NumberTransformable",
			Methods: []*ifacereader.Method{
				{
					Name: "TransformInt",
					Inputs: []*ifacereader.Tuple{
						{VariableName: "i", Type: "int"},
					},
					Outputs: []*ifacereader.Tuple{
						{VariableName: "", Type: "int"},
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
