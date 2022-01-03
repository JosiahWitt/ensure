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
						{VariableName: "prefix", PackagePaths: []string{}, Type: "string"},
						{VariableName: "str", PackagePaths: []string{}, Type: "string"},
					},
					Outputs: []*ifacereader.Tuple{
						{VariableName: "", PackagePaths: []string{}, Type: "string"},
						{VariableName: "", PackagePaths: []string{}, Type: "error"},
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
						{VariableName: "i", PackagePaths: []string{}, Type: "int"},
					},
					Outputs: []*ifacereader.Tuple{
						{VariableName: "", PackagePaths: []string{}, Type: "int"},
					},
				},
				{
					Name: "TransformFloat64",
					Inputs: []*ifacereader.Tuple{
						{VariableName: "f", PackagePaths: []string{}, Type: "float64"},
					},
					Outputs: []*ifacereader.Tuple{
						{VariableName: "", PackagePaths: []string{}, Type: "float64"},
					},
				},
			},
		},
	},
}
