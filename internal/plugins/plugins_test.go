package plugins_test

import (
	"reflect"
	"testing"

	"github.com/JosiahWitt/ensure"
	"github.com/JosiahWitt/ensure/internal/plugins"
)

func TestNoopAfterEntry(t *testing.T) {
	ensure := ensure.New(t)

	err := plugins.NoopAfterEntry{}.AfterEntry(nil, reflect.Value{}, 0)
	ensure(err).IsNotError()
}
