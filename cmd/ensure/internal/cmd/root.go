package cmd

import (
	"context"
	"log"

	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ensurefile"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/mockgen"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/mockwrite"
	"github.com/urfave/cli/v3"
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
	cliApp := &cli.Command{
		Name:    "ensure",
		Usage:   "A balanced test framework for Go 1.14+.",
		Version: a.Version,

		Commands: []*cli.Command{
			a.mocksCmd(),
		},
	}

	return cliApp.Run(context.Background(), args)
}
