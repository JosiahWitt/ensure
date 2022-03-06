package mockwrite_test

import (
	"errors"
	"io/fs"
	"io/ioutil"
	"log"
	"testing"

	"github.com/JosiahWitt/ensure"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ensurefile"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/mockgen"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/mocks/mock_fswrite"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/mockwrite"
	"github.com/JosiahWitt/ensure/ensurepkg"
)

func TestWriteMocks(t *testing.T) {
	ensure := ensure.New(t)

	type Mocks struct {
		FS *mock_fswrite.MockWritable
	}

	table := []struct {
		Name string

		Config         *ensurefile.Config
		GeneratedMocks []*mockgen.PackageMock

		ExpectedError error

		Mocks      *Mocks
		SetupMocks func(*Mocks)
		Subject    *mockwrite.MockWriter
	}{
		{
			Name: "when provided non-internal mocks",

			Config: &ensurefile.Config{
				RootPath:   "/my/root",
				ModulePath: "github.com/my/mod",

				Mocks: &ensurefile.MockConfig{
					PrimaryDestination:  "primary-mocks",
					InternalDestination: "internal-mocks",
					TidyAfterGenerate:   true,

					Packages: []*ensurefile.MockPackage{
						{
							Path:       "github.com/my/mod/pkgs/pkg1",
							Interfaces: []string{"Iface1", "Iface2"},
						},
						{
							Path:       "github.com/my/mod/pkgs/pkg2",
							Interfaces: []string{"Iface3", "Iface4"},
						},
					},
				},
			},

			GeneratedMocks: []*mockgen.PackageMock{
				{
					Package: &ifacereader.Package{
						Name: "pkg1",
						Path: "github.com/my/mod/pkgs/pkg1",
						// Other fields are unused by this package
					},

					FileContents: "pkg1 file!",
				},
				{
					Package: &ifacereader.Package{
						Name: "pkg2",
						Path: "github.com/my/mod/pkgs/pkg2",
						// Other fields are unused by this package
					},

					FileContents: "pkg2 file!",
				},
			},

			SetupMocks: func(m *Mocks) {
				m.FS.EXPECT().MkdirAll("/my/root/primary-mocks/github.com/my/mod/pkgs/mock_pkg1", fs.FileMode(0775))
				m.FS.EXPECT().WriteFile("/my/root/primary-mocks/github.com/my/mod/pkgs/mock_pkg1/mock_pkg1.go", "pkg1 file!", fs.FileMode(0664))

				m.FS.EXPECT().MkdirAll("/my/root/primary-mocks/github.com/my/mod/pkgs/mock_pkg2", fs.FileMode(0775))
				m.FS.EXPECT().WriteFile("/my/root/primary-mocks/github.com/my/mod/pkgs/mock_pkg2/mock_pkg2.go", "pkg2 file!", fs.FileMode(0664))
			},
		},
		{
			Name: "when provided internal mocks",

			Config: &ensurefile.Config{
				RootPath:   "/my/root",
				ModulePath: "github.com/my/mod",

				Mocks: &ensurefile.MockConfig{
					PrimaryDestination:  "primary-mocks",
					InternalDestination: "internal-mocks",
					TidyAfterGenerate:   true,

					Packages: []*ensurefile.MockPackage{
						{
							Path:       "github.com/my/mod/lib/internal/pkgs/pkg1",
							Interfaces: []string{"Iface1", "Iface2"},
						},
						{
							// Doubly nested internal package
							Path:       "github.com/my/mod/lib/internal/pkgs/internal/thing/pkg2",
							Interfaces: []string{"Iface3", "Iface4"},
						},
					},
				},
			},

			GeneratedMocks: []*mockgen.PackageMock{
				{
					Package: &ifacereader.Package{
						Name: "pkg1",
						Path: "github.com/my/mod/lib/internal/pkgs/pkg1",
						// Other fields are unused by this package
					},

					FileContents: "pkg1 file!",
				},
				{
					Package: &ifacereader.Package{
						// Doubly nested internal package
						Name: "pkg2",
						Path: "github.com/my/mod/lib/internal/pkgs/internal/thing/pkg2",
						// Other fields are unused by this package
					},

					FileContents: "pkg2 file!",
				},
			},

			SetupMocks: func(m *Mocks) {
				m.FS.EXPECT().MkdirAll("/my/root/lib/internal/internal-mocks/pkgs/mock_pkg1", fs.FileMode(0775))
				m.FS.EXPECT().WriteFile("/my/root/lib/internal/internal-mocks/pkgs/mock_pkg1/mock_pkg1.go", "pkg1 file!", fs.FileMode(0664))

				m.FS.EXPECT().MkdirAll("/my/root/lib/internal/pkgs/internal/internal-mocks/thing/mock_pkg2", fs.FileMode(0775))
				m.FS.EXPECT().WriteFile("/my/root/lib/internal/pkgs/internal/internal-mocks/thing/mock_pkg2/mock_pkg2.go", "pkg2 file!", fs.FileMode(0664))
			},
		},
		{
			Name: "when provided mocks external to package",

			Config: &ensurefile.Config{
				RootPath:   "/my/root",
				ModulePath: "github.com/my/mod",

				Mocks: &ensurefile.MockConfig{
					PrimaryDestination:  "primary-mocks",
					InternalDestination: "internal-mocks",
					TidyAfterGenerate:   true,

					Packages: []*ensurefile.MockPackage{
						{
							Path:       "github.com/not/my/mod/pkgs/pkg1",
							Interfaces: []string{"Iface1", "Iface2"},
						},
						{
							Path:       "github.com/not/my/mod/pkgs/pkg2",
							Interfaces: []string{"Iface3", "Iface4"},
						},
					},
				},
			},

			GeneratedMocks: []*mockgen.PackageMock{
				{
					Package: &ifacereader.Package{
						Name: "pkg1",
						Path: "github.com/not/my/mod/pkgs/pkg1",
						// Other fields are unused by this package
					},

					FileContents: "pkg1 file!",
				},
				{
					Package: &ifacereader.Package{
						Name: "pkg2",
						Path: "github.com/not/my/mod/pkgs/pkg2",
						// Other fields are unused by this package
					},

					FileContents: "pkg2 file!",
				},
			},

			SetupMocks: func(m *Mocks) {
				m.FS.EXPECT().MkdirAll("/my/root/primary-mocks/github.com/not/my/mod/pkgs/mock_pkg1", fs.FileMode(0775))
				m.FS.EXPECT().WriteFile("/my/root/primary-mocks/github.com/not/my/mod/pkgs/mock_pkg1/mock_pkg1.go", "pkg1 file!", fs.FileMode(0664))

				m.FS.EXPECT().MkdirAll("/my/root/primary-mocks/github.com/not/my/mod/pkgs/mock_pkg2", fs.FileMode(0775))
				m.FS.EXPECT().WriteFile("/my/root/primary-mocks/github.com/not/my/mod/pkgs/mock_pkg2/mock_pkg2.go", "pkg2 file!", fs.FileMode(0664))
			},
		},
		{
			Name: "when provided internal, non-internal and external mocks",

			Config: &ensurefile.Config{
				RootPath:   "/my/root",
				ModulePath: "github.com/my/mod",

				Mocks: &ensurefile.MockConfig{
					PrimaryDestination:  "primary-mocks",
					InternalDestination: "internal-mocks",
					TidyAfterGenerate:   true,

					Packages: []*ensurefile.MockPackage{
						{
							Path:       "github.com/my/mod/lib/internal/pkgs/pkg1",
							Interfaces: []string{"Iface1", "Iface2"},
						},
						{
							Path:       "github.com/my/mod/pkgs/pkg2",
							Interfaces: []string{"Iface3", "Iface4"},
						},
						{
							Path:       "github.com/not/my/mod/pkgs/pkg1",
							Interfaces: []string{"Iface5", "Iface5"},
						},
					},
				},
			},

			GeneratedMocks: []*mockgen.PackageMock{
				{
					Package: &ifacereader.Package{
						Name: "pkg1",
						Path: "github.com/my/mod/lib/internal/pkgs/pkg1",
						// Other fields are unused by this package
					},

					FileContents: "pkg1 file!",
				},
				{
					Package: &ifacereader.Package{
						Name: "pkg2",
						Path: "github.com/my/mod/pkgs/pkg2",
						// Other fields are unused by this package
					},

					FileContents: "pkg2 file!",
				},
				{
					Package: &ifacereader.Package{
						Name: "pkg1",
						Path: "github.com/not/my/mod/pkgs/pkg1",
						// Other fields are unused by this package
					},

					FileContents: "pkg2 file!",
				},
			},

			SetupMocks: func(m *Mocks) {
				m.FS.EXPECT().MkdirAll("/my/root/lib/internal/internal-mocks/pkgs/mock_pkg1", fs.FileMode(0775))
				m.FS.EXPECT().WriteFile("/my/root/lib/internal/internal-mocks/pkgs/mock_pkg1/mock_pkg1.go", "pkg1 file!", fs.FileMode(0664))

				m.FS.EXPECT().MkdirAll("/my/root/primary-mocks/github.com/my/mod/pkgs/mock_pkg2", fs.FileMode(0775))
				m.FS.EXPECT().WriteFile("/my/root/primary-mocks/github.com/my/mod/pkgs/mock_pkg2/mock_pkg2.go", "pkg2 file!", fs.FileMode(0664))

				m.FS.EXPECT().MkdirAll("/my/root/primary-mocks/github.com/not/my/mod/pkgs/mock_pkg1", fs.FileMode(0775))
				m.FS.EXPECT().WriteFile("/my/root/primary-mocks/github.com/not/my/mod/pkgs/mock_pkg1/mock_pkg1.go", "pkg2 file!", fs.FileMode(0664))
			},
		},
		{
			Name: "when provided mocks for a v2 package",

			Config: &ensurefile.Config{
				RootPath:   "/my/root",
				ModulePath: "github.com/my/mod",

				Mocks: &ensurefile.MockConfig{
					PrimaryDestination:  "primary-mocks",
					InternalDestination: "internal-mocks",
					TidyAfterGenerate:   true,

					// Mocks for these two paths will collide in the same file.
					// Hopefully this doesn't happen much in the real world.
					Packages: []*ensurefile.MockPackage{
						{
							Path:       "github.com/my/pkg1/v2",
							Interfaces: []string{"Iface1", "Iface2"},
						},
						{
							Path:       "github.com/my/pkg1/v2/pkg1",
							Interfaces: []string{"Iface3", "Iface4"},
						},
					},
				},
			},

			GeneratedMocks: []*mockgen.PackageMock{
				{
					Package: &ifacereader.Package{
						Name: "pkg1",
						Path: "github.com/my/pkg1/v2",
						// Other fields are unused by this package
					},

					FileContents: "pkg1 file!",
				},
				{
					Package: &ifacereader.Package{
						Name: "pkg1",
						Path: "github.com/my/pkg1/v2/pkg1",
						// Other fields are unused by this package
					},

					FileContents: "nested pkg1 file!",
				},
			},

			SetupMocks: func(m *Mocks) {
				m.FS.EXPECT().MkdirAll("/my/root/primary-mocks/github.com/my/pkg1/v2/mock_pkg1", fs.FileMode(0775))
				m.FS.EXPECT().WriteFile("/my/root/primary-mocks/github.com/my/pkg1/v2/mock_pkg1/mock_pkg1.go", "pkg1 file!", fs.FileMode(0664))

				m.FS.EXPECT().MkdirAll("/my/root/primary-mocks/github.com/my/pkg1/v2/mock_pkg1", fs.FileMode(0775))
				m.FS.EXPECT().WriteFile("/my/root/primary-mocks/github.com/my/pkg1/v2/mock_pkg1/mock_pkg1.go", "nested pkg1 file!", fs.FileMode(0664))
			},
		},
		{
			Name: "returns ErrInternalPackageOutsideModule when provided internal mock that does not belong to the current module",

			Config: &ensurefile.Config{
				RootPath:   "/my/root",
				ModulePath: "github.com/my/mod",

				Mocks: &ensurefile.MockConfig{
					PrimaryDestination:  "primary-mocks",
					InternalDestination: "internal-mocks",
					TidyAfterGenerate:   true,

					Packages: []*ensurefile.MockPackage{
						{
							Path:       "github.com/not/my/mod/internal/pkgs/pkg1",
							Interfaces: []string{"Iface1", "Iface2"},
						},
						{
							Path:       "github.com/my/mod/pkgs/pkg2",
							Interfaces: []string{"Iface3", "Iface4"},
						},
					},
				},
			},

			GeneratedMocks: []*mockgen.PackageMock{
				{
					Package: &ifacereader.Package{
						Name: "pkg1",
						Path: "github.com/not/my/mod/internal/pkgs/pkg1",
						// Other fields are unused by this package
					},

					FileContents: "pkg1 file!",
				},
				{
					Package: &ifacereader.Package{
						Name: "pkg2",
						Path: "github.com/my/mod/pkgs/pkg2",
						// Other fields are unused by this package
					},

					FileContents: "pkg2 file!",
				},
			},

			ExpectedError: mockwrite.ErrInternalPackageOutsideModule,

			SetupMocks: func(m *Mocks) {
				// It continues on writing other files even if one fails
				m.FS.EXPECT().MkdirAll("/my/root/primary-mocks/github.com/my/mod/pkgs/mock_pkg2", fs.FileMode(0775))
				m.FS.EXPECT().WriteFile("/my/root/primary-mocks/github.com/my/mod/pkgs/mock_pkg2/mock_pkg2.go", "pkg2 file!", fs.FileMode(0664))
			},
		},
		{
			Name: "returns ErrUnableToCreateDir when unable to create the directories",

			Config: &ensurefile.Config{
				RootPath:   "/my/root",
				ModulePath: "github.com/my/mod",

				Mocks: &ensurefile.MockConfig{
					PrimaryDestination:  "primary-mocks",
					InternalDestination: "internal-mocks",
					TidyAfterGenerate:   true,

					Packages: []*ensurefile.MockPackage{
						{
							Path:       "github.com/my/mod/pkgs/pkg1",
							Interfaces: []string{"Iface1", "Iface2"},
						},
						{
							Path:       "github.com/my/mod/pkgs/pkg2",
							Interfaces: []string{"Iface3", "Iface4"},
						},
					},
				},
			},

			GeneratedMocks: []*mockgen.PackageMock{
				{
					Package: &ifacereader.Package{
						Name: "pkg1",
						Path: "github.com/my/mod/pkgs/pkg1",
						// Other fields are unused by this package
					},

					FileContents: "pkg1 file!",
				},
				{
					Package: &ifacereader.Package{
						Name: "pkg2",
						Path: "github.com/my/mod/pkgs/pkg2",
						// Other fields are unused by this package
					},

					FileContents: "pkg2 file!",
				},
			},

			ExpectedError: mockwrite.ErrUnableToCreateDir,

			SetupMocks: func(m *Mocks) {
				m.FS.EXPECT().
					MkdirAll("/my/root/primary-mocks/github.com/my/mod/pkgs/mock_pkg1", fs.FileMode(0775)).
					Return(errors.New("a file with that path exists"))

				// It continues on writing other files even if one fails
				m.FS.EXPECT().MkdirAll("/my/root/primary-mocks/github.com/my/mod/pkgs/mock_pkg2", fs.FileMode(0775))
				m.FS.EXPECT().WriteFile("/my/root/primary-mocks/github.com/my/mod/pkgs/mock_pkg2/mock_pkg2.go", "pkg2 file!", fs.FileMode(0664))
			},
		},
		{
			Name: "returns ErrUnableToCreateFile when unable to create the directories",

			Config: &ensurefile.Config{
				RootPath:   "/my/root",
				ModulePath: "github.com/my/mod",

				Mocks: &ensurefile.MockConfig{
					PrimaryDestination:  "primary-mocks",
					InternalDestination: "internal-mocks",
					TidyAfterGenerate:   true,

					Packages: []*ensurefile.MockPackage{
						{
							Path:       "github.com/my/mod/pkgs/pkg1",
							Interfaces: []string{"Iface1", "Iface2"},
						},
						{
							Path:       "github.com/my/mod/pkgs/pkg2",
							Interfaces: []string{"Iface3", "Iface4"},
						},
					},
				},
			},

			GeneratedMocks: []*mockgen.PackageMock{
				{
					Package: &ifacereader.Package{
						Name: "pkg1",
						Path: "github.com/my/mod/pkgs/pkg1",
						// Other fields are unused by this package
					},

					FileContents: "pkg1 file!",
				},
				{
					Package: &ifacereader.Package{
						Name: "pkg2",
						Path: "github.com/my/mod/pkgs/pkg2",
						// Other fields are unused by this package
					},

					FileContents: "pkg2 file!",
				},
			},

			ExpectedError: mockwrite.ErrUnableToCreateFile,

			SetupMocks: func(m *Mocks) {
				m.FS.EXPECT().MkdirAll("/my/root/primary-mocks/github.com/my/mod/pkgs/mock_pkg1", fs.FileMode(0775))
				m.FS.EXPECT().
					WriteFile("/my/root/primary-mocks/github.com/my/mod/pkgs/mock_pkg1/mock_pkg1.go", "pkg1 file!", fs.FileMode(0664)).
					Return(errors.New("permission denied"))

				// It continues on writing other files even if one fails
				m.FS.EXPECT().MkdirAll("/my/root/primary-mocks/github.com/my/mod/pkgs/mock_pkg2", fs.FileMode(0775))
				m.FS.EXPECT().WriteFile("/my/root/primary-mocks/github.com/my/mod/pkgs/mock_pkg2/mock_pkg2.go", "pkg2 file!", fs.FileMode(0664))
			},
		},
	}

	ensure.RunTableByIndex(table, func(ensure ensurepkg.Ensure, i int) {
		entry := table[i]
		entry.Subject.Logger = log.New(ioutil.Discard, "", 0)

		err := entry.Subject.WriteMocks(entry.Config, entry.GeneratedMocks)
		ensure(err).IsError(entry.ExpectedError)
	})
}
