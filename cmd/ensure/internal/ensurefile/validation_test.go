package ensurefile_test

import (
	"strings"

	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ensurefile"
	"github.com/JosiahWitt/ensure/ensurepkg"
)

func appendValidationEntriesToLoadConfigEntries(ensure ensurepkg.Ensure, table []loadConfigEntry) []loadConfigEntry {
	validationEntries := []struct {
		Name string

		Config string

		ExpectedConfig *ensurefile.Config
		ExpectedError  error
	}{
		{
			Name: "returns ErrMissingMockConfig when mocks are not provided",

			Config: "",

			ExpectedError: ensurefile.ErrMissingMockConfig,
		},
		{
			Name: "returns ErrMissingPackages when mock packages are nil",

			Config: `
				mocks:
					primaryDestination: primary
					internalDestination: int
			`,

			ExpectedError: ensurefile.ErrMissingPackages,
		},
		{
			Name: "returns ErrMissingPackages when mock packages are empty",

			Config: `
				mocks:
					primaryDestination: primary
					internalDestination: int
					packages: []
			`,

			ExpectedError: ensurefile.ErrMissingPackages,
		},
		{
			Name: "returns ErrDuplicatePackagePath when package path is listed twice",

			Config: `
				mocks:
					primaryDestination: primary
					internalDestination: int
					packages:
						- path: path/1
							interfaces: [iface1, iface2]
						- path: path/2
							interfaces: [iface3, iface4]
						- path: path/1
							interfaces: [iface5, iface6]
			`,

			ExpectedError: ensurefile.ErrDuplicatePackagePath,
		},
		{
			Name: "returns no error when package paths are listed once",

			Config: `
				mocks:
					primaryDestination: primary
					internalDestination: int
					tidyAfterGenerate: false
					packages:
						- path: path/1
							interfaces: [iface1, iface2]
						- path: path/2
							interfaces: [iface3, iface4]
						- path: path/3
							interfaces: [iface5, iface6]
			`,

			ExpectedConfig: &ensurefile.Config{
				RootPath:   defaultRootPath,
				ModulePath: defaultModulePath,

				Mocks: &ensurefile.MockConfig{
					PrimaryDestination:   "primary",
					InternalDestination:  "int",
					RawTidyAfterGenerate: boolPtr(false),
					TidyAfterGenerate:    false,

					Packages: []*ensurefile.MockPackage{
						{
							Path:       "path/1",
							Interfaces: []string{"iface1", "iface2"},
						},
						{
							Path:       "path/2",
							Interfaces: []string{"iface3", "iface4"},
						},
						{
							Path:       "path/3",
							Interfaces: []string{"iface5", "iface6"},
						},
					},
				},
			},
		},
		{
			Name: "sets primaryDestination to 'internal/mocks' when it is not set",

			Config: `
				mocks:
					# primaryDestination not set
					internalDestination: int
					tidyAfterGenerate: false
					packages:
						- path: path/1
							interfaces: [iface1, iface2]
			`,

			ExpectedConfig: &ensurefile.Config{
				RootPath:   defaultRootPath,
				ModulePath: defaultModulePath,

				Mocks: &ensurefile.MockConfig{
					PrimaryDestination:   "internal/mocks",
					InternalDestination:  "int",
					RawTidyAfterGenerate: boolPtr(false),
					TidyAfterGenerate:    false,

					Packages: []*ensurefile.MockPackage{
						{
							Path:       "path/1",
							Interfaces: []string{"iface1", "iface2"},
						},
					},
				},
			},
		},
		{
			Name: "sets internalDestination to 'mocks' when it is not set",

			Config: `
				mocks:
					primaryDestination: primary
					# internalDestination not set
					tidyAfterGenerate: false
					packages:
						- path: path/1
							interfaces: [iface1, iface2]
			`,

			ExpectedConfig: &ensurefile.Config{
				RootPath:   defaultRootPath,
				ModulePath: defaultModulePath,

				Mocks: &ensurefile.MockConfig{
					PrimaryDestination:   "primary",
					InternalDestination:  "mocks",
					RawTidyAfterGenerate: boolPtr(false),
					TidyAfterGenerate:    false,

					Packages: []*ensurefile.MockPackage{
						{
							Path:       "path/1",
							Interfaces: []string{"iface1", "iface2"},
						},
					},
				},
			},
		},
		{
			Name: "sets tidyAfterGenerate to 'true' when it is not set",

			Config: `
				mocks:
					primaryDestination: primary
					internalDestination: int
					# tidyAfterGenerate not set
					packages:
						- path: path/1
							interfaces: [iface1, iface2]
			`,

			ExpectedConfig: &ensurefile.Config{
				RootPath:   defaultRootPath,
				ModulePath: defaultModulePath,

				Mocks: &ensurefile.MockConfig{
					PrimaryDestination:   "primary",
					InternalDestination:  "int",
					RawTidyAfterGenerate: boolPtr(true),
					TidyAfterGenerate:    true,

					Packages: []*ensurefile.MockPackage{
						{
							Path:       "path/1",
							Interfaces: []string{"iface1", "iface2"},
						},
					},
				},
			},
		},
		{
			Name: "does not override tidyAfterGenerate if it is 'false'",

			Config: `
				mocks:
					primaryDestination: primary
					internalDestination: int
					tidyAfterGenerate: false
					packages:
						- path: path/1
							interfaces: [iface1, iface2]
			`,

			ExpectedConfig: &ensurefile.Config{
				RootPath:   defaultRootPath,
				ModulePath: defaultModulePath,

				Mocks: &ensurefile.MockConfig{
					PrimaryDestination:   "primary",
					InternalDestination:  "int",
					RawTidyAfterGenerate: boolPtr(false),
					TidyAfterGenerate:    false,

					Packages: []*ensurefile.MockPackage{
						{
							Path:       "path/1",
							Interfaces: []string{"iface1", "iface2"},
						},
					},
				},
			},
		},
		{
			Name: "allows tidyAfterGenerate to be set to 'true'",

			Config: `
				mocks:
					primaryDestination: primary
					internalDestination: int
					tidyAfterGenerate: true
					packages:
						- path: path/1
							interfaces: [iface1, iface2]
			`,

			ExpectedConfig: &ensurefile.Config{
				RootPath:   defaultRootPath,
				ModulePath: defaultModulePath,

				Mocks: &ensurefile.MockConfig{
					PrimaryDestination:   "primary",
					InternalDestination:  "int",
					RawTidyAfterGenerate: boolPtr(true),
					TidyAfterGenerate:    true,

					Packages: []*ensurefile.MockPackage{
						{
							Path:       "path/1",
							Interfaces: []string{"iface1", "iface2"},
						},
					},
				},
			},
		},
	}

	for _, entry := range validationEntries {
		table = append(table, loadConfigEntry{
			Name: entry.Name,
			PWD:  defaultRootPath,

			ExpectedConfig: entry.ExpectedConfig,
			ExpectedError:  entry.ExpectedError,

			SetupMocks: mapFS{
				"my/app/go.mod":      defaultGoModFile,
				"my/app/.ensure.yml": strings.ReplaceAll(entry.Config, "\t", "  "), // YAML doesn't like tabs, so replace them with double spaces
			}.setupMocks,
		})
	}

	return table
}
