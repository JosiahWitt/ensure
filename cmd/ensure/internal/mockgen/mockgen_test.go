package mockgen_test

import (
	"errors"
	"os"
	"testing"

	"github.com/JosiahWitt/ensure"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ensurefile"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/mockgen"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/mocks/mock_fswrite"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/mocks/mock_runcmd"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/runcmd"
	"github.com/JosiahWitt/ensure/ensurepkg"
	"github.com/golang/mock/gomock"
)

const (
	expectedDirPerm  = os.FileMode(0775)
	expectedFilePerm = os.FileMode(0664)
)

func TestGenerateMocks(t *testing.T) {
	ensure := ensure.New(t)

	type Mocks struct {
		CmdRun  *mock_runcmd.MockRunnerIface
		FSWrite *mock_fswrite.MockFSWriteIface
	}

	table := []struct {
		Name          string
		Config        *ensurefile.Config
		ExpectedError error

		Mocks      *Mocks
		SetupMocks func(*Mocks)
		Subject    *mockgen.Generator
	}{
		{
			Name: "with simple valid config",
			Config: &ensurefile.Config{
				RootPath:   "/root/path",
				ModulePath: "github.com/my/mod",
				Mocks: &ensurefile.MockConfig{
					PrimaryDestination:  "primary_mocks",
					InternalDestination: "internal_mocks",
					Packages: []*ensurefile.Package{
						{
							Path:       "github.com/some/pkg/abc",
							Interfaces: []string{"Iface1", "Iface2"},
						},
						{
							Path:       "github.com/some/pkg/xyz",
							Interfaces: []string{"Iface2", "Iface3"},
						},
					},
				},
			},

			SetupMocks: func(m *Mocks) {
				const expectedMockFile1 = `<abc mock stuff here>

// NEW creates a MockIface1.
func (*MockIface1) NEW(ctrl *gomock.Controller) *MockIface1 {
	return NewMockIface1(ctrl)
}

// NEW creates a MockIface2.
func (*MockIface2) NEW(ctrl *gomock.Controller) *MockIface2 {
	return NewMockIface2(ctrl)
}
`
				const expectedMockFile2 = `<xyz mock stuff here>

// NEW creates a MockIface2.
func (*MockIface2) NEW(ctrl *gomock.Controller) *MockIface2 {
	return NewMockIface2(ctrl)
}

// NEW creates a MockIface3.
func (*MockIface3) NEW(ctrl *gomock.Controller) *MockIface3 {
	return NewMockIface3(ctrl)
}
`

				gomock.InOrder(
					// Package 1

					m.CmdRun.EXPECT().Exec(&runcmd.ExecParams{
						PWD:  "/root/path",
						CMD:  "mockgen",
						Args: []string{"github.com/some/pkg/abc", "Iface1,Iface2"},
					}).Return("<abc mock stuff here>\n", nil),

					m.FSWrite.EXPECT().
						MkdirAll("/root/path/primary_mocks/github.com/some/pkg/mock_abc", expectedDirPerm).
						Return(nil),

					m.FSWrite.EXPECT().
						WriteFile(
							"/root/path/primary_mocks/github.com/some/pkg/mock_abc/mock_abc.go",
							expectedMockFile1,
							expectedFilePerm,
						).
						Return(nil),

					// Package 2

					m.CmdRun.EXPECT().Exec(&runcmd.ExecParams{
						PWD:  "/root/path",
						CMD:  "mockgen",
						Args: []string{"github.com/some/pkg/xyz", "Iface2,Iface3"},
					}).Return("<xyz mock stuff here>\n", nil),

					m.FSWrite.EXPECT().
						MkdirAll("/root/path/primary_mocks/github.com/some/pkg/mock_xyz", expectedDirPerm).
						Return(nil),

					m.FSWrite.EXPECT().
						WriteFile(
							"/root/path/primary_mocks/github.com/some/pkg/mock_xyz/mock_xyz.go",
							expectedMockFile2,
							expectedFilePerm,
						).
						Return(nil),
				)
			},
		},

		{
			Name: "with simple valid config with nested internal packages",
			Config: &ensurefile.Config{
				RootPath:   "/root/path",
				ModulePath: "github.com/my/mod",
				Mocks: &ensurefile.MockConfig{
					PrimaryDestination:  "primary_mocks",
					InternalDestination: "internal_mocks",
					Packages: []*ensurefile.Package{
						{
							Path:       "github.com/some/pkg/abc",
							Interfaces: []string{"Iface1"},
						},
						{
							Path:       "github.com/my/mod/layer1/layer2/internal/layer3/layer4/internal/layer5/layer6/xyz",
							Interfaces: []string{"Iface2"},
						},
					},
				},
			},

			SetupMocks: func(m *Mocks) {
				const expectedMockFile1 = `<abc mock stuff here>

// NEW creates a MockIface1.
func (*MockIface1) NEW(ctrl *gomock.Controller) *MockIface1 {
	return NewMockIface1(ctrl)
}
`
				const expectedMockFile2 = `<internal xyz mock stuff here>

// NEW creates a MockIface2.
func (*MockIface2) NEW(ctrl *gomock.Controller) *MockIface2 {
	return NewMockIface2(ctrl)
}
`

				gomock.InOrder(
					// Package 1

					m.CmdRun.EXPECT().Exec(&runcmd.ExecParams{
						PWD:  "/root/path",
						CMD:  "mockgen",
						Args: []string{"github.com/some/pkg/abc", "Iface1"},
					}).Return("<abc mock stuff here>\n", nil),

					m.FSWrite.EXPECT().
						MkdirAll("/root/path/primary_mocks/github.com/some/pkg/mock_abc", expectedDirPerm).
						Return(nil),

					m.FSWrite.EXPECT().
						WriteFile(
							"/root/path/primary_mocks/github.com/some/pkg/mock_abc/mock_abc.go",
							expectedMockFile1,
							expectedFilePerm,
						).
						Return(nil),

					// Package 2

					m.CmdRun.EXPECT().Exec(&runcmd.ExecParams{
						PWD:  "/root/path/layer1/layer2/internal/layer3/layer4",
						CMD:  "mockgen",
						Args: []string{"github.com/my/mod/layer1/layer2/internal/layer3/layer4/internal/layer5/layer6/xyz", "Iface2"},
					}).Return("<internal xyz mock stuff here>\n", nil),

					m.FSWrite.EXPECT().
						MkdirAll("/root/path/layer1/layer2/internal/layer3/layer4/internal/internal_mocks/layer5/layer6/mock_xyz", expectedDirPerm).
						Return(nil),

					m.FSWrite.EXPECT().
						WriteFile(
							"/root/path/layer1/layer2/internal/layer3/layer4/internal/internal_mocks/layer5/layer6/mock_xyz/mock_xyz.go",
							expectedMockFile2,
							expectedFilePerm,
						).
						Return(nil),
				)
			},
		},

		{
			Name: "with simple valid config: default primaryDestination",
			Config: &ensurefile.Config{
				RootPath:   "/root/path",
				ModulePath: "github.com/my/mod",
				Mocks: &ensurefile.MockConfig{
					InternalDestination: "internal_mocks",
					Packages: []*ensurefile.Package{
						{
							Path:       "github.com/some/pkg/abc",
							Interfaces: []string{"Iface1"},
						},
					},
				},
			},

			SetupMocks: func(m *Mocks) {
				const expectedMockFile1 = `<abc mock stuff here>

// NEW creates a MockIface1.
func (*MockIface1) NEW(ctrl *gomock.Controller) *MockIface1 {
	return NewMockIface1(ctrl)
}
`

				gomock.InOrder(
					m.CmdRun.EXPECT().Exec(&runcmd.ExecParams{
						PWD:  "/root/path",
						CMD:  "mockgen",
						Args: []string{"github.com/some/pkg/abc", "Iface1"},
					}).Return("<abc mock stuff here>\n", nil),

					m.FSWrite.EXPECT().
						MkdirAll("/root/path/internal/mocks/github.com/some/pkg/mock_abc", expectedDirPerm).
						Return(nil),

					m.FSWrite.EXPECT().
						WriteFile(
							"/root/path/internal/mocks/github.com/some/pkg/mock_abc/mock_abc.go",
							expectedMockFile1,
							expectedFilePerm,
						).
						Return(nil),
				)
			},
		},

		{
			Name: "with simple valid config: default internalDestination",
			Config: &ensurefile.Config{
				RootPath:   "/root/path",
				ModulePath: "github.com/my/mod",
				Mocks: &ensurefile.MockConfig{
					PrimaryDestination: "primary_mocks",
					Packages: []*ensurefile.Package{
						{
							Path:       "github.com/my/mod/abc/internal/abc",
							Interfaces: []string{"Iface1"},
						},
					},
				},
			},

			SetupMocks: func(m *Mocks) {
				const expectedMockFile1 = `<abc mock stuff here>

// NEW creates a MockIface1.
func (*MockIface1) NEW(ctrl *gomock.Controller) *MockIface1 {
	return NewMockIface1(ctrl)
}
`

				gomock.InOrder(
					m.CmdRun.EXPECT().Exec(&runcmd.ExecParams{
						PWD:  "/root/path/abc",
						CMD:  "mockgen",
						Args: []string{"github.com/my/mod/abc/internal/abc", "Iface1"},
					}).Return("<abc mock stuff here>\n", nil),

					m.FSWrite.EXPECT().
						MkdirAll("/root/path/abc/internal/mocks/mock_abc", expectedDirPerm).
						Return(nil),

					m.FSWrite.EXPECT().
						WriteFile(
							"/root/path/abc/internal/mocks/mock_abc/mock_abc.go",
							expectedMockFile1,
							expectedFilePerm,
						).
						Return(nil),
				)
			},
		},

		{
			Name:          "when unable to run mockgen",
			ExpectedError: mockgen.ErrMockGenFailed,
			Config: &ensurefile.Config{
				RootPath:   "/root/path",
				ModulePath: "github.com/my/mod",
				Mocks: &ensurefile.MockConfig{
					Packages: []*ensurefile.Package{
						{
							Path:       "github.com/some/pkg/abc",
							Interfaces: []string{"Iface1"},
						},
					},
				},
			},

			SetupMocks: func(m *Mocks) {
				m.CmdRun.EXPECT().Exec(&runcmd.ExecParams{
					PWD:  "/root/path",
					CMD:  "mockgen",
					Args: []string{"github.com/some/pkg/abc", "Iface1"},
				}).Return("", errors.New("mockgen error"))
			},
		},

		{
			Name:          "when unable to create directory",
			ExpectedError: mockgen.ErrUnableToCreateDir,
			Config: &ensurefile.Config{
				RootPath:   "/root/path",
				ModulePath: "github.com/my/mod",
				Mocks: &ensurefile.MockConfig{
					Packages: []*ensurefile.Package{
						{
							Path:       "github.com/some/pkg/abc",
							Interfaces: []string{"Iface1"},
						},
					},
				},
			},

			SetupMocks: func(m *Mocks) {
				m.CmdRun.EXPECT().Exec(&runcmd.ExecParams{
					PWD:  "/root/path",
					CMD:  "mockgen",
					Args: []string{"github.com/some/pkg/abc", "Iface1"},
				}).Return("<abc mock stuff here>\n", nil)

				m.FSWrite.EXPECT().
					MkdirAll("/root/path/internal/mocks/github.com/some/pkg/mock_abc", expectedDirPerm).
					Return(errors.New("couldn't create the directory"))
			},
		},

		{
			Name:          "when unable to create file",
			ExpectedError: mockgen.ErrUnableToCreateDir,
			Config: &ensurefile.Config{
				RootPath:   "/root/path",
				ModulePath: "github.com/my/mod",
				Mocks: &ensurefile.MockConfig{
					Packages: []*ensurefile.Package{
						{
							Path:       "github.com/some/pkg/abc",
							Interfaces: []string{"Iface1"},
						},
					},
				},
			},

			SetupMocks: func(m *Mocks) {
				m.CmdRun.EXPECT().Exec(&runcmd.ExecParams{
					PWD:  "/root/path",
					CMD:  "mockgen",
					Args: []string{"github.com/some/pkg/abc", "Iface1"},
				}).Return("<abc mock stuff here>\n", nil)

				m.FSWrite.EXPECT().
					MkdirAll("/root/path/internal/mocks/github.com/some/pkg/mock_abc", expectedDirPerm).
					Return(nil)

				m.FSWrite.EXPECT().
					WriteFile(
						"/root/path/internal/mocks/github.com/some/pkg/mock_abc/mock_abc.go",
						gomock.Any(),
						expectedFilePerm,
					).
					Return(errors.New("some write failure"))
			},
		},

		{
			Name:          "when missing mocks",
			ExpectedError: mockgen.ErrMissingMockConfig,
			Config: &ensurefile.Config{
				RootPath:   "/root/path",
				ModulePath: "github.com/my/mod",
				Mocks:      nil, // Missing mocks
			},
		},

		{
			Name:          "when missing package mocks",
			ExpectedError: mockgen.ErrMissingPackages,
			Config: &ensurefile.Config{
				RootPath:   "/root/path",
				ModulePath: "github.com/my/mod",
				Mocks: &ensurefile.MockConfig{
					Packages: []*ensurefile.Package{}, // Missing package mocks
				},
			},
		},

		{
			Name:          "when package mock missing path",
			ExpectedError: mockgen.ErrMissingPackagePath,
			Config: &ensurefile.Config{
				RootPath:   "/root/path",
				ModulePath: "github.com/my/mod",
				Mocks: &ensurefile.MockConfig{
					Packages: []*ensurefile.Package{
						{
							Path:       "", // Missing path
							Interfaces: []string{"Iface1"},
						},
					},
				},
			},
		},

		{
			Name:          "when package mock missing interfaces",
			ExpectedError: mockgen.ErrMissingPackageInterfaces,
			Config: &ensurefile.Config{
				RootPath:   "/root/path",
				ModulePath: "github.com/my/mod",
				Mocks: &ensurefile.MockConfig{
					Packages: []*ensurefile.Package{
						{
							Path:       "github.com/my/pkg",
							Interfaces: nil, // Missing interfaces
						},
					},
				},
			},
		},

		{
			Name:          "when package path duplicated",
			ExpectedError: mockgen.ErrDuplicatePackagePath,
			Config: &ensurefile.Config{
				RootPath:   "/root/path",
				ModulePath: "github.com/my/mod",
				Mocks: &ensurefile.MockConfig{
					Packages: []*ensurefile.Package{
						{
							Path:       "github.com/my/pkg",
							Interfaces: []string{"Iface1"},
						},
						{
							Path:       "github.com/my/other/pkg",
							Interfaces: []string{"Iface1"},
						},
						{
							Path:       "github.com/my/pkg", // Duplicate
							Interfaces: []string{"Iface1"},
						},
					},
				},
			},
		},

		{
			Name:          "when internal package is outside module",
			ExpectedError: mockgen.ErrInternalPackageOutsideModule,
			Config: &ensurefile.Config{
				RootPath:   "/root/path",
				ModulePath: "github.com/my/mod",
				Mocks: &ensurefile.MockConfig{
					Packages: []*ensurefile.Package{
						{
							Path:       "github.com/some/pkg/internal/abc",
							Interfaces: []string{"Iface1"},
						},
					},
				},
			},
		},
	}

	ensure.RunTableByIndex(table, func(ensure ensurepkg.Ensure, i int) {
		entry := table[i]

		err := entry.Subject.GenerateMocks(entry.Config)
		ensure(err).IsError(err)
	})
}
