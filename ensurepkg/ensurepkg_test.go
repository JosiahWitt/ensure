package ensurepkg_test

import (
	"testing"

	"github.com/JosiahWitt/ensure"
	"github.com/JosiahWitt/ensure/ensurepkg"
	"github.com/JosiahWitt/ensure/internal/mocks/mock_ensurepkg"
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
			mockT.EXPECT().Fatalf("Do not call `ensurepkg.InternalCreateDoNotCallDirectly(t)` directly. Instead use `ensure.New(t)`."),
		)

		ensurepkg.InternalCreateDoNotCallDirectly(mockT)
	})
}

func TestEnsureFailf(t *testing.T) {
	mockT := setupMockTWithCleanupCheck(t)

	const message = "my message %s"
	const arg1 = 123

	mockT.EXPECT().Fatalf(message, arg1).After(
		mockT.EXPECT().Helper(),
	)

	ensure := ensure.New(mockT)
	ensure.Failf(message, arg1)
}

func TestEnsureT(t *testing.T) {
	mockT := setupMockTWithCleanupCheck(t)

	ensure := ensure.New(mockT)
	if ensure.T() != mockT {
		t.Error("T() does not equal mockT")
	}
}

func TestEnsureCleanupCheck(t *testing.T) {
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
			mockT.EXPECT().Fatalf("Found ensure(<actual>) without chained assertion."),
		)
		cleanupFn()
	})
}

func setupMockT(t *testing.T) *mock_ensurepkg.MockT {
	t.Helper()
	ctrl := gomock.NewController(t)
	return mock_ensurepkg.NewMockT(ctrl)
}

func setupMockTWithCleanupCheck(t *testing.T) *mock_ensurepkg.MockT {
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
