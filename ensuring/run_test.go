package ensuring_test

import (
	"testing"

	"github.com/JosiahWitt/ensure"
	"github.com/JosiahWitt/ensure/ensuring"
	"github.com/JosiahWitt/ensure/ensuring/internal/testhelper"
	"github.com/JosiahWitt/ensure/internal/mocks/mock_testctx"
	"github.com/JosiahWitt/ensure/internal/testctx"
	"github.com/golang/mock/gomock"
)

func TestERun(t *testing.T) {
	ctrl := gomock.NewController(t)

	outerMockT := setupMockTWithCleanupCheck(t)
	outerMockT.EXPECT().Helper()

	outerMockCtx := mock_testctx.NewMockContext(ctrl)
	testhelper.SetTestContext(t, outerMockT, outerMockCtx)

	innerMockT := setupMockTWithCleanupCheck(t)
	innerMockT.EXPECT().Helper().Times(2)

	innerMockCtx := mock_testctx.NewMockContext(ctrl)
	innerMockCtx.EXPECT().T().Return(innerMockT)
	testhelper.SetTestContext(t, innerMockT, innerMockCtx)

	const name = "my test name"

	outerMockCtx.EXPECT().Run(name, gomock.Any()).
		Do(func(name string, fn func(ctx testctx.Context)) {
			fn(innerMockCtx)
		})

	var innerEnsure ensuring.E
	outerEnsure := ensure.New(outerMockT)
	outerEnsure.Run(name, func(ensure ensuring.E) {
		innerEnsure = ensure
	})

	innerMockT.EXPECT().Fatalf("inner")
	innerEnsure.Failf("inner")
}
