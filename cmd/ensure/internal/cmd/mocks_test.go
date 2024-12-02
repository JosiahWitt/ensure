package cmd_test

import (
	"errors"
	"io"
	"log"
	"testing"

	"github.com/JosiahWitt/ensure"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/cmd"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ensurefile"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/mockgen"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/mocks/mock_ensurefile"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/mocks/mock_ifacereader"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/mocks/mock_mockgen"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/mocks/mock_mockwrite"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/uniqpkg"
	"github.com/JosiahWitt/ensure/ensuring"
)

func TestMocksGenerate(t *testing.T) {
	ensure := ensure.New(t)

	type Mocks struct {
		EnsureFileLoader *mock_ensurefile.MockLoaderIface
		IfaceReader      *mock_ifacereader.MockReadable
		MockGen          *mock_mockgen.MockGenerator
		MockWriter       *mock_mockwrite.MockWritable
	}

	exampleError := errors.New("something went wrong")
	defaultWd := func() (string, error) {
		return "/test", nil
	}

	table := []struct {
		Name          string
		ExpectedError error
		Flags         []string

		Getwd      func() (string, error)
		Mocks      *Mocks
		SetupMocks func(*Mocks)
		Subject    *cmd.App
	}{
		{
			Name:  "with valid execution",
			Getwd: defaultWd,
			SetupMocks: func(m *Mocks) {
				m.EnsureFileLoader.EXPECT().LoadConfig("/test").Return(buildConfig(configNoop), nil)

				pkgsImports := uniqpkg.New()

				m.IfaceReader.EXPECT().
					ReadPackages(buildIfaceReaderPackagesInput(), pkgsImports).
					Return(buildIfaceReaderPackagesOutput(), nil)

				m.MockGen.EXPECT().
					GenerateMocks(buildIfaceReaderPackagesOutput(), pkgsImports).
					Return(buildGeneratedMocks(), nil)

				m.MockWriter.EXPECT().WriteMocks(buildConfig(configNoop), buildGeneratedMocks())
			},
		},
		{
			Name:  "with valid execution and tidy mocks is enabled",
			Getwd: defaultWd,
			SetupMocks: func(m *Mocks) {
				m.EnsureFileLoader.EXPECT().LoadConfig("/test").Return(buildConfig(configTidyEnabled), nil)

				pkgsImports := uniqpkg.New()

				m.IfaceReader.EXPECT().
					ReadPackages(buildIfaceReaderPackagesInput(), pkgsImports).
					Return(buildIfaceReaderPackagesOutput(), nil)

				m.MockGen.EXPECT().
					GenerateMocks(buildIfaceReaderPackagesOutput(), pkgsImports).
					Return(buildGeneratedMocks(), nil)

				m.MockWriter.EXPECT().WriteMocks(buildConfig(configTidyEnabled), buildGeneratedMocks())

				m.MockWriter.EXPECT().TidyMocks(buildConfig(configTidyEnabled), buildIfaceReaderPackagesOutput())
			},
		},
		{
			Name:  "returns error when unable to get workind directory",
			Getwd: func() (string, error) { return "", exampleError },

			ExpectedError: exampleError,
		},
		{
			Name:  "returns error when unable to load config",
			Getwd: defaultWd,

			ExpectedError: exampleError,

			SetupMocks: func(m *Mocks) {
				m.EnsureFileLoader.EXPECT().LoadConfig("/test").Return(nil, exampleError)
			},
		},
		{
			Name:  "returns error when unable to read packages",
			Getwd: defaultWd,

			ExpectedError: exampleError,

			SetupMocks: func(m *Mocks) {
				m.EnsureFileLoader.EXPECT().LoadConfig("/test").Return(buildConfig(configTidyEnabled), nil)

				pkgsImports := uniqpkg.New()

				m.IfaceReader.EXPECT().
					ReadPackages(buildIfaceReaderPackagesInput(), pkgsImports).
					Return(nil, exampleError)
			},
		},
		{
			Name:  "returns error when unable to generate mocks",
			Getwd: defaultWd,

			ExpectedError: exampleError,

			SetupMocks: func(m *Mocks) {
				m.EnsureFileLoader.EXPECT().LoadConfig("/test").Return(buildConfig(configTidyEnabled), nil)

				pkgsImports := uniqpkg.New()

				m.IfaceReader.EXPECT().
					ReadPackages(buildIfaceReaderPackagesInput(), pkgsImports).
					Return(buildIfaceReaderPackagesOutput(), nil)

				m.MockGen.EXPECT().
					GenerateMocks(buildIfaceReaderPackagesOutput(), pkgsImports).
					Return(nil, exampleError)
			},
		},
		{
			Name:  "returns error when unable to write mocks",
			Getwd: defaultWd,

			ExpectedError: exampleError,

			SetupMocks: func(m *Mocks) {
				m.EnsureFileLoader.EXPECT().LoadConfig("/test").Return(buildConfig(configNoop), nil)

				pkgsImports := uniqpkg.New()

				m.IfaceReader.EXPECT().
					ReadPackages(buildIfaceReaderPackagesInput(), pkgsImports).
					Return(buildIfaceReaderPackagesOutput(), nil)

				m.MockGen.EXPECT().
					GenerateMocks(buildIfaceReaderPackagesOutput(), pkgsImports).
					Return(buildGeneratedMocks(), nil)

				m.MockWriter.EXPECT().WriteMocks(buildConfig(configNoop), buildGeneratedMocks()).Return(exampleError)
			},
		},
		{
			Name:  "returns error when unable to tidy mocks",
			Getwd: defaultWd,

			ExpectedError: exampleError,

			SetupMocks: func(m *Mocks) {
				m.EnsureFileLoader.EXPECT().LoadConfig("/test").Return(buildConfig(configTidyEnabled), nil)

				pkgsImports := uniqpkg.New()

				m.IfaceReader.EXPECT().
					ReadPackages(buildIfaceReaderPackagesInput(), pkgsImports).
					Return(buildIfaceReaderPackagesOutput(), nil)

				m.MockGen.EXPECT().
					GenerateMocks(buildIfaceReaderPackagesOutput(), pkgsImports).
					Return(buildGeneratedMocks(), nil)

				m.MockWriter.EXPECT().WriteMocks(buildConfig(configTidyEnabled), buildGeneratedMocks())

				m.MockWriter.EXPECT().TidyMocks(buildConfig(configTidyEnabled), buildIfaceReaderPackagesOutput()).Return(exampleError)
			},
		},
	}

	ensure.RunTableByIndex(table, func(ensure ensuring.E, i int) {
		entry := table[i]
		entry.Subject.Getwd = entry.Getwd
		entry.Subject.Logger = log.New(io.Discard, "", 0)

		err := entry.Subject.Run(append([]string{"ensure", "mocks", "generate"}, entry.Flags...))
		ensure(err).IsError(entry.ExpectedError)
	})
}

func TestMocksTidy(t *testing.T) {
	ensure := ensure.New(t)

	type Mocks struct {
		EnsureFileLoader *mock_ensurefile.MockLoaderIface
		IfaceReader      *mock_ifacereader.MockReadable
		MockGen          *mock_mockgen.MockGenerator
		MockWriter       *mock_mockwrite.MockWritable
	}

	exampleError := errors.New("something went wrong")
	defaultWd := func() (string, error) {
		return "/test", nil
	}

	table := []struct {
		Name          string
		ExpectedError error

		Getwd      func() (string, error)
		Mocks      *Mocks
		SetupMocks func(*Mocks)
		Subject    *cmd.App
	}{
		{
			Name:  "with valid execution",
			Getwd: defaultWd,
			SetupMocks: func(m *Mocks) {
				m.EnsureFileLoader.EXPECT().LoadConfig("/test").Return(buildConfig(configNoop), nil)

				pkgsImports := uniqpkg.New()

				m.IfaceReader.EXPECT().
					ReadPackages(buildIfaceReaderPackagesInput(), pkgsImports).
					Return(buildIfaceReaderPackagesOutput(), nil)

				m.MockWriter.EXPECT().
					TidyMocks(buildConfig(configNoop), buildIfaceReaderPackagesOutput()).
					Return(nil)
			},
		},
		{
			Name:          "returns error when unable to load working directory",
			Getwd:         func() (string, error) { return "", exampleError },
			ExpectedError: exampleError,
		},
		{
			Name:          "returns error when unable to load config",
			Getwd:         defaultWd,
			ExpectedError: exampleError,
			SetupMocks: func(m *Mocks) {
				m.EnsureFileLoader.EXPECT().LoadConfig("/test").Return(nil, exampleError)
			},
		},
		{
			Name:          "returns error when unable to load packages",
			Getwd:         defaultWd,
			ExpectedError: exampleError,
			SetupMocks: func(m *Mocks) {
				m.EnsureFileLoader.EXPECT().LoadConfig("/test").Return(buildConfig(configNoop), nil)

				pkgsImports := uniqpkg.New()

				m.IfaceReader.EXPECT().
					ReadPackages(buildIfaceReaderPackagesInput(), pkgsImports).
					Return(buildIfaceReaderPackagesOutput(), exampleError)
			},
		},
		{
			Name:          "returns error when unable to tidy mocks",
			Getwd:         defaultWd,
			ExpectedError: exampleError,
			SetupMocks: func(m *Mocks) {
				m.EnsureFileLoader.EXPECT().LoadConfig("/test").Return(buildConfig(configNoop), nil)

				pkgsImports := uniqpkg.New()

				m.IfaceReader.EXPECT().
					ReadPackages(buildIfaceReaderPackagesInput(), pkgsImports).
					Return(buildIfaceReaderPackagesOutput(), nil)

				m.MockWriter.EXPECT().
					TidyMocks(buildConfig(configNoop), buildIfaceReaderPackagesOutput()).
					Return(exampleError)
			},
		},
	}

	ensure.RunTableByIndex(table, func(ensure ensuring.E, i int) {
		entry := table[i]
		entry.Subject.Getwd = entry.Getwd
		entry.Subject.Logger = log.New(io.Discard, "", 0)

		err := entry.Subject.Run([]string{"ensure", "mocks", "tidy"})
		ensure(err).IsError(entry.ExpectedError)
	})
}

func configNoop(config *ensurefile.Config) {}
func configTidyEnabled(config *ensurefile.Config) {
	config.Mocks.TidyAfterGenerate = true
}

func buildConfig(modify func(config *ensurefile.Config)) *ensurefile.Config {
	config := &ensurefile.Config{
		RootPath: "/some/root/path",
		Mocks: &ensurefile.MockConfig{
			Packages: []*ensurefile.MockPackage{
				{
					Path:       "pkgs/pkg1",
					Interfaces: []string{"Iface1", "Iface2"},
				},
				{
					Path:       "pkgs/pkg2",
					Interfaces: []string{"Iface3", "Iface4"},
				},
			},
		},
	}

	modify(config)
	return config
}

func buildIfaceReaderPackagesInput() []*ifacereader.PackageDetails {
	return []*ifacereader.PackageDetails{
		{
			Path:       "pkgs/pkg1",
			Interfaces: []string{"Iface1", "Iface2"},
		},
		{
			Path:       "pkgs/pkg2",
			Interfaces: []string{"Iface3", "Iface4"},
		},
	}
}

func buildIfaceReaderPackagesOutput() []*ifacereader.Package {
	return []*ifacereader.Package{
		{
			Name: "p1",
			Path: "pkgs/pkg1",
			Interfaces: []*ifacereader.Interface{
				{
					Name: "Iface1",
					Methods: []*ifacereader.Method{
						{Name: "Method"},
					},
				},
				{
					Name: "Iface2",
					Methods: []*ifacereader.Method{
						{Name: "Method"},
					},
				},
			},
		},
		{
			Name: "p2",
			Path: "pkgs/pkg2",
			Interfaces: []*ifacereader.Interface{
				{
					Name: "Iface3",
					Methods: []*ifacereader.Method{
						{Name: "Method"},
					},
				},
				{
					Name: "Iface4",
					Methods: []*ifacereader.Method{
						{Name: "Method"},
					},
				},
			},
		},
	}
}

func buildGeneratedMocks() []*mockgen.PackageMock {
	return []*mockgen.PackageMock{
		{
			Package:      buildIfaceReaderPackagesOutput()[0],
			FileContents: "pkg1 file!",
		},
		{
			Package:      buildIfaceReaderPackagesOutput()[1],
			FileContents: "pkg2 file!",
		},
	}
}
