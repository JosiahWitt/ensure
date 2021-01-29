package cmd

import (
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ensurefile"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/mockgen"
	"github.com/urfave/cli/v2"
)

// App is the CLI application for ensure.
type App struct {
	Version string

	Getwd            func() (string, error)
	EnsureFileLoader ensurefile.LoaderIface
	MockGenerator    mockgen.GeneratorIface
}

// Run the application given the os.Args array.
func (a *App) Run(args []string) error {
	cliApp := &cli.App{
		Name:    "ensure",
		Usage:   "A balanced test framework for Go 1.14+.",
		Version: a.Version,

		Commands: []*cli.Command{
			a.generateCmd(),
		},
	}

	return cliApp.Run(args)
}
