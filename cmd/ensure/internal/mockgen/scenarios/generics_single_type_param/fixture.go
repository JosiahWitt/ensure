package generics_single_type_param

import (
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader"
)

var Package = &ifacereader.Package{
	Name: "pkg1",
	Path: "pkgs/pkg1",
	Interfaces: []*ifacereader.Interface{
		{
			Name:       "Identifier",
			TypeParams: []*ifacereader.TypeParam{{Name: "T", Type: "any"}},
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
			},
		},
	},
}
