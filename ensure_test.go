package ensure_test

import (
	"testing"

	"github.com/JosiahWitt/ensure"
	"github.com/JosiahWitt/ensure/internal/mocks/mock_testctx"
	"go.uber.org/mock/gomock"
)

func TestNew(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockT := mock_testctx.NewMockT(ctrl)
	ensure := ensure.New(mockT)

	if ensure == nil {
		t.Error("expected ensure not to be nil")
	}
}
