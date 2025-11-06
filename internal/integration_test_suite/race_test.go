package integration_test_suite_test

import (
	"sync"
	"testing"

	"github.com/JosiahWitt/ensure"
	"go.uber.org/mock/gomock"
)

func TestGoMockControllerConcurrency(t *testing.T) {
	ensure := ensure.New(t)

	const numGoroutines = 10
	controllers := make(chan *gomock.Controller, numGoroutines)

	var wg sync.WaitGroup
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			controllers <- ensure.GoMockController()
		}()
	}

	wg.Wait()
	close(controllers)

	var first *gomock.Controller
	for controller := range controllers {
		if first == nil {
			first = controller
		} else if controller != first {
			t.Errorf("Not all controllers are the same instance")
		}
	}
}
