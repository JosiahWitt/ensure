package cmd_test

import (
	"errors"
	"testing"

	"github.com/JosiahWitt/ensure"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/cmd"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ensurefile"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/mocks/mock_ensurefile"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/mocks/mock_mockgen"
	"github.com/JosiahWitt/ensure/ensurepkg"
)

func TestGenerateMocks(t *testing.T) {
	ensure := ensure.New(t)

	type Mocks struct {
		EnsureFileLoader *mock_ensurefile.MockLoaderIface
		MockGen          *mock_mockgen.MockGeneratorIface
	}

	exampleError := errors.New("something went wrong")
	defaultWd := func() (string, error) {
		return "/test", nil
	}

	table := []struct {
		Name          string
		ExpectedError error

		Getwd      func() (string, error)
		Mocks      *Mocks
		SetupMocks func(*Mocks)
		Subject    *cmd.App
	}{
		{
			Name:  "with valid execution",
			Getwd: defaultWd,
			SetupMocks: func(m *Mocks) {
				m.EnsureFileLoader.EXPECT().
					LoadConfig("/test").
					Return(&ensurefile.Config{
						RootPath: "/some/root/path",
					}, nil)

				m.MockGen.EXPECT().
					GenerateMocks(&ensurefile.Config{
						RootPath: "/some/root/path",
					}).
					Return(nil)
			},
		},

		{
			Name:          "when error loading working directory",
			Getwd:         func() (string, error) { return "", exampleError },
			ExpectedError: exampleError,
		},

		{
			Name:          "when cannot load config",
			Getwd:         defaultWd,
			ExpectedError: exampleError,
			SetupMocks: func(m *Mocks) {
				m.EnsureFileLoader.EXPECT().LoadConfig("/test").Return(nil, exampleError)
			},
		},

		{
			Name:          "when cannot generate mocks",
			Getwd:         defaultWd,
			ExpectedError: exampleError,
			SetupMocks: func(m *Mocks) {
				m.EnsureFileLoader.EXPECT().
					LoadConfig("/test").
					Return(&ensurefile.Config{
						RootPath: "/some/root/path",
					}, nil)

				m.MockGen.EXPECT().
					GenerateMocks(&ensurefile.Config{
						RootPath: "/some/root/path",
					}).
					Return(exampleError)
			},
		},
	}

	ensure.RunTableByIndex(table, func(ensure ensurepkg.Ensure, i int) {
		entry := table[i]
		entry.Subject.Getwd = entry.Getwd

		err := entry.Subject.Run([]string{"ensure", "generate", "mocks"})
		ensure(err).IsError(entry.ExpectedError)
	})
}
