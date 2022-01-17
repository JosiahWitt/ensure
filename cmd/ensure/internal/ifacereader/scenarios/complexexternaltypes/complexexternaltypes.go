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

	ExpectedMethods: []*ifacereader.Method{
		{
			Name: "InterfaceWithNestedTypes",
			Inputs: []*ifacereader.Tuple{
				{
					VariableName: "a",
					PackagePaths: []string{example1.PackagePath, example2.PackagePath},
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
					PackagePaths: []string{example1.PackagePath},
					Type:         "[]func(m *example1.Message)",
				},
				{
					VariableName: "b",
					PackagePaths: []string{example1.PackagePath, example2.PackagePath},
					Type:         "[][]map[example1.String]*example2.User",
					Variadic:     true,
				},
			},
			Outputs: []*ifacereader.Tuple{
				{
					VariableName: "x",
					PackagePaths: []string{},
					Type:         "string",
				},
				{
					VariableName: "y",
					PackagePaths: []string{example2.PackagePath},
					Type:         "[]*struct{ID example2.Float64}",
				},
				{
					VariableName: "z",
					PackagePaths: []string{},
					Type:         "error",
				},
			},
		},
	},
}
