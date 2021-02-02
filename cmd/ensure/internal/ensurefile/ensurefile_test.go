package ensurefile_test

import (
	"errors"
	"testing"

	"bursavich.dev/fs-shim/io/fs"
	"github.com/JosiahWitt/ensure"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ensurefile"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/mocks/bursavich.dev/fs-shim/io/mock_fs"
	"github.com/JosiahWitt/ensure/ensurepkg"
	"github.com/golang/mock/gomock"
)

func TestLoadConfig(t *testing.T) {
	ensure := ensure.New(t)

	const defaultGoModFile = "module github.com/my/app"

	type Mocks struct {
		FS *mock_fs.MockReadFileFS
	}

	type mapFS map[string]interface{}
	setupMapFS := func(mapFS mapFS) func(*Mocks) {
		return func(m *Mocks) {
			m.FS.EXPECT().ReadFile(gomock.Any()).AnyTimes().
				DoAndReturn(func(name string) ([]byte, error) {
					rawData, ok := mapFS[name]
					if !ok {
						return nil, fs.ErrNotExist
					}

					switch data := rawData.(type) {
					case string:
						return []byte(data), nil
					case error:
						return nil, data
					default:
						return nil, errors.New("unknown type")
					}
				})
		}
	}

	table := []struct {
		Name string
		PWD  string

		ExpectedConfig *ensurefile.Config
		ExpectedError  error

		Mocks      *Mocks
		SetupMocks func(*Mocks)
		Subject    *ensurefile.Loader
	}{
		{
			Name: "with valid config in current directory",
			PWD:  "/my/app",
			ExpectedConfig: &ensurefile.Config{
				RootPath:   "/my/app",
				ModulePath: "github.com/my/app",
				Mocks: &ensurefile.MockConfig{
					PrimaryDestination:  "internal/mocks",
					InternalDestination: "mocks",
					Packages: []*ensurefile.Package{
						{
							Path: "github.com/my/app/some/pkg",
							Interfaces: []string{
								"Iface1",
								"Iface2",
							},
						},
					},
				},
			},

			SetupMocks: setupMapFS(mapFS{
				"my/app/go.mod":      defaultGoModFile,
				"my/app/.ensure.yml": ensurefile.ExampleFile,
			}),
		},

		{
			Name: "with valid config in parent directory",
			PWD:  "/my/app/some/nested/pkg",
			ExpectedConfig: &ensurefile.Config{
				RootPath:   "/my/app",
				ModulePath: "github.com/my/app",
				Mocks: &ensurefile.MockConfig{
					PrimaryDestination:  "internal/mocks",
					InternalDestination: "mocks",
					Packages: []*ensurefile.Package{
						{
							Path: "github.com/my/app/some/pkg",
							Interfaces: []string{
								"Iface1",
								"Iface2",
							},
						},
					},
				},
			},

			SetupMocks: setupMapFS(mapFS{
				"my/app/go.mod":      defaultGoModFile,
				"my/app/.ensure.yml": ensurefile.ExampleFile,
			}),
		},

		{
			Name:          "when missing go.mod file",
			PWD:           "/my/app",
			ExpectedError: ensurefile.ErrCannotFindGoModule,

			SetupMocks: setupMapFS(mapFS{
				"my/app/.ensure.yml": ensurefile.ExampleFile,
			}),
		},

		{
			Name:          "when cannot open go.mod file",
			PWD:           "/my/app",
			ExpectedError: ensurefile.ErrCannotOpenFile,

			SetupMocks: setupMapFS(mapFS{
				"my/app/go.mod":      fs.ErrPermission,
				"my/app/.ensure.yml": ensurefile.ExampleFile,
			}),
		},

		{
			Name:          "when cannot parse go.mod file",
			PWD:           "/my/app",
			ExpectedError: ensurefile.ErrCannotParseGoModule,

			SetupMocks: setupMapFS(mapFS{
				"my/app/go.mod":      "something is broken",
				"my/app/.ensure.yml": ensurefile.ExampleFile,
			}),
		},

		{
			Name:          "when cannot find .ensure.yml file",
			PWD:           "/my/app",
			ExpectedError: ensurefile.ErrCannotOpenFile,

			SetupMocks: setupMapFS(mapFS{
				"my/app/go.mod": defaultGoModFile,
			}),
		},

		{
			Name:          "when cannot open .ensure.yml file",
			PWD:           "/my/app",
			ExpectedError: ensurefile.ErrCannotOpenFile,

			SetupMocks: setupMapFS(mapFS{
				"my/app/go.mod":      defaultGoModFile,
				"my/app/.ensure.yml": fs.ErrPermission,
			}),
		},

		{
			Name:          "when cannot parse .ensure.yml file",
			PWD:           "/my/app",
			ExpectedError: ensurefile.ErrCannotUnmarshalFile,

			SetupMocks: setupMapFS(mapFS{
				"my/app/go.mod":      defaultGoModFile,
				"my/app/.ensure.yml": "{{{{{{ Not YAML",
			}),
		},
	}

	ensure.RunTableByIndex(table, func(ensure ensurepkg.Ensure, i int) {
		entry := table[i]

		config, err := entry.Subject.LoadConfig(entry.PWD)
		ensure(err).IsError(entry.ExpectedError)
		ensure(config).Equals(entry.ExpectedConfig)
	})
}
