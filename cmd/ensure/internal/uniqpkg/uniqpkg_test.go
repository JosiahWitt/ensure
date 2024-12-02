package uniqpkg_test

import (
	"go/types"
	"testing"

	"github.com/JosiahWitt/ensure"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/uniqpkg"
	"github.com/JosiahWitt/ensure/ensuring"
	"golang.org/x/tools/go/packages"
)

func TestNew(t *testing.T) {
	ensure := ensure.New(t)

	imports := uniqpkg.New()
	ensure(imports).Equals(&uniqpkg.UniquePackagePaths{})
}

func TestGeneratePackageName(t *testing.T) {
	ensure := ensure.New(t)

	type packageResult struct {
		Path    string
		Imports []*uniqpkg.ImportDetails
	}
	type result map[string]*packageResult

	table := []struct {
		Name string

		Runner func(p *uniqpkg.UniquePackagePaths)

		ExpectedResult result

		Subject *uniqpkg.UniquePackagePaths
	}{
		{
			Name: "inside one package scope with single imported package referenced once",

			Runner: func(p *uniqpkg.UniquePackagePaths) {
				scopePkg1 := &packages.Package{PkgPath: "pkgs/scopepkg1"}
				p.GeneratePackageName(scopePkg1, types.NewPackage("pkgs/importedpkg1", "importedpkg1"))
			},

			ExpectedResult: result{
				"pkgs/scopepkg1": {
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
				scopePkg1 := &packages.Package{PkgPath: "pkgs/scopepkg1"}
				importedPkg1 := types.NewPackage("pkgs/importedpkg1", "importedpkg1")

				p.GeneratePackageName(scopePkg1, importedPkg1)
				p.GeneratePackageName(scopePkg1, importedPkg1)
				p.GeneratePackageName(scopePkg1, importedPkg1)
			},

			ExpectedResult: result{
				"pkgs/scopepkg1": {
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
				scopePkg1 := &packages.Package{PkgPath: "pkgs/scopepkg1"}
				importedPkg1 := types.NewPackage("pkgs/importedpkg1", "importedpkg1")
				importedPkg2 := types.NewPackage("pkgs/importedpkg2", "importedpkg2")

				p.GeneratePackageName(scopePkg1, importedPkg1)
				p.GeneratePackageName(scopePkg1, importedPkg2)
				p.GeneratePackageName(scopePkg1, importedPkg1)
				p.GeneratePackageName(scopePkg1, importedPkg2)
			},

			ExpectedResult: result{
				"pkgs/scopepkg1": {
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
				scopePkg1 := &packages.Package{PkgPath: "pkgs/scopepkg1"}
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

			ExpectedResult: result{
				"pkgs/scopepkg1": {
					Path: "pkgs/scopepkg1",
					Imports: []*uniqpkg.ImportDetails{
						{Name: "importedpkg1", Path: "pkgs/importedpkg1"},
						{Name: "importedpkg2", Path: "pkgs/importedpkg2"},
					},
				},
				"pkgs/scopepkg2": {
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
				scopePkg1 := &packages.Package{PkgPath: "pkgs/scopepkg1"}
				scopePkg2 := &packages.Package{Name: "scopepkg2", PkgPath: "pkgs/scopepkg2"}

				p.GeneratePackageName(scopePkg1, types.NewPackage("pkgs/1/common", "common"))
				p.GeneratePackageName(scopePkg2, types.NewPackage("pkgs/2/common", "common"))
			},

			ExpectedResult: result{
				"pkgs/scopepkg1": {
					Path: "pkgs/scopepkg1",
					Imports: []*uniqpkg.ImportDetails{
						{Name: "common", Path: "pkgs/1/common"},
					},
				},
				"pkgs/scopepkg2": {
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
				scopePkg1 := &packages.Package{PkgPath: "pkgs/scopepkg1"}

				p.GeneratePackageName(scopePkg1, types.NewPackage("pkgs/1/common", "common"))
				p.GeneratePackageName(scopePkg1, types.NewPackage("pkgs/2/common", "common"))
			},

			ExpectedResult: result{
				"pkgs/scopepkg1": {
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
				scopePkg1 := &packages.Package{PkgPath: "pkgs/scopepkg1"}

				p.GeneratePackageName(scopePkg1, types.NewPackage("pkgs/1/common", "common"))
				p.GeneratePackageName(scopePkg1, types.NewPackage("pkgs/2/common", "common"))
				p.GeneratePackageName(scopePkg1, types.NewPackage("pkgs/3/common", "common"))

				// To show that the generated names are repeatable
				p.GeneratePackageName(scopePkg1, types.NewPackage("pkgs/2/common", "common"))
				p.GeneratePackageName(scopePkg1, types.NewPackage("pkgs/3/common", "common"))
				p.GeneratePackageName(scopePkg1, types.NewPackage("pkgs/1/common", "common"))
			},

			ExpectedResult: result{
				"pkgs/scopepkg1": {
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

	ensure.RunTableByIndex(table, func(ensure ensuring.E, i int) {
		entry := table[i]

		entry.Runner(entry.Subject)

		actual := result{}
		for k := range entry.ExpectedResult {
			pkg := entry.Subject.ForPackage(k)

			actual[k] = &packageResult{
				Path: pkg.Path,

				Imports: pkg.Imports(),
			}
		}

		ensure(actual).Equals(entry.ExpectedResult)
	})
}

func TestForPackage(t *testing.T) {
	ensure := ensure.New(t)

	const pkg1Path = "pkgs/pkg1"
	const pkg2Path = "pkgs/pkg2"

	pkgs := uniqpkg.New()

	// Add an import for pkg1
	{
		pkg := pkgs.ForPackage(pkg1Path)
		ensure(pkg.Path).Equals(pkg1Path)
		ensure(pkg.Imports()).IsEmpty()
		pkg.AddImport("pkgs/imported", "imported")
	}

	// Show the imports are separate between packages
	{
		pkg := pkgs.ForPackage(pkg2Path)
		ensure(pkg.Path).Equals(pkg2Path)
		ensure(pkg.Imports()).IsEmpty()
		pkg.AddImport("pkgs/other", "other")
		ensure(pkg.Imports()).Equals([]*uniqpkg.ImportDetails{{Path: "pkgs/other", Name: "other"}})
	}

	// Show that it finds the existing pkg1 correctly
	{
		pkg := pkgs.ForPackage(pkg1Path)
		ensure(pkg.Imports()).Equals([]*uniqpkg.ImportDetails{{Path: "pkgs/imported", Name: "imported"}})
		ensure(pkg.Path).Equals(pkg1Path)
	}
}

func TestAddImport(t *testing.T) {
	ensure := ensure.New(t)

	createPkg := func() *uniqpkg.Package {
		pkgs := uniqpkg.New()
		return pkgs.ForPackage("my/pkg")
	}

	ensure.Run("supports adding multiple imports", func(ensure ensuring.E) {
		pkg := createPkg()

		pkg.AddImport("pkgs/pkg1", "pkg1")
		pkg.AddImport("pkgs/pkg2", "pkg2")

		ensure(pkg.Imports()).Equals([]*uniqpkg.ImportDetails{
			{Path: "pkgs/pkg1", Name: "pkg1"},
			{Path: "pkgs/pkg2", Name: "pkg2"},
		})
	})

	ensure.Run("supports adding the same import several times", func(ensure ensuring.E) {
		pkg := createPkg()

		pkg.AddImport("pkgs/pkg1", "pkg1")
		pkg.AddImport("pkgs/pkg2", "pkg2")
		pkg.AddImport("pkgs/pkg1", "pkg1")
		pkg.AddImport("pkgs/pkg2", "pkg2")

		ensure(pkg.Imports()).Equals([]*uniqpkg.ImportDetails{
			{Path: "pkgs/pkg1", Name: "pkg1"},
			{Path: "pkgs/pkg2", Name: "pkg2"},
		})
	})

	ensure.Run("supports adding several imports with the same name", func(ensure ensuring.E) {
		pkg := createPkg()

		pkg.AddImport("pkgs/1/pkg1", "pkg1")
		pkg.AddImport("pkgs/1/pkg2", "pkg2")
		pkg.AddImport("pkgs/2/pkg1", "pkg1")
		pkg.AddImport("pkgs/2/pkg2", "pkg2")

		ensure(pkg.Imports()).Equals([]*uniqpkg.ImportDetails{
			{Path: "pkgs/1/pkg1", Name: "pkg1"},
			{Path: "pkgs/1/pkg2", Name: "pkg2"},
			{Path: "pkgs/2/pkg1", Name: "pkg12", IsAlias: true},
			{Path: "pkgs/2/pkg2", Name: "pkg22", IsAlias: true},
		})
	})
}

func TestImports(t *testing.T) {
	ensure := ensure.New(t)

	ensure.Run("imports are in sorted order", func(ensure ensuring.E) {
		pkgs := uniqpkg.New()
		pkg := pkgs.ForPackage("my/pkg")

		pkg.AddImport("zzz", "zzz")
		pkg.AddImport("aaa", "aaa")
		pkg.AddImport("ccc", "ccc")
		pkg.AddImport("bbb", "bbb")

		ensure(pkg.Imports()).Equals([]*uniqpkg.ImportDetails{
			{Path: "aaa", Name: "aaa"},
			{Path: "bbb", Name: "bbb"},
			{Path: "ccc", Name: "ccc"},
			{Path: "zzz", Name: "zzz"},
		})
	})
}
