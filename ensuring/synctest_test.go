//go:build go1.25

package ensuring_test

import (
	"testing"

	"github.com/JosiahWitt/ensure/ensuring"
)

func TestERunSync(t *testing.T) {
	runConfig{
		isSync: true,
		prepare: func(ensure ensuring.E) func(string, func(ensuring.E)) {
			return ensure.RunSync
		},
	}.test(t)
}

func TestERunTableByIndexSync(t *testing.T) {
	runTableConfig{
		isSync: true,
		prepare: func(ensure ensuring.E) func(table interface{}, fn func(ensure ensuring.E, i int)) {
			return ensure.RunTableByIndexSync
		},
	}.test(t)
}
