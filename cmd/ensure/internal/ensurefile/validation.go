package ensurefile

import "github.com/JosiahWitt/erk"

const (
	defaultMockPrimaryDestination  = "internal/mocks"
	defaultMockInternalDestination = "mocks"
)

type ErkInvalidConfig struct{ erk.DefaultKind }

var (
	ErrMissingMockConfig = erk.New(ErkInvalidConfig{}, "Missing `mocks` config in .ensure.yml file. For example:\n\n"+ExampleFile)

	ErrMissingPackages = erk.New(ErkInvalidConfig{},
		"No mocks to generate. Please add some to `mocks.packages` in .ensure.yml file. For example:\n\n"+ExampleFile,
	)
	ErrDuplicatePackagePath = erk.New(ErkInvalidConfig{}, "Found duplicate package path: {{.packagePath}} in .ensure.yml file. Package paths must be unique.")
)

// validateAndSetDefaults validates the ensurefile config file and sets defaults for missing values.
func (c *Config) validateAndSetDefaults() error {
	if c.Mocks == nil {
		return ErrMissingMockConfig
	}

	if c.Mocks.PrimaryDestination == "" {
		c.Mocks.PrimaryDestination = defaultMockPrimaryDestination
	}

	if c.Mocks.InternalDestination == "" {
		c.Mocks.InternalDestination = defaultMockInternalDestination
	}

	if c.Mocks.RawTidyAfterGenerate == nil {
		tidyAfterGenerate := true
		c.Mocks.RawTidyAfterGenerate = &tidyAfterGenerate
	}
	c.Mocks.TidyAfterGenerate = *c.Mocks.RawTidyAfterGenerate

	packages := c.Mocks.Packages
	if len(packages) < 1 {
		return ErrMissingPackages
	}

	// Ensure no duplicate package paths, since the last one would overwrite the first while generating
	packagePaths := map[string]bool{}
	for _, pkg := range packages {
		if _, ok := packagePaths[pkg.Path]; ok {
			return erk.WithParams(ErrDuplicatePackagePath, erk.Params{
				"packagePath": pkg.Path,
			})
		}

		packagePaths[pkg.Path] = true
	}

	return nil
}
