package main

import (
	"fmt"
	"log"
	"os"

	"github.com/JosiahWitt/ensure/cmd/ensure/internal/cmd"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ensurefile"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/fswrite"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/mockgen"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/mockwrite"
)

//nolint:gochecknoglobals // Allows injecting the version
// Version of the CLI.
// Should be tied to the release version.
var Version = "0.3.1"

func main() {
	if err := run(); err != nil {
		fmt.Printf("ERROR: %v\n", err) //nolint:forbidigo // Allow printing error messages
		os.Exit(1)
	}
}

func run() error {
	logger := log.New(os.Stdout, "", 0)

	mockGen, err := mockgen.New()
	if err != nil {
		return err
	}

	app := cmd.App{
		Version: Version,

		Logger:           logger,
		Getwd:            os.Getwd,
		EnsureFileLoader: &ensurefile.Loader{FS: os.DirFS("")},
		InterfaceReader:  &ifacereader.InterfaceReader{},
		MockGenerator:    mockGen,
		MockWriter: &mockwrite.MockWriter{
			FSWrite: &fswrite.FSWrite{},
			Logger:  logger,
		},
	}

	return app.Run(os.Args)
}
