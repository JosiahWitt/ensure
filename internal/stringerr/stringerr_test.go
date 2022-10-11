package stringerr_test

import (
	"errors"
	"testing"

	"github.com/JosiahWitt/ensure"
	"github.com/JosiahWitt/ensure/ensurepkg"
	"github.com/JosiahWitt/ensure/internal/stringerr"
)

func TestNewf(t *testing.T) {
	ensure := ensure.New(t)

	err := stringerr.Newf("some %s string", "formatted")
	ensure(err.Error()).Equals("some formatted string")
}

func TestNewGroup(t *testing.T) {
	ensure := ensure.New(t)

	err := stringerr.NewGroup("my prefix", []error{errors.New("first"), errors.New("second")})
	ensure(err.Error()).Equals("my prefix:\n - first\n - second")
}

func TestIs(t *testing.T) {
	ensure := ensure.New(t)

	ensure.Run("when equal", func(ensure ensurepkg.Ensure) {
		err := stringerr.Newf("some %s string", "formatted")
		ensure(errors.Is(err, errors.New("some formatted string"))).IsTrue()
	})

	ensure.Run("when not equal", func(ensure ensurepkg.Ensure) {
		err := stringerr.Newf("some %s string", "formatted")
		ensure(errors.Is(err, errors.New("some formatted thing"))).IsFalse()
	})

	ensure.Run("when comparing against nil", func(ensure ensurepkg.Ensure) {
		err := stringerr.Newf("some %s string", "formatted")
		ensure(errors.Is(err, nil)).IsFalse()
	})
}
