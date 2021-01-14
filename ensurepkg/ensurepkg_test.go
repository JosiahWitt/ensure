package ensurepkg_test

import (
	"strings"
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
		ctrl := gomock.NewController(t)
		mockT := mock_ensurepkg.NewMockT(ctrl)
		ensure := ensure.New(mockT)

		if ensure == nil {
			t.Error("expected ensure not to be nil")
		}
	})

	t.Run("when called directly", func(t *testing.T) {
		defer func() {
			rec := recover()
			if rec == nil {
				t.Error("expected panic, got none")
				return
			}

			if !strings.HasPrefix(rec.(string), "Do not call ensurepkg.New directly. Instead use ensure.New. Called ensurepkg.New from:") {
				t.Error("incorrect panic message")
			}
		}()

		ctrl := gomock.NewController(t)
		mockT := mock_ensurepkg.NewMockT(ctrl)
		ensurepkg.New(mockT)
	})
}

func TestEnsureFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockT := mock_ensurepkg.NewMockT(ctrl)
	mockT.EXPECT().Fail().After(
		mockT.EXPECT().Helper(),
	)

	ensure := ensure.New(mockT)
	ensure.Fail()
}

func TestEnsureT(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockT := mock_ensurepkg.NewMockT(ctrl)

	ensure := ensure.New(mockT)
	if ensure.T() != mockT {
		t.Error("T() does not equal mockT")
	}
}
