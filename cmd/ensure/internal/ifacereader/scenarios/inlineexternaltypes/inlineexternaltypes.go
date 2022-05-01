package inlineexternaltypes

import (
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader/fixtures/externaltypes/example1"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader/fixtures/externaltypes/example2"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader/scenarios/base"
)

type Fixture interface {
	ExternalIOFuncWithInputs(a func(m *example1.Message)) func(u *example2.User)
	ExternalIOFuncWithOutputs(a func() *example1.Message) func() *example2.User
	ExternalIOFuncWithIO(a func(m *example1.Message) *example2.User) func(f example2.Float64) example1.String
	ExternalIOFuncWithIOSamePackage(a func(m *example1.Message) example1.String) func(u *example2.User) example2.Float64

	ExternalIOStruct(a *struct{ m *example1.Message }) *struct{ m *example2.User }
	ExternalIOStructWithMultiplePackages(a *struct {
		m *example1.Message
		f example2.Float64
	}) *struct {
		m *example2.User
		s example1.String
	}
	ExternalIOStructWithSamePackage(a *struct {
		m *example1.Message
		s example1.String
	}) *struct {
		m *example2.User
		f example2.Float64
	}

	ExternalIOInterface(a interface{ Input(m *example1.Message) }) interface{ Output() *example2.User }
	ExternalIOInterfaceWithMultiplePackages(a interface {
		Method1(m *example1.Message) example2.Float64
		Method2(u *example2.User) example1.String
	}) interface {
		Method1(f example2.Float64) *example1.Message
		Method2(s example1.String) *example2.User
	}
	ExternalIOInterfaceWithSamePackage(a interface {
		IO(m *example1.Message) example1.String
	}) interface {
		IO(f example2.Float64) *example2.User
	}
}

var FixtureDetails = &base.ScenarioDetails{
	Fixture: (*Fixture)(nil),

	ExpectedPackagePaths: []string{example1.PackagePath, example2.PackagePath},

	ExpectedMethods: []*ifacereader.Method{
		{
			Name: "ExternalIOFuncWithInputs",
			Inputs: []*ifacereader.Tuple{
				{VariableName: "a", Type: "func(m *example1.Message)"},
			},
			Outputs: []*ifacereader.Tuple{
				{VariableName: "", Type: "func(u *example2.User)"},
			},
		},
		{
			Name: "ExternalIOFuncWithOutputs",
			Inputs: []*ifacereader.Tuple{
				{VariableName: "a", Type: "func() *example1.Message"},
			},
			Outputs: []*ifacereader.Tuple{
				{VariableName: "", Type: "func() *example2.User"},
			},
		},
		{
			Name: "ExternalIOFuncWithIO",
			Inputs: []*ifacereader.Tuple{
				{VariableName: "a", Type: "func(m *example1.Message) *example2.User"},
			},
			Outputs: []*ifacereader.Tuple{
				{VariableName: "", Type: "func(f example2.Float64) example1.String"},
			},
		},
		{
			Name: "ExternalIOFuncWithIOSamePackage",
			Inputs: []*ifacereader.Tuple{
				{VariableName: "a", Type: "func(m *example1.Message) example1.String"},
			},
			Outputs: []*ifacereader.Tuple{
				{VariableName: "", Type: "func(u *example2.User) example2.Float64"},
			},
		},

		{
			Name: "ExternalIOStruct",
			Inputs: []*ifacereader.Tuple{
				{VariableName: "a", Type: "*struct{m *example1.Message}"},
			},
			Outputs: []*ifacereader.Tuple{
				{VariableName: "", Type: "*struct{m *example2.User}"},
			},
		},
		{
			Name: "ExternalIOStructWithMultiplePackages",
			Inputs: []*ifacereader.Tuple{
				{VariableName: "a", Type: "*struct{m *example1.Message; f example2.Float64}"},
			},
			Outputs: []*ifacereader.Tuple{
				{VariableName: "", Type: "*struct{m *example2.User; s example1.String}"},
			},
		},
		{
			Name: "ExternalIOStructWithSamePackage",
			Inputs: []*ifacereader.Tuple{
				{VariableName: "a", Type: "*struct{m *example1.Message; s example1.String}"},
			},
			Outputs: []*ifacereader.Tuple{
				{VariableName: "", Type: "*struct{m *example2.User; f example2.Float64}"},
			},
		},

		{
			Name: "ExternalIOInterface",
			Inputs: []*ifacereader.Tuple{
				{VariableName: "a", Type: "interface{Input(m *example1.Message)}"},
			},
			Outputs: []*ifacereader.Tuple{
				{VariableName: "", Type: "interface{Output() *example2.User}"},
			},
		},
		{
			Name: "ExternalIOInterfaceWithMultiplePackages",
			Inputs: []*ifacereader.Tuple{
				{
					VariableName: "a",
					Type:         "interface{Method1(m *example1.Message) example2.Float64; Method2(u *example2.User) example1.String}",
				},
			},
			Outputs: []*ifacereader.Tuple{
				{
					VariableName: "",
					Type:         "interface{Method1(f example2.Float64) *example1.Message; Method2(s example1.String) *example2.User}",
				},
			},
		},
		{
			Name: "ExternalIOInterfaceWithSamePackage",
			Inputs: []*ifacereader.Tuple{
				{
					VariableName: "a",
					Type:         "interface{IO(m *example1.Message) example1.String}",
				},
			},
			Outputs: []*ifacereader.Tuple{
				{
					VariableName: "",
					Type:         "interface{IO(f example2.Float64) *example2.User}",
				},
			},
		},
	},
}
