package mockwrite

import (
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ensurefile"
	"github.com/JosiahWitt/erk"
)

type (
	ErkUnableToTidy struct{ erk.DefaultKind }
)

var (
	ErrTidyUnableToList    = erk.New(ErkUnableToTidy{}, "Could not list files recursively for '{{.path}}': {{.err}}")
	ErrTidyUnableToCleanup = erk.New(ErkUnableToTidy{}, "Could not delete '{{.path}}': {{.err}}")
)

// TidyMocks removes any files other than those that are expected to exist in the mock directories.
func (g *MockWriter) TidyMocks(config *ensurefile.Config) error {
	mockDestinations, err := computeMockDestinations(config)
	if err != nil {
		return err
	}

	mockDestinationsByMockDir := mockDestinations.byFullMockDir()
	pathsToDelete := []string{}

	for mockDir, mockDests := range mockDestinationsByMockDir {
		recursivePaths, err := g.FSWrite.ListRecursive(mockDir)
		if err != nil {
			return erk.WrapWith(ErrTidyUnableToList, err, erk.Params{
				"path": mockDir,
			})
		}

		// Any recursive path that isn't a prefix to a mock destination can be deleted
		for _, recursivePath := range recursivePaths {
			if !mockDests.hasFullPathPrefix(recursivePath) {
				pathsToDelete = append(pathsToDelete, recursivePath)
			}
		}
	}

	if len(pathsToDelete) > 0 {
		g.Logger.Println("Tidying mocks:")
		for _, pathToDelete := range pathsToDelete {
			g.Logger.Printf(" - Removing: %s\n", pathToDelete)

			if err := g.FSWrite.RemoveAll(pathToDelete); err != nil {
				return erk.WrapWith(ErrTidyUnableToCleanup, err, erk.Params{
					"path": pathToDelete,
				})
			}
		}
	}

	return nil
}
