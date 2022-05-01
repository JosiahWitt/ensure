package complexexternaltypes

import (
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader/fixtures/externaltypes/example1"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader/fixtures/externaltypes/example2"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader/scenarios/base"
)

type Fixture interface {
	InterfaceWithNestedTypes(a interface {
		Method([]func(m *example1.Message) map[example1.String]*example2.User) *struct{ ID []example2.Float64 }
	})

	Variadic(a []func(m *example1.Message), b ...[]map[example1.String]*example2.User) (x string, y []*struct{ ID example2.Float64 }, z error)
}

var FixtureDetails = &base.ScenarioDetails{
	Fixture: (*Fixture)(nil),

	ExpectedPackagePaths: []string{example1.PackagePath, example2.PackagePath},

	ExpectedMethods: []*ifacereader.Method{
		{
			Name: "InterfaceWithNestedTypes",
			Inputs: []*ifacereader.Tuple{
				{
					VariableName: "a",
					Type:         "interface{Method([]func(m *example1.Message) map[example1.String]*example2.User) *struct{ID []example2.Float64}}",
				},
			},
			Outputs: []*ifacereader.Tuple{},
		},
		{
			Name: "Variadic",
			Inputs: []*ifacereader.Tuple{
				{
					VariableName: "a",
					Type:         "[]func(m *example1.Message)",
				},
				{
					VariableName: "b",
					Type:         "[][]map[example1.String]*example2.User",
					Variadic:     true,
				},
			},
			Outputs: []*ifacereader.Tuple{
				{
					VariableName: "x",
					Type:         "string",
				},
				{
					VariableName: "y",
					Type:         "[]*struct{ID example2.Float64}",
				},
				{
					VariableName: "z",
					Type:         "error",
				},
			},
		},
	},
}
