//go:build go1.25

package testctx_test

import (
	"fmt"
	"testing"

	"github.com/JosiahWitt/ensure/internal/mocks/mock_testctx"
	"github.com/JosiahWitt/ensure/internal/testctx"
	"github.com/golang/mock/gomock"
)

func TestSync(t *testing.T) {
	mockSync := testctx.SetupMockSync(t)

	ctrl := gomock.NewController(t)
	outerT := mock_testctx.NewMockT(ctrl)
	outerT.EXPECT().Helper()

	wrapEnsure := func(t testctx.T) interface{} { return fmt.Sprintf("%T", t) }
	ctx := testctx.New(outerT, wrapEnsure)

	var called bool
	ctx.(testctx.SyncableContext).Sync(func(ctx testctx.Context) {
		actualInnerT := ctx.T()
		neq(t, actualInnerT, nil)          // It shouldn't be nil, indicating the callback wasn't called
		neq(t, actualInnerT, &testing.T{}) // It shouldn't be empty, indicating Helper() wasn't called
		neq(t, actualInnerT, outerT)       // It shouldn't be the outerT

		// Show wrapEnsure was promoted correctly
		eq(t, ctx.Ensure(), "*testing.T")

		called = true
	})

	// Verify the callback was actually executed
	eq(t, called, true)

	eq(t, len(mockSync.Calls), 1)
	eq(t, mockSync.Calls[0], outerT)
}
