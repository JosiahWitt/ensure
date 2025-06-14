//go:build go1.25

package testctx

import (
	"testing"
	"testing/synctest"
)

func TestBaseSync(t *testing.T) {
	sync.sync(t, func(t *testing.T) {
		// Shows synctest.Test was called, since Wait panics if it's not in a "bubble"
		synctest.Wait()
	})
}

func SetupMockSync(t *testing.T) *MockSync {
	originalSync := sync
	t.Cleanup(func() {
		sync = originalSync
	})

	s := &MockSync{}
	sync = s
	return s
}

type MockSync struct {
	Calls []T
}

func (s *MockSync) sync(t T, fn func(t *testing.T)) {
	s.Calls = append(s.Calls, t)
	fn(&testing.T{})
}
