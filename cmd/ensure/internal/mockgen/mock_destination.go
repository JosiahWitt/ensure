package mockgen

import (
	"path/filepath"
	"strings"

	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ensurefile"
	"github.com/JosiahWitt/erk"
)

type ErkMockDestination struct{ erk.DefaultKind }

var ErrInternalPackageOutsideModule = erk.New(ErkMockDestination{},
	"Cannot generate mock of internal package, since package '{{.packagePath}}' is not in the current module '{{.modulePath}}'",
)

type mockDestination struct {
	PWD            string
	MockDir        string
	rawPackagePath string
}

func computeMockDestination(config *ensurefile.Config, pkg *ensurefile.Package) (*mockDestination, error) {
	const internalPart = "internal/"

	// Check if package is internal
	idx := strings.LastIndex(pkg.Path, internalPart)
	if idx < 0 {
		return &mockDestination{
			PWD:            config.RootPath,
			MockDir:        config.Mocks.PrimaryDestination,
			rawPackagePath: pkg.Path,
		}, nil
	}

	if !strings.HasPrefix(pkg.Path, config.ModulePath) {
		return nil, erk.WithParams(ErrInternalPackageOutsideModule, erk.Params{
			"packagePath": pkg.Path,
			"modulePath":  config.ModulePath,
		})
	}

	// Remove both the module path prefix, and the last internal/... suffix
	pkgPathPrefix := strings.TrimPrefix(pkg.Path[:idx], config.ModulePath)

	// Everything after the last internal/...
	pkgPathSuffix := pkg.Path[idx+len(internalPart):]

	return &mockDestination{
		PWD:            filepath.Join(config.RootPath, pkgPathPrefix),
		MockDir:        filepath.Join(internalPart, config.Mocks.InternalDestination),
		rawPackagePath: pkgPathSuffix,
	}, nil
}

func (dest *mockDestination) fullPath() string {
	originalPackageName := filepath.Base(dest.rawPackagePath)
	mockPackageName := "mock_" + originalPackageName
	destPkgFile := filepath.Join(filepath.Dir(dest.rawPackagePath), mockPackageName, mockPackageName+".go")

	return filepath.Join(dest.PWD, dest.MockDir, destPkgFile)
}
