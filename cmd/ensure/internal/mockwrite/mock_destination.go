package mockwrite

import (
	"path/filepath"
	"strings"

	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ensurefile"
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
	rawPackagePath string
}

func computeMockDestinations(config *ensurefile.Config) (mockDestinations, error) {
	errGroup := erg.NewAs(ErrMultipleWriteFailures)

	destinations := mockDestinations{}
	for _, pkg := range config.Mocks.Packages {
		dest, err := computeMockDestination(config, pkg.Path)
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

func computeMockDestination(config *ensurefile.Config, pkgPath string) (*mockDestination, error) {
	const internalPart = "internal/"

	// Check if package is internal
	idx := strings.LastIndex(pkgPath, internalPart)
	if idx < 0 {
		return &mockDestination{
			PWD:            config.RootPath,
			MockDir:        config.Mocks.PrimaryDestination,
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
		rawPackagePath: pkgPathSuffix,
	}, nil
}

func (dest *mockDestination) fullPath() string {
	originalPackageName := filepath.Base(dest.rawPackagePath)
	mockPackageName := "mock_" + originalPackageName
	destPkgFile := filepath.Join(filepath.Dir(dest.rawPackagePath), mockPackageName, mockPackageName+".go")

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
