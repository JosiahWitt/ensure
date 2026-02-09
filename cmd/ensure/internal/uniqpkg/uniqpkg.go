// Package uniqpkg supports generating unique package names by creating aliases as necessary.
package uniqpkg

import (
	"cmp"
	"fmt"
	"go/types"
	"slices"

	"golang.org/x/tools/go/packages"
)

// UniquePackagePaths supports generating unique package names by creating aliases as necessary.
type UniquePackagePaths struct {
	packages map[string]*Package
}

// Package stores the imports for a single package.
type Package struct {
	Path string

	byPath map[string]*ImportDetails
	byName map[string]*ImportDetails
}

// ImportDetails stores the details of a single import within a package.
type ImportDetails struct {
	Name    string
	Path    string
	IsAlias bool // Indicates that the name is not the default name, and thus should be specified as an alias.
}

func New() *UniquePackagePaths {
	return &UniquePackagePaths{}
}

// GeneratePackageName supports generating unique package names by creating aliases as necessary.
// This method is expected to be called from within ifacereader, as it satisfies the PackageNameGenerator interface.
func (p *UniquePackagePaths) GeneratePackageName(scopePackage *packages.Package, importedPackage *types.Package) string {
	pkg := p.ForPackage(scopePackage.PkgPath)
	details := pkg.AddImport(importedPackage.Path(), importedPackage.Name())
	return details.Name
}

// ForPackage finds or adds the package identified by the path.
func (p *UniquePackagePaths) ForPackage(path string) *Package {
	if p.packages == nil {
		p.packages = make(map[string]*Package)
	}

	pkg := p.packages[path]
	if pkg == nil {
		pkg = &Package{
			Path: path,

			byPath: make(map[string]*ImportDetails),
			byName: make(map[string]*ImportDetails),
		}

		p.packages[path] = pkg
	}

	return pkg
}

// AddImport adds the package path and name as an import, generating an alias if necessary.
func (pkg *Package) AddImport(packagePath, packageName string) *ImportDetails {
	if details := pkg.byPath[packagePath]; details != nil {
		return details
	}

	// Start at 1, so the first alias starts with a suffix of 2.
	// eg. mypkg, mypkg2, mypkg3, etc.
	for i := 1; ; i++ {
		isAlias := i != 1 // After one iteration, we are creating an alias

		generatedPackageName := packageName
		if isAlias {
			generatedPackageName = fmt.Sprintf("%s%d", packageName, i)
		}

		// When we find a name that doesn't exist yet, we use it as the import details
		if details := pkg.byName[generatedPackageName]; details == nil {
			details := &ImportDetails{
				Name:    generatedPackageName,
				Path:    packagePath,
				IsAlias: isAlias,
			}

			pkg.byName[details.Name] = details
			pkg.byPath[details.Path] = details

			return details
		}
	}
}

// Imports returns the list of imports for the Package in sorted order.
func (pkg *Package) Imports() []*ImportDetails {
	imports := make([]*ImportDetails, 0, len(pkg.byPath))

	for _, details := range pkg.byPath {
		imports = append(imports, details)
	}

	slices.SortFunc(imports, func(a, b *ImportDetails) int {
		return cmp.Compare(a.Path, b.Path)
	})

	return imports
}
