package mockwrite

import (
	"path"
	"path/filepath"
	"strings"

	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ensurefile"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader"
	"github.com/JosiahWitt/erk"
	"github.com/JosiahWitt/erk/erg"
)

type ErkMockDestination struct{ erk.DefaultKind }

var ErrInternalPackageOutsideModule = erk.New(ErkMockDestination{},
	"Cannot generate mock of internal package, since package '{{.packagePath}}' is not in the current module '{{.modulePath}}'",
)

type mockDestinations []*mockDestination

type mockDestination struct {
	PWD            string
	MockDir        string
	packageName    string
	rawPackagePath string
}

func computeMockDestinations(config *ensurefile.Config, packages []*ifacereader.Package) (mockDestinations, error) {
	errGroup := erg.NewAs(ErrMultipleWriteFailures)

	destinations := mockDestinations{}
	for _, pkg := range packages {
		dest, err := computeMockDestination(config, pkg.Name, pkg.Path)
		if err != nil {
			errGroup = erg.Append(errGroup, err)
			continue
		}

		destinations = append(destinations, dest)
	}

	if erg.Any(errGroup) {
		return nil, errGroup
	}

	return destinations, nil
}

func computeMockDestination(config *ensurefile.Config, pkgName, pkgPath string) (*mockDestination, error) {
	const internalPart = "internal/"

	// Check if package is internal
	idx := strings.LastIndex(pkgPath, internalPart)
	if idx < 0 {
		return &mockDestination{
			PWD:            config.RootPath,
			MockDir:        config.Mocks.PrimaryDestination,
			packageName:    pkgName,
			rawPackagePath: pkgPath,
		}, nil
	}

	if !strings.HasPrefix(pkgPath, config.ModulePath) {
		return nil, erk.WithParams(ErrInternalPackageOutsideModule, erk.Params{
			"packagePath": pkgPath,
			"modulePath":  config.ModulePath,
		})
	}

	// Remove both the module path prefix, and the last internal/... suffix
	pkgPathPrefix := strings.TrimPrefix(pkgPath[:idx], config.ModulePath)

	// Everything after the last internal/...
	pkgPathSuffix := pkgPath[idx+len(internalPart):]

	return &mockDestination{
		PWD:            filepath.Join(config.RootPath, pkgPathPrefix),
		MockDir:        filepath.Join(internalPart, config.Mocks.InternalDestination),
		packageName:    pkgName,
		rawPackagePath: pkgPathSuffix,
	}, nil
}

func (dest *mockDestination) fullPath() string {
	packagePathPrefix, pathPackageName := path.Split(dest.rawPackagePath)

	// If the last part of the package path doesn't match the package name, keep the full package path as the prefix.
	// For example, for something like github.com/xyz/abc/v2, we'll generate the package as github.com/xyz/abc/v2/mock_abc.
	// This does mean that github.com/xyz/abc/v2 and github.com/xyz/abc/v2/abc would clash, but hopefully that's unlikely in the real world.
	if pathPackageName != dest.packageName {
		packagePathPrefix = filepath.Join(packagePathPrefix, pathPackageName)
	}

	mockPackageName := "mock_" + dest.packageName
	destPkgFile := filepath.Join(packagePathPrefix, mockPackageName, mockPackageName+".go")

	return filepath.Join(dest.PWD, dest.MockDir, destPkgFile)
}

func (dests mockDestinations) byFullMockDir() map[string]mockDestinations {
	byMockDir := map[string]mockDestinations{}
	for _, dest := range dests {
		key := filepath.Join(dest.PWD, dest.MockDir)
		byMockDir[key] = append(byMockDir[key], dest)
	}

	return byMockDir
}

func (dests mockDestinations) hasFullPathPrefix(prefix string) bool {
	for _, dest := range dests {
		if strings.HasPrefix(dest.fullPath(), prefix) {
			return true
		}
	}

	return false
}
