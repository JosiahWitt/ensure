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
	runConfig{
		prepare: func(ensure ensuring.E) func(string, func(ensuring.E)) {
			return ensure.Run
		},
	}.test(t)
}

func TestERunParallel(t *testing.T) {
	runConfig{
		isParallel: true,
		prepare: func(ensure ensuring.E) func(string, func(ensuring.E)) {
			return ensure.RunParallel
		},
	}.test(t)
}

type runConfig struct {
	isParallel bool
	isSync     bool
	prepare    func(ensure ensuring.E) func(string, func(ensuring.E))
}

func (cfg runConfig) test(t *testing.T) {
	ctrl := gomock.NewController(t)

	const name = "my test name"

	outerMockT := setupMockTWithCleanupCheck(t)
	outerMockT.EXPECT().Helper()

	outerMockCtx := mock_testctx.NewMockContext(ctrl)
	testhelper.SetTestContext(t, outerMockT, outerMockCtx)

	innerMockT := setupMockTWithCleanupCheck(t)
	innerMockT.EXPECT().Helper().Times(2)

	if cfg.isParallel {
		innerMockT.EXPECT().Parallel()
	}

	innerMockCtx := mock_testctx.NewMockContext(ctrl)
	innerMockCtx.EXPECT().T().Return(innerMockT)
	testhelper.SetTestContext(t, innerMockT, innerMockCtx)

	if cfg.isSync {
		preSyncInnerMockT := setupMockT(t)
		preSyncInnerMockT.EXPECT().Helper()

		preSyncInnerMockCtx := mock_testctx.NewMockSyncableContext(ctrl)
		preSyncInnerMockCtx.EXPECT().T().Return(preSyncInnerMockT)
		testhelper.SetTestContext(t, preSyncInnerMockT, preSyncInnerMockCtx)

		outerMockCtx.EXPECT().Run(name, gomock.Any()).
			Do(execFuncParamWithName(preSyncInnerMockCtx))

		preSyncInnerMockCtx.EXPECT().Sync(gomock.Any()).
			Do(execFuncParam(innerMockCtx))
	} else {
		outerMockCtx.EXPECT().Run(name, gomock.Any()).
			Do(execFuncParamWithName(innerMockCtx))
	}

	var innerEnsure ensuring.E
	outerEnsure := ensure.New(outerMockT)
	run := cfg.prepare(outerEnsure)
	run(name, func(ensure ensuring.E) {
		innerEnsure = ensure
	})

	innerMockT.EXPECT().Fatalf("inner")
	innerEnsure.Failf("inner")
}

func execFuncParam(param testctx.Context) func(fn func(ctx testctx.Context)) {
	return func(fn func(ctx testctx.Context)) {
		fn(param)
	}
}

func execFuncParamWithName(param testctx.Context) func(name string, fn func(ctx testctx.Context)) {
	return func(name string, fn func(ctx testctx.Context)) {
		fn(param)
	}
}
