package ensurepkg_test

import (
	"testing"

	"github.com/JosiahWitt/ensure"
	"github.com/JosiahWitt/ensure/ensurepkg"
	"github.com/golang/mock/gomock"
)

func TestEnsureRun(t *testing.T) {
	mockT := setupMockTWithCleanupCheck(t)
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
		t.Fatalf("Expected providedTestingInput to be wrapped by ensure")
	}
}
