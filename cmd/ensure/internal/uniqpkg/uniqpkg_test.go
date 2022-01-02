package uniqpkg_test

import (
	"go/types"
	"testing"

	"github.com/JosiahWitt/ensure"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/uniqpkg"
	"github.com/JosiahWitt/ensure/ensurepkg"
	"golang.org/x/tools/go/packages"
)

func TestGeneratePackageName(t *testing.T) {
	ensure := ensure.New(t)

	table := []struct {
		Name string

		Runner func(p *uniqpkg.UniquePackagePaths)

		ExpectedResult []*uniqpkg.Package

		Subject *uniqpkg.UniquePackagePaths
	}{
		{
			Name: "inside one package scope with single imported package referenced once",

			Runner: func(p *uniqpkg.UniquePackagePaths) {
				scopePkg1 := &packages.Package{Name: "scopepkg1", PkgPath: "pkgs/scopepkg1"}
				p.GeneratePackageName(scopePkg1, types.NewPackage("pkgs/importedpkg1", "importedpkg1"))
			},

			ExpectedResult: []*uniqpkg.Package{
				{
					Name: "scopepkg1",
					Path: "pkgs/scopepkg1",
					Imports: []*uniqpkg.ImportDetails{
						{Name: "importedpkg1", Path: "pkgs/importedpkg1"},
					},
				},
			},
		},
		{
			Name: "inside one package scope with single imported package referenced multiple times",

			Runner: func(p *uniqpkg.UniquePackagePaths) {
				scopePkg1 := &packages.Package{Name: "scopepkg1", PkgPath: "pkgs/scopepkg1"}
				importedPkg1 := types.NewPackage("pkgs/importedpkg1", "importedpkg1")

				p.GeneratePackageName(scopePkg1, importedPkg1)
				p.GeneratePackageName(scopePkg1, importedPkg1)
				p.GeneratePackageName(scopePkg1, importedPkg1)
			},

			ExpectedResult: []*uniqpkg.Package{
				{
					Name: "scopepkg1",
					Path: "pkgs/scopepkg1",
					Imports: []*uniqpkg.ImportDetails{
						{Name: "importedpkg1", Path: "pkgs/importedpkg1"},
					},
				},
			},
		},
		{
			Name: "inside one package scope with multiple imported packages referenced multiple times",

			Runner: func(p *uniqpkg.UniquePackagePaths) {
				scopePkg1 := &packages.Package{Name: "scopepkg1", PkgPath: "pkgs/scopepkg1"}
				importedPkg1 := types.NewPackage("pkgs/importedpkg1", "importedpkg1")
				importedPkg2 := types.NewPackage("pkgs/importedpkg2", "importedpkg2")

				p.GeneratePackageName(scopePkg1, importedPkg1)
				p.GeneratePackageName(scopePkg1, importedPkg2)
				p.GeneratePackageName(scopePkg1, importedPkg1)
				p.GeneratePackageName(scopePkg1, importedPkg2)
			},

			ExpectedResult: []*uniqpkg.Package{
				{
					Name: "scopepkg1",
					Path: "pkgs/scopepkg1",
					Imports: []*uniqpkg.ImportDetails{
						{Name: "importedpkg1", Path: "pkgs/importedpkg1"},
						{Name: "importedpkg2", Path: "pkgs/importedpkg2"},
					},
				},
			},
		},
		{
			Name: "inside multiple package scopes with multiple imported packages referenced multiple times",

			Runner: func(p *uniqpkg.UniquePackagePaths) {
				scopePkg1 := &packages.Package{Name: "scopepkg1", PkgPath: "pkgs/scopepkg1"}
				scopePkg2 := &packages.Package{Name: "scopepkg2", PkgPath: "pkgs/scopepkg2"}
				importedPkg1 := types.NewPackage("pkgs/importedpkg1", "importedpkg1")
				importedPkg2 := types.NewPackage("pkgs/importedpkg2", "importedpkg2")
				importedPkg3 := types.NewPackage("pkgs/importedpkg3", "importedpkg3")

				p.GeneratePackageName(scopePkg1, importedPkg1)
				p.GeneratePackageName(scopePkg1, importedPkg2)
				p.GeneratePackageName(scopePkg2, importedPkg1)
				p.GeneratePackageName(scopePkg2, importedPkg2)

				p.GeneratePackageName(scopePkg1, importedPkg1)
				p.GeneratePackageName(scopePkg1, importedPkg2)
				p.GeneratePackageName(scopePkg2, importedPkg1)
				p.GeneratePackageName(scopePkg2, importedPkg2)

				// Only pkg3 is imported in scope 2
				p.GeneratePackageName(scopePkg2, importedPkg3)
				p.GeneratePackageName(scopePkg2, importedPkg3)
			},

			ExpectedResult: []*uniqpkg.Package{
				{
					Name: "scopepkg1",
					Path: "pkgs/scopepkg1",
					Imports: []*uniqpkg.ImportDetails{
						{Name: "importedpkg1", Path: "pkgs/importedpkg1"},
						{Name: "importedpkg2", Path: "pkgs/importedpkg2"},
					},
				},
				{
					Name: "scopepkg2",
					Path: "pkgs/scopepkg2",
					Imports: []*uniqpkg.ImportDetails{
						{Name: "importedpkg1", Path: "pkgs/importedpkg1"},
						{Name: "importedpkg2", Path: "pkgs/importedpkg2"},
						{Name: "importedpkg3", Path: "pkgs/importedpkg3"},
					},
				},
			},
		},
		{
			Name: "two package scopes import a common package name with different package paths",

			Runner: func(p *uniqpkg.UniquePackagePaths) {
				scopePkg1 := &packages.Package{Name: "scopepkg1", PkgPath: "pkgs/scopepkg1"}
				scopePkg2 := &packages.Package{Name: "scopepkg2", PkgPath: "pkgs/scopepkg2"}

				p.GeneratePackageName(scopePkg1, types.NewPackage("pkgs/1/common", "common"))
				p.GeneratePackageName(scopePkg2, types.NewPackage("pkgs/2/common", "common"))
			},

			ExpectedResult: []*uniqpkg.Package{
				{
					Name: "scopepkg1",
					Path: "pkgs/scopepkg1",
					Imports: []*uniqpkg.ImportDetails{
						{Name: "common", Path: "pkgs/1/common"},
					},
				},
				{
					Name: "scopepkg2",
					Path: "pkgs/scopepkg2",
					Imports: []*uniqpkg.ImportDetails{
						{Name: "common", Path: "pkgs/2/common"},
					},
				},
			},
		},
		{
			Name: "one package scope imports two packages with the same name but different paths",

			Runner: func(p *uniqpkg.UniquePackagePaths) {
				scopePkg1 := &packages.Package{Name: "scopepkg1", PkgPath: "pkgs/scopepkg1"}

				p.GeneratePackageName(scopePkg1, types.NewPackage("pkgs/1/common", "common"))
				p.GeneratePackageName(scopePkg1, types.NewPackage("pkgs/2/common", "common"))
			},

			ExpectedResult: []*uniqpkg.Package{
				{
					Name: "scopepkg1",
					Path: "pkgs/scopepkg1",
					Imports: []*uniqpkg.ImportDetails{
						{Name: "common", Path: "pkgs/1/common"},
						{Name: "common2", Path: "pkgs/2/common", IsAlias: true},
					},
				},
			},
		},
		{
			Name: "one package scope imports multiple packages with the same name but different paths",

			Runner: func(p *uniqpkg.UniquePackagePaths) {
				scopePkg1 := &packages.Package{Name: "scopepkg1", PkgPath: "pkgs/scopepkg1"}

				p.GeneratePackageName(scopePkg1, types.NewPackage("pkgs/1/common", "common"))
				p.GeneratePackageName(scopePkg1, types.NewPackage("pkgs/2/common", "common"))
				p.GeneratePackageName(scopePkg1, types.NewPackage("pkgs/3/common", "common"))

				// To show that the generated names are repeatable
				p.GeneratePackageName(scopePkg1, types.NewPackage("pkgs/2/common", "common"))
				p.GeneratePackageName(scopePkg1, types.NewPackage("pkgs/3/common", "common"))
				p.GeneratePackageName(scopePkg1, types.NewPackage("pkgs/1/common", "common"))
			},

			ExpectedResult: []*uniqpkg.Package{
				{
					Name: "scopepkg1",
					Path: "pkgs/scopepkg1",
					Imports: []*uniqpkg.ImportDetails{
						{Name: "common", Path: "pkgs/1/common"},
						{Name: "common2", Path: "pkgs/2/common", IsAlias: true},
						{Name: "common3", Path: "pkgs/3/common", IsAlias: true},
					},
				},
			},
		},
	}

	ensure.RunTableByIndex(table, func(ensure ensurepkg.Ensure, i int) {
		entry := table[i]

		entry.Runner(entry.Subject)
		ensure(entry.Subject.Export()).Equals(entry.ExpectedResult)
	})
}
