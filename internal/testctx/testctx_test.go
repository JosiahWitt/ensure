package testctx_test

import (
	"reflect"
	"testing"

	"github.com/JosiahWitt/ensure/internal/mocks/mock_testctx"
	"github.com/JosiahWitt/ensure/internal/testctx"
	"github.com/golang/mock/gomock"
	"github.com/kr/pretty"
)

func TestNew(t *testing.T) {
	mockT := struct {
		testctx.T
		unique string
	}{unique: "hello"}

	ctx := testctx.New(mockT)
	eq(t, ctx.T(), mockT)
}

func TestT(t *testing.T) {
	mockT := struct {
		testctx.T
		unique string
	}{unique: "hello"}

	ctx := testctx.New(mockT)
	eq(t, ctx.T(), mockT)
}

func TestRun(t *testing.T) {
	ctrl := gomock.NewController(t)
	outerT := mock_testctx.NewMockT(ctrl)
	outerT.EXPECT().Helper()

	outerT.EXPECT().Run("everything works", gomock.Any()).Do(func(_ string, fn func(t *testing.T)) {
		fn(&testing.T{})
	})

	ctx := testctx.New(outerT)

	var actualInnerT *testing.T
	ctx.Run("everything works", func(ctx testctx.Context) {
		actualInnerT = ctx.T().(*testing.T)
	})

	neq(t, actualInnerT, nil)          // It shouldn't be nil, indicating the callback wasn't called
	neq(t, actualInnerT, &testing.T{}) // It shouldn't be empty, indicating Helper() wasn't called
	neq(t, actualInnerT, outerT)       // It shouldn't be the outerT
}

func TestGoMockController(t *testing.T) {
	t.Run("GoMock controller is memoized", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockT := mock_testctx.NewMockT(ctrl)

		mockT.EXPECT().Helper().MinTimes(2)
		mockT.EXPECT().Cleanup(gomock.Any()).Do(func(fn func()) {
			fn()
		}).Times(2) // We call it once and gomock.NewController calls it once

		ctx := testctx.New(mockT)
		mockCtrl := ctx.GoMockController()
		eq(t, mockCtrl.T, mockT)

		// Should return the same result if called twice
		mockCtrl2 := ctx.GoMockController()
		eq(t, mockCtrl2.T, mockT)

		eq(t, mockCtrl == mockCtrl2, true) // Both point to the same instance
	})

	t.Run("GoMock controller is finished", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockT := mock_testctx.NewMockT(ctrl)

		var cleanupFn func()
		gomock.InOrder(
			mockT.EXPECT().Helper(),
			mockT.EXPECT().Cleanup(gomock.Any()).Do(func(fn func()) {
				cleanupFn = fn
			}),
		)

		mockT.EXPECT().Helper().AnyTimes()
		mockT.EXPECT().Cleanup(gomock.Any()).AnyTimes()

		ctx := testctx.New(mockT)
		mockCtrl := ctx.GoMockController()

		// SomeMethod is never "called", and should be noticed during cleanup
		exampleType := &exampleTypeWithMethod{}
		mockCtrl.RecordCall(exampleType, "SomeMethod", true)

		mockT.EXPECT().Errorf(gomock.Any(), gomock.Any()).MinTimes(1)
		cleanupFn()
	})
}

func eq(t *testing.T, a, b interface{}) {
	t.Helper()
	if !reflect.DeepEqual(a, b) {
		t.Fatalf(pretty.Sprintf("% #v should equal % #v", a, b))
	}
}

func neq(t *testing.T, a, b interface{}) {
	t.Helper()
	if reflect.DeepEqual(a, b) {
		t.Fatalf(pretty.Sprintf("% #v should not equal % #v", a, b))
	}
}

type exampleTypeWithMethod struct{}

func (*exampleTypeWithMethod) SomeMethod(param bool) {}
