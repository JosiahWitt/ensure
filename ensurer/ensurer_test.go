package ensurer_test

import (
	"testing"

	"github.com/JosiahWitt/ensure"
	"github.com/JosiahWitt/ensure/ensurer"
	"github.com/JosiahWitt/ensure/ensurer/internal/testhelper"
	"github.com/JosiahWitt/ensure/internal/mocks/mock_testctx"
	"github.com/JosiahWitt/ensure/internal/testctx"
	"github.com/golang/mock/gomock"
)

func TestGoVersion(t *testing.T) {
	t.Cleanup(func() {}) // This ensures the Cleanup function is present (Go 1.14+) so gomock controller tests don't fail silently
}

func TestNew(t *testing.T) {
	t.Run("when called through ensure.New", func(t *testing.T) {
		mockT := setupMockT(t)
		ensure := ensure.New(mockT)

		if ensure == nil {
			t.Error("expected ensure not to be nil")
		}
	})

	t.Run("when called directly", func(t *testing.T) {
		mockT := setupMockT(t)

		gomock.InOrder(
			mockT.EXPECT().Helper(),
			mockT.EXPECT().Fatalf("Do not call `ensurer.InternalCreateDoNotCallDirectly(t)` directly. Instead use `ensure := ensure.New(t)`."),
		)

		ensurer.InternalCreateDoNotCallDirectly(mockT)
	})
}

func TestNestedNew(t *testing.T) {
	t.Run("inner and outer scopes use the correct testing scope", func(t *testing.T) {
		outerMockT := setupMockT(t)
		outerMockT.EXPECT().Helper().AnyTimes()
		outerMockT.EXPECT().Cleanup(gomock.Any())
		outerEnsure := ensure.New(outerMockT)

		innerMockT := setupMockT(t)
		innerMockT.EXPECT().Helper().AnyTimes()
		innerMockT.EXPECT().Cleanup(gomock.Any())
		innerEnsure := outerEnsure.New(innerMockT)

		outerMockT.EXPECT().Fatalf("still outer")
		outerEnsure.Failf("still outer")

		innerMockT.EXPECT().Fatalf("inner")
		innerEnsure.Failf("inner")
	})

	t.Run("works with a real *testing.T", func(t *testing.T) {
		testhelper.AllowAnyTestContexts(t)
		originalName := t.Name()
		outerEnsure := ensure.New(t)

		t.Run("check nested ensure.New", func(t *testing.T) {
			innerEnsure := outerEnsure.New(t) // Uses the nested New method

			if outerEnsure.T().Name() == innerEnsure.T().Name() {
				t.Errorf("The testing context should not be the same between the inner and outer ensure")
			}
		})

		if outerEnsure.T().Name() != originalName {
			t.Errorf("The testing context should not be changed when outerEnsure.New() is used")
		}
	})
}

func TestEFailf(t *testing.T) {
	mockT := setupMockTWithCleanupCheck(t)

	const message = "my message %s"
	const arg1 = 123

	mockT.EXPECT().Fatalf(message, arg1).After(
		mockT.EXPECT().Helper(),
	)

	ensure := ensure.New(mockT)
	ensure.Failf(message, arg1)
}

func TestET(t *testing.T) {
	t.Run("when provided a non *testing.T instance", func(t *testing.T) {
		mockT := setupMockTWithCleanupCheck(t)

		gomock.InOrder(
			mockT.EXPECT().Helper(),
			mockT.EXPECT().Fatalf("An instance of *testing.T was not provided to ensure.New(t), thus T() cannot be used."),
		)

		ensure := ensure.New(mockT)
		ensure.T()
	})

	t.Run("when provided a *testing.T instance", func(t *testing.T) {
		testhelper.AllowAnyTestContexts(t)
		ensure := ensure.New(t)

		if ensure.T().Name() != t.Name() {
			t.Fatalf("Expected the same *testing.T instance as the one that was provided")
		}
	})

	t.Run("when provided a *testing.T instance and using ensure.Run", func(t *testing.T) {
		testhelper.AllowAnyTestContexts(t)
		ensure := ensure.New(t)
		outerName := t.Name()

		ensure.Run("inner", func(ensure ensurer.E) {
			if ensure.T().Name() != outerName+"/inner" {
				t.Fatalf("Expected to be able to use T() inside ensure.Run")
			}
		})
	})
}

func TestEGoMockController(t *testing.T) {
	mockT := setupMockTWithCleanupCheck(t)
	mockT.EXPECT().Cleanup(gomock.Any()).AnyTimes() // Setup by GoMock Controller
	mockT.EXPECT().Helper().AnyTimes()              // Setup by GoMock Controller

	ensure := ensure.New(mockT)
	firstController := ensure.GoMockController()
	if firstController == nil {
		t.Error("firstController == nil")
	}

	// Memoized across GoMockController() calls
	secondController := ensure.GoMockController()
	if firstController != secondController {
		t.Errorf("firstController != secondController: %p != %p", firstController, secondController)
	}
}

func TestECleanupCheck(t *testing.T) {
	t.Run("when test was run", func(t *testing.T) {
		mockT := setupMockT(t)

		var cleanupFn func()
		gomock.InOrder(
			mockT.EXPECT().Helper(),
			mockT.EXPECT().Cleanup(gomock.Any()).Do(func(fn func()) {
				cleanupFn = fn
			}),

			mockT.EXPECT().Helper(), // For IsTrue call
		)

		ensure := ensure.New(mockT)
		ensure(true).IsTrue()
		cleanupFn()
	})

	t.Run("when test was not run", func(t *testing.T) {
		mockT := setupMockT(t)

		var cleanupFn func()
		gomock.InOrder(
			mockT.EXPECT().Helper(),
			mockT.EXPECT().Cleanup(gomock.Any()).Do(func(fn func()) {
				cleanupFn = fn
			}),
		)

		ensure := ensure.New(mockT)
		ensure(true)

		gomock.InOrder(
			mockT.EXPECT().Helper(),
			mockT.EXPECT().Errorf("Found ensure(<actual>) without chained assertion."),
		)
		cleanupFn()
	})
}

func setupMockT(t *testing.T) *mock_testctx.MockT {
	t.Helper()
	ctrl := gomock.NewController(t)
	mockT := mock_testctx.NewMockT(ctrl)

	testhelper.SetTestContext(t, mockT, testctx.New(mockT))
	return mockT
}

func setupMockTWithCleanupCheck(t *testing.T) *mock_testctx.MockT {
	t.Helper()
	mockT := setupMockT(t)

	gomock.InOrder(
		mockT.EXPECT().Helper(),
		mockT.EXPECT().Cleanup(gomock.Any()).Do(func(fn func()) {
			t.Cleanup(fn)
		}),
	)

	return mockT
}
