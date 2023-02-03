//nolint:testpackage // Only used for the init function below.
package ensurepkg

import "github.com/JosiahWitt/ensure/ensurepkg/internal/testhelper"

//nolint:gochecknoinits // Only to make testing easier.
func init() {
	// Initializes the unexported newTestContext variable to use the test implementation.
	// This allows us to continue to keep the tests in the separate testing package and
	// keep the newTestContext variable unexported.
	newTestContext = testhelper.NewTestContext
}
