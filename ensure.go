// Package ensure is a balanced testing framework.
// It supports modern Go 1.13+ error comparisons (via errors.Is), and provides easy to read diffs (via deep.Equal).
//
// Most of the implementation is in the ensurepkg package.
// ensure.New should be used to create an instance of the ensure framework,
// which allows shadowing the "ensure" package (like with the t variable in tests).
// This provides easy test refactoring, while still being able to access the underlying types via the ensurepkg package.
//
// For example:
//  func TestMyFunction(t *testing.T) {
//   ensure := ensure.New(t)
//   ...
//
//   t.Run("my subtest", func(t *testing.T) {
//	   ensure := ensure.New(t) // This is using the shadowed version of ensure, and can easily be refactored
//     ...
//
//     ensure("abc").Equals("abc") // To ensure a value is correct, use ensure as a function
//     ensure.Fail() // Methods can be called directly on ensure
//   })
//  }
package ensure

import "github.com/JosiahWitt/ensure/ensurepkg"

// New creates an instance of the ensure test framework using the current testing context.
func New(t ensurepkg.T) ensurepkg.Ensure {
	return ensurepkg.New(t)
}
