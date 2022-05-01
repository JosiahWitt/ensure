package iterableexternaltypes

import (
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader/fixtures/externaltypes/example1"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader/fixtures/externaltypes/example2"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader/scenarios/base"
)

type Fixture interface {
	ExternalIOSlice(a []*example1.Message) []*example2.User
	ExternalIOArray(a [2]*example1.Message) [3]*example2.User

	ExternalIOMapValue(a map[string]*example1.Message) map[int]*example2.User
	ExternalIOMapKey(a map[example2.Float64]string) map[example1.String]string
	ExternalIOMapKeyAndValue(a map[example2.Float64]*example1.Message) map[example1.String]*example2.User
	ExternalIOMapKeyAndValueSamePackage(a map[example1.String]*example1.Message) map[example2.Float64]*example2.User

	ExternalIOChan(a chan *example1.Message) chan *example2.User
}

var FixtureDetails = &base.ScenarioDetails{
	Fixture: (*Fixture)(nil),

	ExpectedPackagePaths: []string{example1.PackagePath, example2.PackagePath},

	ExpectedMethods: []*ifacereader.Method{
		{
			Name: "ExternalIOSlice",
			Inputs: []*ifacereader.Tuple{
				{VariableName: "a", Type: "[]*example1.Message"},
			},
			Outputs: []*ifacereader.Tuple{
				{VariableName: "", Type: "[]*example2.User"},
			},
		},
		{
			Name: "ExternalIOArray",
			Inputs: []*ifacereader.Tuple{
				{VariableName: "a", Type: "[2]*example1.Message"},
			},
			Outputs: []*ifacereader.Tuple{
				{VariableName: "", Type: "[3]*example2.User"},
			},
		},

		{
			Name: "ExternalIOMapValue",
			Inputs: []*ifacereader.Tuple{
				{VariableName: "a", Type: "map[string]*example1.Message"},
			},
			Outputs: []*ifacereader.Tuple{
				{VariableName: "", Type: "map[int]*example2.User"},
			},
		},
		{
			Name: "ExternalIOMapKey",
			Inputs: []*ifacereader.Tuple{
				{VariableName: "a", Type: "map[example2.Float64]string"},
			},
			Outputs: []*ifacereader.Tuple{
				{VariableName: "", Type: "map[example1.String]string"},
			},
		},
		{
			Name: "ExternalIOMapKeyAndValue",
			Inputs: []*ifacereader.Tuple{
				{VariableName: "a", Type: "map[example2.Float64]*example1.Message"},
			},
			Outputs: []*ifacereader.Tuple{
				{VariableName: "", Type: "map[example1.String]*example2.User"},
			},
		},
		{
			Name: "ExternalIOMapKeyAndValueSamePackage",
			Inputs: []*ifacereader.Tuple{
				{VariableName: "a", Type: "map[example1.String]*example1.Message"},
			},
			Outputs: []*ifacereader.Tuple{
				{VariableName: "", Type: "map[example2.Float64]*example2.User"},
			},
		},

		{
			Name: "ExternalIOChan",
			Inputs: []*ifacereader.Tuple{
				{VariableName: "a", Type: "chan *example1.Message"},
			},
			Outputs: []*ifacereader.Tuple{
				{VariableName: "", Type: "chan *example2.User"},
			},
		},
	},
}
