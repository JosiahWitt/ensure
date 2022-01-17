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

	EmptyInterface(a interface{}) interface{}
}

var FixtureDetails = &base.ScenarioDetails{
	Fixture: (*Fixture)(nil),

	ExpectedMethods: []*ifacereader.Method{
		{
			Name: "Bool",
			Inputs: []*ifacereader.Tuple{
				{VariableName: "a", PackagePaths: []string{}, Type: "bool"},
			},
			Outputs: []*ifacereader.Tuple{
				{VariableName: "", PackagePaths: []string{}, Type: "bool"},
			},
		},
		{
			Name: "String",
			Inputs: []*ifacereader.Tuple{
				{VariableName: "a", PackagePaths: []string{}, Type: "string"},
			},
			Outputs: []*ifacereader.Tuple{
				{VariableName: "", PackagePaths: []string{}, Type: "string"},
			},
		},
		{
			Name: "Int",
			Inputs: []*ifacereader.Tuple{
				{VariableName: "a", PackagePaths: []string{}, Type: "int"},
			},
			Outputs: []*ifacereader.Tuple{
				{VariableName: "", PackagePaths: []string{}, Type: "int"},
			},
		},
		{
			Name: "Float64",
			Inputs: []*ifacereader.Tuple{
				{VariableName: "a", PackagePaths: []string{}, Type: "float64"},
			},
			Outputs: []*ifacereader.Tuple{
				{VariableName: "", PackagePaths: []string{}, Type: "float64"},
			},
		},

		{
			Name: "EmptyInterface",
			Inputs: []*ifacereader.Tuple{
				{VariableName: "a", PackagePaths: []string{}, Type: "interface{}"},
			},
			Outputs: []*ifacereader.Tuple{
				{VariableName: "", PackagePaths: []string{}, Type: "interface{}"},
			},
		},
	},
}
