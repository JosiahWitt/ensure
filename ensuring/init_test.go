//nolint:testpackage // Only used for the init function below.
package ensuring

import "github.com/JosiahWitt/ensure/ensuring/internal/testhelper"

//nolint:gochecknoinits // Only to make testing easier.
func init() {
	// Initializes the unexported newTestContextFunc variable to use the test implementation.
	// This allows us to continue to keep the tests in the separate testing package and keep
	// the newTestContextFunc variable unexported.
	newTestContextFunc = testhelper.NewTestContext
}
