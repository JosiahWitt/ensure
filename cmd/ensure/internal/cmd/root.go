package cmd

import (
	"log"

	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ensurefile"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/mockgen"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/mockwrite"
	"github.com/urfave/cli/v2"
)

// App is the CLI application for ensure.
type App struct {
	Version string

	Logger           *log.Logger
	Getwd            func() (string, error)
	EnsureFileLoader ensurefile.LoaderIface
	InterfaceReader  ifacereader.Readable
	MockGenerator    mockgen.Generator
	MockWriter       mockwrite.Writable
}

// Run the application given the os.Args array.
func (a *App) Run(args []string) error {
	cliApp := &cli.App{
		Name:    "ensure",
		Usage:   "A balanced test framework for Go 1.14+.",
		Version: a.Version,

		ExitErrHandler: func(context *cli.Context, err error) {}, // Bubble up error

		Commands: []*cli.Command{
			a.mocksCmd(),
		},
	}

	return cliApp.Run(args)
}
