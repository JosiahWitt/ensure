package ensurepkg_test

import (
	"testing"

	"github.com/JosiahWitt/ensure"
	"github.com/JosiahWitt/ensure/ensurepkg"
	"github.com/JosiahWitt/ensure/tests/mocks/mock_ensurepkg"
	"github.com/golang/mock/gomock"
)

func TestEnsureRun(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockT := mock_ensurepkg.NewMockT(ctrl)
	mockT.EXPECT().Helper().Times(2)

	const name = "my test name"
	providedT := &testing.T{}

	mockT.EXPECT().Run(name, gomock.Any()).
		Do(func(name string, fn func(t *testing.T)) {
			fn(providedT)
		})

	var actualParam ensurepkg.Ensure
	fn := func(ensure ensurepkg.Ensure) {
		actualParam = ensure
	}

	ensure := ensure.New(mockT)
	ensure.Run(name, fn)

	if actualParam == nil || actualParam.T() != providedT {
		t.Errorf("Expected providedT to be wrapped by ensure")
	}
}
