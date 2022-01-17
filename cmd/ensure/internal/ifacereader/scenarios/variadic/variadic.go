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
				{VariableName: "prefix", PackagePaths: []string{}, Type: "string"},
				{VariableName: "strs", PackagePaths: []string{}, Type: "[]string", Variadic: true},
			},
			Outputs: []*ifacereader.Tuple{
				{VariableName: "", PackagePaths: []string{}, Type: "[]string"},
			},
		},
	},
}
