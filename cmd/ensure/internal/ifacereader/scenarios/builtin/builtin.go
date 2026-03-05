package builtin

import (
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader/scenarios/base"
)

type Fixture interface {
	Bool(a bool) bool
	String(a string) string
	Int(a int) int
	Float64(a float64) float64

	EmptyInterface(a any) any
}

var FixtureDetails = &base.ScenarioDetails{
	Fixture: (*Fixture)(nil),

	ExpectedMethods: []*ifacereader.Method{
		{
			Name: "Bool",
			Inputs: []*ifacereader.Tuple{
				{VariableName: "a", Type: "bool"},
			},
			Outputs: []*ifacereader.Tuple{
				{VariableName: "", Type: "bool"},
			},
		},
		{
			Name: "String",
			Inputs: []*ifacereader.Tuple{
				{VariableName: "a", Type: "string"},
			},
			Outputs: []*ifacereader.Tuple{
				{VariableName: "", Type: "string"},
			},
		},
		{
			Name: "Int",
			Inputs: []*ifacereader.Tuple{
				{VariableName: "a", Type: "int"},
			},
			Outputs: []*ifacereader.Tuple{
				{VariableName: "", Type: "int"},
			},
		},
		{
			Name: "Float64",
			Inputs: []*ifacereader.Tuple{
				{VariableName: "a", Type: "float64"},
			},
			Outputs: []*ifacereader.Tuple{
				{VariableName: "", Type: "float64"},
			},
		},

		{
			Name: "EmptyInterface",
			Inputs: []*ifacereader.Tuple{
				{VariableName: "a", Type: "any"},
			},
			Outputs: []*ifacereader.Tuple{
				{VariableName: "", Type: "any"},
			},
		},
	},
}
