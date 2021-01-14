package ensurepkg_test

import (
	"testing"

	"github.com/JosiahWitt/ensure"
	"github.com/JosiahWitt/ensure/ensurepkg"
	"github.com/JosiahWitt/ensure/internal/mocks/mock_ensurepkg"
	"github.com/golang/mock/gomock"
)

func TestEnsureRun(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockT := mock_ensurepkg.NewMockT(ctrl)
	mockT.EXPECT().Helper().Times(2)

	const name = "my test name"
	providedTestingInput := &testing.T{}

	mockT.EXPECT().Run(name, gomock.Any()).
		Do(func(name string, fn func(t *testing.T)) {
			fn(providedTestingInput)
		})

	var actualParam ensurepkg.Ensure
	ensure := ensure.New(mockT)
	ensure.Run(name, func(ensure ensurepkg.Ensure) {
		actualParam = ensure
	})

	if actualParam == nil || actualParam.T() != providedTestingInput {
		t.Errorf("Expected providedTestingInput to be wrapped by ensure")
	}
}
