//go:build go1.25

package integration_test_suite_test

import (
	"testing"
	"testing/synctest"
	"time"

	"github.com/JosiahWitt/ensure"
	"github.com/JosiahWitt/ensure/ensuring"
)

func TestRunSync(t *testing.T) {
	t.Run("callback is executed within synctest.Test", func(t *testing.T) {
		ensure := ensure.New(t)

		start := time.Now()
		ensure.RunSync("hello", func(ensure ensuring.E) {
			// Shows synctest.Test was called, since Wait panics if it's not in a "bubble"
			synctest.Wait()

			time.Sleep(15 * time.Second)
		})

		// Shows time.Sleep within RunSync doesn't wait actual clock time
		assertShorterThan(t, time.Since(start), 10*time.Second)
	})

	sharedEnsureRunTests(t, func(ensure ensuring.E) func(string, func(ensuring.E)) {
		return ensure.RunSync
	})
}

func TestRunTableByIndexSync(t *testing.T) {
	t.Run("callback is executed within synctest.Test", func(t *testing.T) {
		ensure := ensure.New(t)

		table := []struct {
			Name string
		}{
			{
				Name: "Hello",
			},
		}

		start := time.Now()
		ensure.RunTableByIndexSync(table, func(ensure ensuring.E, i int) {
			// Shows synctest.Test was called, since Wait panics if it's not in a "bubble"
			synctest.Wait()

			time.Sleep(15 * time.Second)
		})

		// Shows time.Sleep within RunSync doesn't wait actual clock time
		assertShorterThan(t, time.Since(start), 10*time.Second)
	})

	sharedEnsureRunTableByIndexTests(t, func(ensure ensuring.E) func(table any, fn func(ensure ensuring.E, i int)) {
		return ensure.RunTableByIndexSync
	})
}

func assertShorterThan(t *testing.T, actual, maximum time.Duration) {
	if actual > maximum {
		t.Fatalf("%v > %v", actual, maximum)
	}
}
