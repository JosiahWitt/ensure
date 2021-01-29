package cmd

import (
	"github.com/urfave/cli/v2"
)

func (a *App) generateCmd() *cli.Command {
	return &cli.Command{
		Name:  "generate",
		Usage: "prepare to run tests",

		Subcommands: []*cli.Command{
			a.generateMocksCmd(),
		},
	}
}

func (a *App) generateMocksCmd() *cli.Command {
	return &cli.Command{
		Name:  "mocks",
		Usage: "generates GoMocks (https://github.com/golang/mock) for the packages and interfaces listed in .ensure.yml",

		Action: func(c *cli.Context) error {
			pwd, err := a.Getwd()
			if err != nil {
				return err
			}

			config, err := a.EnsureFileLoader.LoadConfig(pwd)
			if err != nil {
				return err
			}

			return a.MockGenerator.GenerateMocks(config)
		},
	}
}
