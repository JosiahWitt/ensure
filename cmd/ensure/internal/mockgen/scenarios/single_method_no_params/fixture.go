package single_method_no_params

import "github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader"

var Package = &ifacereader.Package{
	Name: "noop",
	Path: "pkgs/noop",
	Interfaces: []*ifacereader.Interface{
		{
			Name: "Noopable",
			Methods: []*ifacereader.Method{
				{
					Name:    "Noop",
					Inputs:  []*ifacereader.Tuple{},
					Outputs: []*ifacereader.Tuple{},
				},
			},
		},
	},
}
