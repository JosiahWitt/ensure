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

	ExpectedMethods: []*ifacereader.Method{
		{
			Name: "ExternalIOSlice",
			Inputs: []*ifacereader.Tuple{
				{VariableName: "a", PackagePaths: []string{example1.PackagePath}, Type: "[]*example1.Message"},
			},
			Outputs: []*ifacereader.Tuple{
				{VariableName: "", PackagePaths: []string{example2.PackagePath}, Type: "[]*example2.User"},
			},
		},
		{
			Name: "ExternalIOArray",
			Inputs: []*ifacereader.Tuple{
				{VariableName: "a", PackagePaths: []string{example1.PackagePath}, Type: "[2]*example1.Message"},
			},
			Outputs: []*ifacereader.Tuple{
				{VariableName: "", PackagePaths: []string{example2.PackagePath}, Type: "[3]*example2.User"},
			},
		},

		{
			Name: "ExternalIOMapValue",
			Inputs: []*ifacereader.Tuple{
				{VariableName: "a", PackagePaths: []string{example1.PackagePath}, Type: "map[string]*example1.Message"},
			},
			Outputs: []*ifacereader.Tuple{
				{VariableName: "", PackagePaths: []string{example2.PackagePath}, Type: "map[int]*example2.User"},
			},
		},
		{
			Name: "ExternalIOMapKey",
			Inputs: []*ifacereader.Tuple{
				{VariableName: "a", PackagePaths: []string{example2.PackagePath}, Type: "map[example2.Float64]string"},
			},
			Outputs: []*ifacereader.Tuple{
				{VariableName: "", PackagePaths: []string{example1.PackagePath}, Type: "map[example1.String]string"},
			},
		},
		{
			Name: "ExternalIOMapKeyAndValue",
			Inputs: []*ifacereader.Tuple{
				{VariableName: "a", PackagePaths: []string{example2.PackagePath, example1.PackagePath}, Type: "map[example2.Float64]*example1.Message"},
			},
			Outputs: []*ifacereader.Tuple{
				{VariableName: "", PackagePaths: []string{example1.PackagePath, example2.PackagePath}, Type: "map[example1.String]*example2.User"},
			},
		},
		{
			Name: "ExternalIOMapKeyAndValueSamePackage",
			Inputs: []*ifacereader.Tuple{
				{VariableName: "a", PackagePaths: []string{example1.PackagePath}, Type: "map[example1.String]*example1.Message"},
			},
			Outputs: []*ifacereader.Tuple{
				{VariableName: "", PackagePaths: []string{example2.PackagePath}, Type: "map[example2.Float64]*example2.User"},
			},
		},

		{
			Name: "ExternalIOChan",
			Inputs: []*ifacereader.Tuple{
				{VariableName: "a", PackagePaths: []string{example1.PackagePath}, Type: "chan *example1.Message"},
			},
			Outputs: []*ifacereader.Tuple{
				{VariableName: "", PackagePaths: []string{example2.PackagePath}, Type: "chan *example2.User"},
			},
		},
	},
}
