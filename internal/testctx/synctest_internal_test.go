//go:build go1.25

package testctx

import (
	"testing"
	"testing/synctest"
)

func TestBaseSync(t *testing.T) {
	syncer.sync(t, func(t *testing.T) {
		// Shows synctest.Test was called, since Wait panics if it's not in a "bubble"
		synctest.Wait()
	})
}

func SetupMockSync(t *testing.T) *MockSync {
	originalSync := syncer
	t.Cleanup(func() {
		syncer = originalSync
	})

	s := &MockSync{}
	syncer = s
	return s
}

type MockSync struct {
	Calls []T
}

func (s *MockSync) sync(t T, fn func(t *testing.T)) {
	s.Calls = append(s.Calls, t)
	fn(&testing.T{})
}
