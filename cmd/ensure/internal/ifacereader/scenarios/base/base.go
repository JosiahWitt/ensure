package base

import "github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader"

type ScenarioDetails struct {
	Fixture interface{}

	ExpectedMethods []*ifacereader.Method
}
