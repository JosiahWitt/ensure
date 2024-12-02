package mockwrite_test

import (
	"errors"
	"io"
	"log"
	"testing"

	"github.com/JosiahWitt/ensure"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ensurefile"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/mocks/mock_fswrite"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/mockwrite"
	"github.com/JosiahWitt/ensure/ensuring"
)

func TestTidyMocks(t *testing.T) {
	ensure := ensure.New(t)

	type Mocks struct {
		FSWrite *mock_fswrite.MockWritable
	}

	table := []struct {
		Name string

		Config   *ensurefile.Config
		Packages []*ifacereader.Package

		ExpectedError error

		Mocks      *Mocks
		SetupMocks func(*Mocks)
		Subject    *mockwrite.MockWriter
	}{
		{
			Name: "with files to delete",

			Config: &ensurefile.Config{
				RootPath:   "/root/path",
				ModulePath: "github.com/my/mod",
				Mocks: &ensurefile.MockConfig{
					PrimaryDestination:  "primary_mocks",
					InternalDestination: "internal_mocks",
					Packages: []*ensurefile.MockPackage{
						{
							Path:       "github.com/some/pkg/abc",
							Interfaces: []string{"Iface1"},
						},
						{
							Path:       "github.com/some/pkg/qwerty/v2",
							Interfaces: []string{"Iface2"},
						},
						{
							Path:       "github.com/my/mod/layer1/layer2/internal/layer3/layer4/internal/layer5/layer6/xyz",
							Interfaces: []string{"Iface3"},
						},
					},
				},
			},

			Packages: []*ifacereader.Package{
				{
					Name: "abc",
					Path: "github.com/some/pkg/abc",
					// Other fields are unused by this package
				},
				{
					Name: "qwerty",
					Path: "github.com/some/pkg/qwerty/v2",
					// Other fields are unused by this package
				},
				{
					Name: "xyz",
					Path: "github.com/my/mod/layer1/layer2/internal/layer3/layer4/internal/layer5/layer6/xyz",
					// Other fields are unused by this package
				},
			},

			SetupMocks: func(m *Mocks) {
				const primaryMocksDir = "/root/path/primary_mocks"
				m.FSWrite.EXPECT().ListRecursive(primaryMocksDir).
					Return([]string{
						primaryMocksDir + "/github.com",
						primaryMocksDir + "/github.com/some",
						primaryMocksDir + "/github.com/some/pkg",
						primaryMocksDir + "/github.com/some/pkg/mock_abc",
						primaryMocksDir + "/github.com/some/pkg/mock_abc/mock_abc.go",
						primaryMocksDir + "/github.com/some/pkg/qwerty/v2/mock_qwerty",
						primaryMocksDir + "/github.com/some/pkg/qwerty/v2/mock_qwerty/mock_qwerty.go",

						// Extra files

						primaryMocksDir + "/github.com/some/pkg/mock_qwerty/extra_file.go",
						primaryMocksDir + "/github.com/some/pkg/mock_qwerty/extra_dir",
						primaryMocksDir + "/github.com/some/pkg/mock_qwerty/extra_dir/with_file.go",

						primaryMocksDir + "/github.com/d3l3t3.m3",
						primaryMocksDir + "/somefile.txt",
						primaryMocksDir + "/some",
						primaryMocksDir + "/some/nesting",
						primaryMocksDir + "/some/nesting/file1.txt",
						primaryMocksDir + "/some/nesting/file2.txt",
						primaryMocksDir + "/some/hello.txt",
					}, nil)

				const internalMocksDir = "/root/path/layer1/layer2/internal/layer3/layer4/internal/internal_mocks"
				m.FSWrite.EXPECT().ListRecursive(internalMocksDir).
					Return([]string{
						internalMocksDir + "/layer5",
						internalMocksDir + "/layer5/layer6",
						internalMocksDir + "/layer5/layer6/mock_xyz",
						internalMocksDir + "/layer5/layer6/mock_xyz/mock_xyz.go",

						// Extra files
						internalMocksDir + "/layer5/layer6/mock_xyz/extra123.go",
						internalMocksDir + "/layer5/layer6/mock_xyz/nested",
						internalMocksDir + "/layer5/layer6/mock_xyz/nested/more.go",
						internalMocksDir + "/layer5/hello",
						internalMocksDir + "/layer5/hello/there.hi",
						internalMocksDir + "/garbage.txt",
					}, nil)

				expectedPathsToDelete := []string{
					// Primary mocks
					primaryMocksDir + "/github.com/some/pkg/mock_qwerty/extra_file.go",
					primaryMocksDir + "/github.com/some/pkg/mock_qwerty/extra_dir",
					primaryMocksDir + "/github.com/some/pkg/mock_qwerty/extra_dir/with_file.go",

					primaryMocksDir + "/github.com/d3l3t3.m3",
					primaryMocksDir + "/somefile.txt",
					primaryMocksDir + "/some",
					primaryMocksDir + "/some/nesting",
					primaryMocksDir + "/some/nesting/file1.txt",
					primaryMocksDir + "/some/nesting/file2.txt",
					primaryMocksDir + "/some/hello.txt",

					// Internal mocks
					internalMocksDir + "/layer5/layer6/mock_xyz/extra123.go",
					internalMocksDir + "/layer5/layer6/mock_xyz/nested",
					internalMocksDir + "/layer5/layer6/mock_xyz/nested/more.go",
					internalMocksDir + "/layer5/hello",
					internalMocksDir + "/layer5/hello/there.hi",
					internalMocksDir + "/garbage.txt",
				}

				for _, expectedPathToDelete := range expectedPathsToDelete {
					m.FSWrite.EXPECT().RemoveAll(expectedPathToDelete).Return(nil)
				}
			},
		},

		{
			Name: "when already tidy",
			Config: &ensurefile.Config{
				RootPath:   "/root/path",
				ModulePath: "github.com/my/mod",
				Mocks: &ensurefile.MockConfig{
					PrimaryDestination:  "primary_mocks",
					InternalDestination: "internal_mocks",
					Packages: []*ensurefile.MockPackage{
						{
							Path:       "github.com/some/pkg/abc",
							Interfaces: []string{"Iface1"},
						},
						{
							Path:       "github.com/some/pkg/qwerty/v2",
							Interfaces: []string{"Iface2"},
						},
						{
							Path:       "github.com/my/mod/layer1/layer2/internal/layer3/layer4/internal/layer5/layer6/xyz",
							Interfaces: []string{"Iface3"},
						},
					},
				},
			},

			Packages: []*ifacereader.Package{
				{
					Name: "abc",
					Path: "github.com/some/pkg/abc",
					// Other fields are unused by this package
				},
				{
					Name: "qwerty",
					Path: "github.com/some/pkg/qwerty/v2",
					// Other fields are unused by this package
				},
				{
					Name: "xyz",
					Path: "github.com/my/mod/layer1/layer2/internal/layer3/layer4/internal/layer5/layer6/xyz",
					// Other fields are unused by this package
				},
			},

			SetupMocks: func(m *Mocks) {
				const primaryMocksDir = "/root/path/primary_mocks"
				m.FSWrite.EXPECT().ListRecursive(primaryMocksDir).
					Return([]string{
						primaryMocksDir + "/github.com",
						primaryMocksDir + "/github.com/some",
						primaryMocksDir + "/github.com/some/pkg",
						primaryMocksDir + "/github.com/some/pkg/mock_abc",
						primaryMocksDir + "/github.com/some/pkg/mock_abc/mock_abc.go",
						primaryMocksDir + "/github.com/some/pkg/qwerty/v2/mock_qwerty",
						primaryMocksDir + "/github.com/some/pkg/qwerty/v2/mock_qwerty/mock_qwerty.go",
					}, nil)

				const internalMocksDir = "/root/path/layer1/layer2/internal/layer3/layer4/internal/internal_mocks"
				m.FSWrite.EXPECT().ListRecursive(internalMocksDir).
					Return([]string{
						internalMocksDir + "/layer5",
						internalMocksDir + "/layer5/layer6",
						internalMocksDir + "/layer5/layer6/mock_xyz",
						internalMocksDir + "/layer5/layer6/mock_xyz/mock_xyz.go",
					}, nil)
			},
		},

		{
			Name: "with invalid config: internal package outside module",

			Config: &ensurefile.Config{
				RootPath:   "/root/path",
				ModulePath: "github.com/my/mod",
				Mocks: &ensurefile.MockConfig{
					Packages: []*ensurefile.MockPackage{
						{
							Path:       "github.com/not/my/mod/internal/xyz",
							Interfaces: []string{"Iface1"},
						},
					},
				},
			},

			Packages: []*ifacereader.Package{
				{
					Name: "xyz",
					Path: "github.com/not/my/mod/internal/xyz",
				},
			},

			ExpectedError: mockwrite.ErrInternalPackageOutsideModule,
		},

		{
			Name: "when unable to list files recursively",

			Config: &ensurefile.Config{
				RootPath:   "/root/path",
				ModulePath: "github.com/my/mod",
				Mocks: &ensurefile.MockConfig{
					PrimaryDestination:  "primary_mocks",
					InternalDestination: "internal_mocks",
					Packages: []*ensurefile.MockPackage{
						{
							Path:       "github.com/some/pkg/abc",
							Interfaces: []string{"Iface1"},
						},
					},
				},
			},

			Packages: []*ifacereader.Package{
				{
					Name: "abc",
					Path: "github.com/some/pkg/abc",
				},
			},

			ExpectedError: mockwrite.ErrTidyUnableToList,

			SetupMocks: func(m *Mocks) {
				m.FSWrite.EXPECT().ListRecursive("/root/path/primary_mocks").
					Return(nil, errors.New("you can't do that"))
			},
		},

		{
			Name: "when unable to delete files",

			Config: &ensurefile.Config{
				RootPath:   "/root/path",
				ModulePath: "github.com/my/mod",
				Mocks: &ensurefile.MockConfig{
					PrimaryDestination:  "primary_mocks",
					InternalDestination: "internal_mocks",
					Packages: []*ensurefile.MockPackage{
						{
							Path:       "github.com/abc",
							Interfaces: []string{"Iface1"},
						},
					},
				},
			},

			Packages: []*ifacereader.Package{
				{
					Name: "abc",
					Path: "github.com/abc",
				},
			},

			ExpectedError: mockwrite.ErrTidyUnableToCleanup,

			SetupMocks: func(m *Mocks) {
				const primaryMocksDir = "/root/path/primary_mocks"
				m.FSWrite.EXPECT().ListRecursive(primaryMocksDir).
					Return([]string{
						primaryMocksDir + "/github.com",
						primaryMocksDir + "/github.com/mock_abc",
						primaryMocksDir + "/github.com/mock_abc/mock_abc.go",
						primaryMocksDir + "/github.com/extra1.go",
						primaryMocksDir + "/github.com/extra2.go",
					}, nil)

				m.FSWrite.EXPECT().RemoveAll(primaryMocksDir + "/github.com/extra1.go").Return(nil)
				m.FSWrite.EXPECT().RemoveAll(primaryMocksDir + "/github.com/extra2.go").Return(errors.New("oops"))
			},
		},
	}

	ensure.RunTableByIndex(table, func(ensure ensuring.E, i int) {
		entry := table[i]
		entry.Subject.Logger = log.New(io.Discard, "", 0)

		err := entry.Subject.TidyMocks(entry.Config, entry.Packages)
		ensure(err).IsError(entry.ExpectedError)
	})
}
