// Package ensurepkg contains the implementation for the ensure test framework.
// Use ensure.New to create a new instance of Ensure.
//
// Deprecated: Use the ensuring package instead.
package ensurepkg

import "github.com/JosiahWitt/ensure/ensuring"

// T implements a subset of methods on [testing.T].
// More methods may be added to T with a minor ensure release.
//
// Deprecated: Use [ensuring.T] instead.
type T = ensuring.T

// Ensure ensures the actual value is correct using [Chain].
// Ensure also has methods that can be called directly.
//
// Deprecated: Use [ensuring.E] instead.
type Ensure = ensuring.E

// Chain chains assertions to the ensure function call.
//
// Deprecated: Use [ensuring.Chain] instead.
type Chain = ensuring.Chain
