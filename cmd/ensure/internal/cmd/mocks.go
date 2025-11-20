package cmd

import (
	"context"

	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ensurefile"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/uniqpkg"
	"github.com/urfave/cli/v3"
)

func (a *App) mocksCmd() *cli.Command {
	return &cli.Command{
		Name:  "mocks",
		Usage: "commands related to mocks",

		Commands: []*cli.Command{
			a.mocksGenerateCmd(),
			a.mocksTidyCmd(),
		},
	}
}

func (a *App) mocksGenerateCmd() *cli.Command {
	return &cli.Command{
		Name:  "generate",
		Usage: "generates GoMocks (https://github.com/golang/mock) for the packages and interfaces listed in .ensure.yml",

		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "disable-enhanced-matcher-failures",
				Usage: "Disables the enhanced failure messages for generated mocks, reverting to standard gomock output",
			},
		},

		Action: func(ctx context.Context, cmd *cli.Command) error {
			pwd, err := a.Getwd()
			if err != nil {
				return err
			}

			config, err := a.EnsureFileLoader.LoadConfig(pwd)
			if err != nil {
				return err
			}

			if cmd.Bool("disable-enhanced-matcher-failures") {
				config.Mocks.DisableEnhancedMatcherFailures = true
			}

			pkgList := buildPackageList(config.Mocks.Packages)
			packageImports := uniqpkg.New()

			a.Logger.Println("Reading packages listed in .ensure.yml...")

			pkgs, err := a.InterfaceReader.ReadPackages(pkgList, packageImports)
			if err != nil {
				return err
			}

			a.Logger.Println("Generating mocks...")

			mocks, err := a.MockGenerator.GenerateMocks(pkgs, packageImports, config.Mocks)
			if err != nil {
				return err
			}

			if err := a.MockWriter.WriteMocks(config, mocks); err != nil {
				return err
			}

			if config.Mocks.TidyAfterGenerate {
				if err := a.MockWriter.TidyMocks(config, pkgs); err != nil {
					return err
				}
			}

			return nil
		},
	}
}

func (a *App) mocksTidyCmd() *cli.Command {
	return &cli.Command{
		Name:  "tidy",
		Usage: "removes any files and directories that would not be generated for the packages and interfaces listed in .ensure.yml",

		Action: func(ctx context.Context, cmd *cli.Command) error {
			pwd, err := a.Getwd()
			if err != nil {
				return err
			}

			config, err := a.EnsureFileLoader.LoadConfig(pwd)
			if err != nil {
				return err
			}

			pkgList := buildPackageList(config.Mocks.Packages)
			packageImports := uniqpkg.New()

			a.Logger.Println("Reading packages listed in .ensure.yml...")

			pkgs, err := a.InterfaceReader.ReadPackages(pkgList, packageImports)
			if err != nil {
				return err
			}

			return a.MockWriter.TidyMocks(config, pkgs)
		},
	}
}

func buildPackageList(configPackages []*ensurefile.MockPackage) []*ifacereader.PackageDetails {
	packages := make([]*ifacereader.PackageDetails, 0, len(configPackages))

	for _, pkg := range configPackages {
		packages = append(packages, &ifacereader.PackageDetails{
			Path:       pkg.Path,
			Interfaces: pkg.Interfaces,
		})
	}

	return packages
}
