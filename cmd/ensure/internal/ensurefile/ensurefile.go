// Package ensurefile assists with loading and representing .ensure.yml files.
package ensurefile

import (
	"errors"
	"path/filepath"
	"strings"

	"bursavich.dev/fs-shim/io/fs"
	"github.com/JosiahWitt/erk"
	"golang.org/x/mod/modfile"
	"gopkg.in/yaml.v3"
)

// ExampleFile for use in error messages and CLI help menus.
const ExampleFile = `mocks:
  # Used as the directory path relative to the root of the module
  # for any interfaces that are not within internal directories.
  # Optional, defaults to "internal/mocks".
  primaryDestination: internal/mocks

  # Used as the directory path relative to internal directories within the project.
  # Optional, defaults to "mocks".
  internalDestination: mocks

  # Packages with interfaces for which to generate mocks
  packages:
    - path: github.com/my/app/some/pkg
      interfaces: [Iface1, Iface2]
`

const (
	gomodFileName  = "go.mod"
	configFileName = ".ensure.yml"
)

type ErkCannotLoadConfig struct{ erk.DefaultKind }

var (
	ErrCannotFindGoModule  = erk.New(ErkCannotLoadConfig{}, "Cannot find root go.mod file by searching parent working directories")
	ErrCannotParseGoModule = erk.New(ErkCannotLoadConfig{}, "Cannot read module path from go.mod file: {{.path}}")
	ErrCannotOpenFile      = erk.New(ErkCannotLoadConfig{}, "Cannot open the file '{{.path}}': {{.err}}")
	ErrCannotUnmarshalFile = erk.New(ErkCannotLoadConfig{}, "Cannot parse the file '{{.path}}': {{.err}}")
)

type LoaderIface interface {
	LoadConfig(pwd string) (*Config, error)
}

// Loader allows loading the project's .ensure.yml file.
type Loader struct {
	FS fs.FS
}

var _ LoaderIface = &Loader{}

// Config is the root of the .ensure.yml file.
type Config struct {
	RootPath   string      `yaml:"-"`
	ModulePath string      `yaml:"-"`
	Mocks      *MockConfig `yaml:"mocks"`
}

type MockConfig struct {
	PrimaryDestination  string     `yaml:"primaryDestination"`
	InternalDestination string     `yaml:"internalDestination"`
	Packages            []*Package `yaml:"packages"`
}

type Package struct {
	Path       string   `yaml:"path"`
	Interfaces []string `yaml:"interfaces"`
}

// LoadConfig from the .ensure.yml file that is located in pwd or a parent of pwd.
func (l *Loader) LoadConfig(pwd string) (*Config, error) {
	pwd = strings.TrimPrefix(pwd, "/")
	gomodFilePath := filepath.Join(pwd, gomodFileName)

	gomodFileData, err := fs.ReadFile(l.FS, gomodFilePath)
	if errors.Is(err, fs.ErrNotExist) {
		newPWD := filepath.Dir(pwd)
		if pwd == newPWD {
			return nil, ErrCannotFindGoModule
		}

		return l.LoadConfig(newPWD)
	}

	if err != nil {
		return nil, erk.WrapWith(ErrCannotOpenFile, err, erk.Params{
			"path": gomodFilePath,
		})
	}

	modulePath := modfile.ModulePath(gomodFileData)
	if modulePath == "" {
		return nil, erk.WrapWith(ErrCannotParseGoModule, err, erk.Params{
			"path": gomodFilePath,
		})
	}

	configFilePath := filepath.Join(pwd, configFileName)
	configFileData, err := fs.ReadFile(l.FS, configFilePath)
	if err != nil {
		return nil, erk.WrapWith(ErrCannotOpenFile, err, erk.Params{
			"path": configFilePath,
		})
	}

	config := Config{}
	if err := yaml.Unmarshal(configFileData, &config); err != nil {
		return nil, erk.WrapWith(ErrCannotUnmarshalFile, err, erk.Params{
			"path": configFilePath,
		})
	}

	config.RootPath = "/" + pwd
	config.ModulePath = modulePath
	return &config, nil
}
