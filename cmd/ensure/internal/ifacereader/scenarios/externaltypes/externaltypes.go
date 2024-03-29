package externaltypes

import (
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader/fixtures/externaltypes/example1"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader/fixtures/externaltypes/example2"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader/scenarios/base"
)

type Fixture interface {
	ExternalInput(a *example1.Message)
	ExternalOutput() *example1.Message
	ExternalIO(a *example1.Message) *example1.Message
	ExternalIODifferentPackages(a *example1.Message) *example2.User
}

var FixtureDetails = &base.ScenarioDetails{
	Fixture: (*Fixture)(nil),

	ExpectedPackagePaths: []string{example1.PackagePath, example2.PackagePath},

	ExpectedMethods: []*ifacereader.Method{
		{
			Name: "ExternalInput",
			Inputs: []*ifacereader.Tuple{
				{VariableName: "a", Type: "*example1.Message"},
			},
			Outputs: []*ifacereader.Tuple{},
		},
		{
			Name:   "ExternalOutput",
			Inputs: []*ifacereader.Tuple{},
			Outputs: []*ifacereader.Tuple{
				{VariableName: "", Type: "*example1.Message"},
			},
		},
		{
			Name: "ExternalIO",
			Inputs: []*ifacereader.Tuple{
				{VariableName: "a", Type: "*example1.Message"},
			},
			Outputs: []*ifacereader.Tuple{
				{VariableName: "", Type: "*example1.Message"},
			},
		},
		{
			Name: "ExternalIODifferentPackages",
			Inputs: []*ifacereader.Tuple{
				{VariableName: "a", Type: "*example1.Message"},
			},
			Outputs: []*ifacereader.Tuple{
				{VariableName: "", Type: "*example2.User"},
			},
		},
	},
}
