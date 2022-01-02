// Package uniqpkg supports generating unique package names by creating aliases as necessary.
package uniqpkg

import (
	"fmt"
	"go/types"
	"sort"

	"golang.org/x/tools/go/packages"
)

// UniquePackagePaths supports generating unique package names by creating aliases as necessary.
type UniquePackagePaths struct {
	packagePaths map[string]*internalPackage
}

type internalPackage struct {
	name string
	path string

	byPath map[string]*ImportDetails
	byName map[string]*ImportDetails
}

// Package stores the imports for a single package.
type Package struct {
	Name    string
	Path    string
	Imports []*ImportDetails
}

// ImportDetails stores the details of a single import within a package.
type ImportDetails struct {
	Name    string
	Path    string
	IsAlias bool // Indicates that the name is not the default name, and thus should be specified as an alias.
}

// GeneratePackageName supports generating unique package names by creating aliases as necessary.
func (p *UniquePackagePaths) GeneratePackageName(scopePackage *packages.Package, importedPackage *types.Package) string {
	rawPkg := p.getByScope(scopePackage)

	details := rawPkg.byPath[importedPackage.Path()]
	if details == nil {
		details = p.buildImportDetails(rawPkg, importedPackage)
		rawPkg.byName[details.Name] = details
		rawPkg.byPath[details.Path] = details
	}

	return details.Name
}

func (p *UniquePackagePaths) getByScope(scopePackage *packages.Package) *internalPackage {
	if p.packagePaths == nil {
		p.packagePaths = make(map[string]*internalPackage)
	}

	rawPkg := p.packagePaths[scopePackage.PkgPath]
	if rawPkg == nil {
		rawPkg = &internalPackage{
			name: scopePackage.Name,
			path: scopePackage.PkgPath,

			byPath: make(map[string]*ImportDetails),
			byName: make(map[string]*ImportDetails),
		}

		p.packagePaths[scopePackage.PkgPath] = rawPkg
	}

	return rawPkg
}

func (p *UniquePackagePaths) buildImportDetails(rawPkg *internalPackage, importedPackage *types.Package) *ImportDetails {
	pkgName := importedPackage.Name()
	pkgPath := importedPackage.Path()

	// Start at 1, so the first alias starts with a suffix of 2.
	// eg. mypkg, mypkg2, mypkg3, etc.
	for i := 1; ; i++ {
		isAlias := i != 1 // After one iteration, we are creating an alias

		generatedPkgName := pkgName
		if isAlias {
			generatedPkgName = fmt.Sprintf("%s%d", pkgName, i)
		}

		// When we find a name that doesn't exist yet, we use it as the import details
		if details := rawPkg.byName[generatedPkgName]; details == nil {
			return &ImportDetails{
				Name:    generatedPkgName,
				Path:    pkgPath,
				IsAlias: isAlias,
			}
		}
	}
}

// Export the packages into the external formats.
func (p *UniquePackagePaths) Export() []*Package {
	pkgs := make([]*Package, 0, len(p.packagePaths))

	for _, rawPkg := range p.packagePaths {
		pkgs = append(pkgs, rawPkg.export())
	}

	sort.Slice(pkgs, func(i, j int) bool {
		return pkgs[i].Path < pkgs[j].Path
	})

	return pkgs
}

func (rawPkg *internalPackage) export() *Package {
	imports := make([]*ImportDetails, 0, len(rawPkg.byPath))

	for _, details := range rawPkg.byPath {
		imports = append(imports, details)
	}

	sort.Slice(imports, func(i, j int) bool {
		return imports[i].Path < imports[j].Path
	})

	return &Package{
		Name: rawPkg.name,
		Path: rawPkg.path,

		Imports: imports,
	}
}
