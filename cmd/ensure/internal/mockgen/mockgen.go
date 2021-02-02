package mockgen

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ensurefile"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/fswrite"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/runcmd"
	"github.com/JosiahWitt/erk"
)

const (
	defaultPrimaryDestination  = "internal/mocks"
	defaultInternalDestination = "mocks"
)

type (
	ErkInvalidConfig struct{ erk.DefaultKind }
	ErkMockGenError  struct{ erk.DefaultKind }
	ErkFSWriteError  struct{ erk.DefaultKind }
)

var (
	ErrMissingMockConfig = erk.New(ErkInvalidConfig{}, "Missing `mocks` config in .ensure.yml file. For example:\n\n"+ensurefile.ExampleFile)
	ErrMissingPackages   = erk.New(ErkInvalidConfig{},
		"No mocks to generate. Please add some to `mocks.packages` in .ensure.yml file. For example:\n\n"+ensurefile.ExampleFile,
	)
	ErrDuplicatePackagePath = erk.New(ErkInvalidConfig{}, "Found duplicate package path: {{.packagePath}}. Package paths must be unique.")

	ErrMissingPackagePath       = erk.New(ErkInvalidConfig{}, "Missing `path` key for package.")
	ErrMissingPackageInterfaces = erk.New(ErkInvalidConfig{},
		"Package '{{.packagePath}}' has no interfaces to generate. Please add them using the `interfaces` key.",
	)

	ErrMockGenFailed = erk.New(ErkMockGenError{}, "Could not run mockgen successfully: {{.err}}")

	ErrUnableToCreateDir  = erk.New(ErkFSWriteError{}, "Could not create directory '{{.path}}': {{.err}}")
	ErrUnableToCreateFile = erk.New(ErkFSWriteError{}, "Could not create file '{{.path}}': {{.err}}")
)

type GeneratorIface interface {
	GenerateMocks(config *ensurefile.Config) error
}

type Generator struct {
	CmdRun  runcmd.RunnerIface
	FSWrite fswrite.FSWriteIface
}

var _ GeneratorIface = &Generator{}

// GenerateMocks for the provided configuration.
func (g *Generator) GenerateMocks(config *ensurefile.Config) error {
	if config.Mocks == nil {
		return ErrMissingMockConfig
	}

	if config.Mocks.PrimaryDestination == "" {
		config.Mocks.PrimaryDestination = defaultPrimaryDestination
	}

	if config.Mocks.InternalDestination == "" {
		config.Mocks.InternalDestination = defaultInternalDestination
	}

	packages := config.Mocks.Packages
	if len(packages) < 1 {
		return ErrMissingPackages
	}

	// Ensure no duplicate package paths, since the last one would overwrite the first
	packagePaths := map[string]bool{}
	for _, pkg := range packages {
		if _, ok := packagePaths[pkg.Path]; ok {
			return erk.WithParams(ErrDuplicatePackagePath, erk.Params{
				"packagePath": pkg.Path,
			})
		}

		packagePaths[pkg.Path] = true
	}

	fmt.Println("Generating mocks:") //nolint:forbidigo // Print header
	for _, pkg := range packages {
		if err := g.generateMock(config, pkg); err != nil {
			return err // TODO: group errors
		}
	}

	return nil
}

func (g *Generator) generateMock(config *ensurefile.Config, pkg *ensurefile.Package) error {
	//nolint:forbidigo // Print the mock currently being generated
	fmt.Printf(" - Generating: %s:%s\n", pkg.Path, strings.Join(pkg.Interfaces, ","))

	if pkg.Path == "" {
		return ErrMissingPackagePath
	}

	if len(pkg.Interfaces) < 1 {
		return erk.WithParams(ErrMissingPackageInterfaces, erk.Params{
			"packagePath": pkg.Path,
		})
	}

	mockDestination, err := computeMockDestination(config, pkg)
	if err != nil {
		return err
	}

	result, err := g.CmdRun.Exec(&runcmd.ExecParams{
		PWD: mockDestination.PWD,
		CMD: "mockgen", // TODO: Allow overriding
		Args: []string{
			pkg.Path,
			strings.Join(pkg.Interfaces, ","),
		},
	})
	if err != nil {
		return erk.WrapAs(ErrMockGenFailed, err)
	}

	result += createNEWMethods(pkg.Interfaces)

	mockFilePath := mockDestination.fullPath()
	mockDirPath := filepath.Dir(mockFilePath)

	if err := g.FSWrite.MkdirAll(mockDirPath, 0775); err != nil {
		return erk.WrapWith(ErrUnableToCreateDir, err, erk.Params{
			"path": mockDirPath,
		})
	}

	if err := g.FSWrite.WriteFile(mockFilePath, result, 0664); err != nil {
		return erk.WrapWith(ErrUnableToCreateFile, err, erk.Params{
			"path": mockFilePath,
		})
	}

	return nil
}

func createNEWMethods(interfaces []string) string {
	str := ""

	for _, iface := range interfaces {
		str += fmt.Sprintf(
			"\n// NEW creates a Mock%s.\n"+
				"func (*Mock%s) NEW(ctrl *gomock.Controller) *Mock%s {\n"+
				"\treturn NewMock%s(ctrl)\n"+
				"}\n",
			iface, iface, iface, iface,
		)
	}

	return str
}
