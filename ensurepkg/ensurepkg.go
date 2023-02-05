// Package ensurepkg contains the implementation for the ensure test framework.
// Use ensure.New to create a new instance of Ensure.
//
// Deprecated: Use the ensurer package instead.
package ensurepkg

import "github.com/JosiahWitt/ensure/ensurer"

// T implements a subset of methods on [testing.T].
// More methods may be added to T with a minor ensure release.
//
// Deprecated: Use [ensurer.T] instead.
type T = ensurer.T

// Ensure ensures the actual value is correct using [Chain].
// Ensure also has methods that can be called directly.
//
// Deprecated: Use [ensurer.E] instead.
type Ensure = ensurer.E

// Chain chains assertions to the ensure function call.
//
// Deprecated: Use [ensurer.Chain] instead.
type Chain = ensurer.Chain
