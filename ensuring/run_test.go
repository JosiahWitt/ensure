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
	testERun(t, false, func(ensure ensuring.E) func(string, func(ensuring.E)) {
		return ensure.Run
	})
}

func TestERunParallel(t *testing.T) {
	testERun(t, true, func(ensure ensuring.E) func(string, func(ensuring.E)) {
		return ensure.RunParallel
	})
}

type runBuilderFunc func(ensure ensuring.E) func(string, func(ensuring.E))

func testERun(t *testing.T, isParallel bool, runBuilder runBuilderFunc) {
	ctrl := gomock.NewController(t)

	outerMockT := setupMockTWithCleanupCheck(t)
	outerMockT.EXPECT().Helper()

	outerMockCtx := mock_testctx.NewMockContext(ctrl)
	testhelper.SetTestContext(t, outerMockT, outerMockCtx)

	innerMockT := setupMockTWithCleanupCheck(t)
	innerMockT.EXPECT().Helper().Times(2)

	if isParallel {
		innerMockT.EXPECT().Parallel()
	}

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
	run := runBuilder(outerEnsure)
	run(name, func(ensure ensuring.E) {
		innerEnsure = ensure
	})

	innerMockT.EXPECT().Fatalf("inner")
	innerEnsure.Failf("inner")
}
