// Package ensure is a balanced testing framework for Go 1.14+.
// It supports modern Go 1.13+ error comparisons (via errors.Is), and provides easy to read diffs (via deep.Equal).
//
// Most of the implementation is in the ensurer package.
// ensure.New should be used to create an instance of the ensure framework,
// which allows shadowing the "ensure" package (like with the t variable in tests).
// This provides easy test refactoring, while still being able to access the underlying types via the ensurer package.
//
// For example:
//
//	func TestBasicExample(t *testing.T) {
//	 ensure := ensure.New(t)
//	 ...
//
//	 // Methods can be called on ensure, for example, Run:
//	 ensure.Run("my subtest", func(ensure ensurer.Ensure) {
//	   ...
//
//	 	 // To ensure a value is correct, use ensure as a function:
//	 	 ensure("abc").Equals("abc")
//	 	 ensure(produceError()).IsError(expectedError)
//	 	 ensure(doNotProduceError()).IsNotError()
//	 	 ensure(true).IsTrue()
//	 	 ensure(false).IsFalse()
//	 	 ensure("").IsEmpty()
//
//	   // Failing a test directly:
//	   ensure.Failf("Something went wrong, and we stop the test immediately")
//	 })
//	}
package ensure

import "github.com/JosiahWitt/ensure/ensurer"

// New creates an instance of the ensure test framework using the current testing context.
func New(t ensurer.T) ensurer.E {
	return ensurer.InternalCreateDoNotCallDirectly(t)
}
