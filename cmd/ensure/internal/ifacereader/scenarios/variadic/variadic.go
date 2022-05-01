package variadic

import (
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader/scenarios/base"
)

type Fixture interface {
	Transform(prefix string, strs ...string) []string
}

var FixtureDetails = &base.ScenarioDetails{
	Fixture: (*Fixture)(nil),

	ExpectedMethods: []*ifacereader.Method{
		{
			Name: "Transform",
			Inputs: []*ifacereader.Tuple{
				{VariableName: "prefix", Type: "string"},
				{VariableName: "strs", Type: "[]string", Variadic: true},
			},
			Outputs: []*ifacereader.Tuple{
				{VariableName: "", Type: "[]string"},
			},
		},
	},
}
