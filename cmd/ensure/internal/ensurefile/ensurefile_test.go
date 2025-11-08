package ensurefile_test

import (
	"errors"
	"io/fs"
	"testing"

	"github.com/JosiahWitt/ensure"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ensurefile"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/mocks/io/mock_fs"
	"github.com/JosiahWitt/ensure/ensuring"
	"go.uber.org/mock/gomock"
)

const (
	defaultRootPath   = "/my/app"
	defaultModulePath = "github.com/my/app"
	defaultGoModFile  = "module " + defaultModulePath
)

type mocks struct {
	FS *mock_fs.MockReadFileFS
}

type loadConfigEntry struct {
	Name string
	PWD  string

	ExpectedConfig *ensurefile.Config
	ExpectedError  error

	Mocks      *mocks
	SetupMocks func(*mocks)
	Subject    *ensurefile.Loader
}

func TestLoadConfig(t *testing.T) {
	ensure := ensure.New(t)

	table := []loadConfigEntry{
		{
			Name: "with valid config in current directory",
			PWD:  defaultRootPath,
			ExpectedConfig: &ensurefile.Config{
				RootPath:   defaultRootPath,
				ModulePath: defaultModulePath,
				Mocks: &ensurefile.MockConfig{
					PrimaryDestination:   "internal/mocks",
					InternalDestination:  "mocks",
					RawTidyAfterGenerate: boolPtr(true),
					TidyAfterGenerate:    true,
					Packages: []*ensurefile.MockPackage{
						{
							Path: defaultModulePath + "/some/pkg",
							Interfaces: []string{
								"Iface1",
								"Iface2",
							},
						},
					},
				},
			},

			SetupMocks: mapFS{
				"my/app/go.mod":      defaultGoModFile,
				"my/app/.ensure.yml": ensurefile.ExampleFile,
			}.setupMocks,
		},

		{
			Name: "with valid config in parent directory",
			PWD:  "/my/app/some/nested/pkg",
			ExpectedConfig: &ensurefile.Config{
				RootPath:   defaultRootPath,
				ModulePath: defaultModulePath,
				Mocks: &ensurefile.MockConfig{
					PrimaryDestination:   "internal/mocks",
					InternalDestination:  "mocks",
					RawTidyAfterGenerate: boolPtr(true),
					TidyAfterGenerate:    true,
					Packages: []*ensurefile.MockPackage{
						{
							Path: defaultModulePath + "/some/pkg",
							Interfaces: []string{
								"Iface1",
								"Iface2",
							},
						},
					},
				},
			},

			SetupMocks: mapFS{
				"my/app/go.mod":      defaultGoModFile,
				"my/app/.ensure.yml": ensurefile.ExampleFile,
			}.setupMocks,
		},

		{
			Name:          "when missing go.mod file",
			PWD:           defaultRootPath,
			ExpectedError: ensurefile.ErrCannotFindGoModule,

			SetupMocks: mapFS{
				"my/app/.ensure.yml": ensurefile.ExampleFile,
			}.setupMocks,
		},

		{
			Name:          "when cannot open go.mod file",
			PWD:           defaultRootPath,
			ExpectedError: ensurefile.ErrCannotOpenFile,

			SetupMocks: mapFS{
				"my/app/go.mod":      fs.ErrPermission,
				"my/app/.ensure.yml": ensurefile.ExampleFile,
			}.setupMocks,
		},

		{
			Name:          "when cannot parse go.mod file",
			PWD:           defaultRootPath,
			ExpectedError: ensurefile.ErrCannotParseGoModule,

			SetupMocks: mapFS{
				"my/app/go.mod":      "something is broken",
				"my/app/.ensure.yml": ensurefile.ExampleFile,
			}.setupMocks,
		},

		{
			Name:          "when cannot find .ensure.yml file",
			PWD:           defaultRootPath,
			ExpectedError: ensurefile.ErrCannotOpenFile,

			SetupMocks: mapFS{
				"my/app/go.mod": defaultGoModFile,
			}.setupMocks,
		},

		{
			Name:          "when cannot open .ensure.yml file",
			PWD:           defaultRootPath,
			ExpectedError: ensurefile.ErrCannotOpenFile,

			SetupMocks: mapFS{
				"my/app/go.mod":      defaultGoModFile,
				"my/app/.ensure.yml": fs.ErrPermission,
			}.setupMocks,
		},

		{
			Name:          "when cannot parse .ensure.yml file",
			PWD:           defaultRootPath,
			ExpectedError: ensurefile.ErrCannotUnmarshalFile,

			SetupMocks: mapFS{
				"my/app/go.mod":      defaultGoModFile,
				"my/app/.ensure.yml": "{{{{{{ Not YAML",
			}.setupMocks,
		},
	}

	table = appendValidationEntriesToLoadConfigEntries(table)

	ensure.RunTableByIndex(table, func(ensure ensuring.E, i int) {
		entry := table[i]

		config, err := entry.Subject.LoadConfig(entry.PWD)
		ensure(err).IsError(entry.ExpectedError)
		ensure(config).Equals(entry.ExpectedConfig)
	})
}

func TestPackageString(t *testing.T) {
	ensure := ensure.New(t)

	pkg := ensurefile.MockPackage{
		Path:       "github.com/my/pkg",
		Interfaces: []string{"Iface1", "Iface2"},
	}
	ensure(pkg.String()).Equals("github.com/my/pkg:Iface1,Iface2")
}

type mapFS map[string]interface{}

func (mapFS mapFS) setupMocks(m *mocks) {
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

func boolPtr(b bool) *bool {
	return &b
}
