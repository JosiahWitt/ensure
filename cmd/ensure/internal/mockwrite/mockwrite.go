// Package mockwrite writes the provided mocks to the file system.
package mockwrite

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ensurefile"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/fswrite"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/mockgen"
	"github.com/JosiahWitt/erk"
	"github.com/JosiahWitt/erk/erg"
)

type (
	ErkInvalidConfig    struct{ erk.DefaultKind }
	ErkMultipleFailures struct{ erk.DefaultKind }
	ErkFSWriteError     struct{ erk.DefaultKind }
)

var (
	ErrMissingMockConfig = erk.New(ErkInvalidConfig{}, "Missing `mocks` config in .ensure.yml file. For example:\n\n"+ensurefile.ExampleFile)

	ErrMissingPackages = erk.New(ErkInvalidConfig{},
		"No mocks to generate. Please add some to `mocks.packages` in .ensure.yml file. For example:\n\n"+ensurefile.ExampleFile,
	)
	ErrDuplicatePackagePath = erk.New(ErkInvalidConfig{}, "Found duplicate package path: {{.packagePath}}. Package paths must be unique.")

	ErrMultipleWriteFailures = erk.New(ErkMultipleFailures{}, "Unable to write at least one mock")

	ErrUnableToCreateDir  = erk.New(ErkFSWriteError{}, "Could not create directory '{{.path}}': {{.err}}")
	ErrUnableToCreateFile = erk.New(ErkFSWriteError{}, "Could not create file '{{.path}}': {{.err}}")
)

// Writable writes the provided mocks to the file system.
type Writable interface {
	WriteMocks(config *ensurefile.Config, mocks []*mockgen.PackageMock) error
	TidyMocks(config *ensurefile.Config) error
}

// MockWriter writes the provided mocks to the file system.
type MockWriter struct {
	FSWrite fswrite.Writable
	Logger  *log.Logger
}

var _ Writable = &MockWriter{}

// WriteMocks writes the provided mocks to the file system.
func (w *MockWriter) WriteMocks(config *ensurefile.Config, mocks []*mockgen.PackageMock) error {
	w.Logger.Println("Writing mocks:")
	errs := erg.NewAs(ErrMultipleWriteFailures)

	for _, mock := range mocks {
		if err := w.writeMock(config, mock); err != nil {
			errs = erg.Append(errs, err)
			continue
		}
	}

	if erg.Any(errs) {
		return errs
	}

	return nil
}

func (w *MockWriter) writeMock(config *ensurefile.Config, mock *mockgen.PackageMock) error {
	mockDest, err := computeMockDestination(config, mock.Package.Path)
	if err != nil {
		return err
	}

	mockFilePath := mockDest.fullPath()
	mockDirPath := filepath.Dir(mockFilePath)

	ifaceNames := strings.Join(mock.Package.InterfaceNames(), ",")
	w.Logger.Printf(" - Writing mocks: %s:%s -> %s\n", mock.Package.Path, ifaceNames, mockFilePath)

	if err := w.FSWrite.MkdirAll(mockDirPath, 0775); err != nil { //nolint:gomnd
		return erk.WrapWith(ErrUnableToCreateDir, err, erk.Params{
			"path": mockDirPath,
		})
	}

	if err := w.FSWrite.WriteFile(mockFilePath, mock.FileContents, 0664); err != nil { //nolint:gomnd
		return erk.WrapWith(ErrUnableToCreateFile, err, erk.Params{
			"path": mockFilePath,
		})
	}

	return nil
}
