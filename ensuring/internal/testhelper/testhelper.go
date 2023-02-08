// Package testhelper implements helpers for testing ensure.
package testhelper

import (
	"fmt"
	"testing"

	"github.com/JosiahWitt/ensure/internal/testctx"
)

//nolint:gochecknoglobals // Only used for testing.
var (
	testContexts         = map[testctx.T]testctx.Context{}
	allowAnyTestContexts = false

	checkingWrapEnsure = false // Used to prevent an infinite loop
)

// NewTestContext is called instead of [testctx.New] and is setup in ../../init_test.go.
// This shouldn't be used by anything else.
func NewTestContext(t testctx.T, wrapEnsure testctx.WrapEnsure) testctx.Context {
	checkWrapEnsure(t, wrapEnsure)

	ctx, ok := testContexts[t]
	if !ok {
		if allowAnyTestContexts {
			return testctx.New(t, wrapEnsure)
		}

		panic("Missing mock test context")
	}

	return ctx
}

// SetTestContext connects the provided targetT to the provided context to be surfaced by [NewTestContext].
// It is disconnected when t goes out of scope.
func SetTestContext(t *testing.T, targetT testctx.T, ctx testctx.Context) {
	t.Helper()
	testContexts[targetT] = ctx
	t.Cleanup(func() {
		delete(testContexts, targetT)
	})
}

// AllowAnyTestContexts allows any test contexts to be used for the scope of t.
// It causes [NewTestContext] to fallback to the default [testctx.New] implementation.
func AllowAnyTestContexts(t *testing.T) {
	t.Helper()
	allowAnyTestContexts = true
	t.Cleanup(func() {
		allowAnyTestContexts = false
	})
}

func checkWrapEnsure(t testctx.T, wrapEnsure testctx.WrapEnsure) {
	// Prevent an infinite loop, since NewTestContext will be called by wrapEnsure
	if checkingWrapEnsure {
		return
	}

	checkingWrapEnsure = true
	defer func() { checkingWrapEnsure = false }()

	ensure := wrapEnsure(t)
	if ensure == nil || fmt.Sprintf("%T", ensure) != "ensuring.E" {
		panic(fmt.Sprintf("wrapEnsure doesn't function correctly: %[1]v (%[1]T)", ensure))
	}
}
