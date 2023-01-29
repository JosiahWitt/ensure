package stringerr_test

import (
	"errors"
	"strings"
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

func TestNewfIs(t *testing.T) {
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

func TestNewGroup(t *testing.T) {
	ensure := ensure.New(t)

	ensure.Run("with no nested errors", func(ensure ensurepkg.Ensure) {
		err := stringerr.NewGroup("my prefix", []error{errors.New("first"), errors.New("second")})

		ensure(err.Error()).Equals(w(`
my prefix:
 - first
 - second
`))
	})

	ensure.Run("with nested group", func(ensure ensurepkg.Ensure) {
		grpErr := stringerr.NewGroup("nested prefix", []error{errors.New("nested first"), errors.New("nested second")})
		err := stringerr.NewGroup("my prefix", []error{errors.New("first"), grpErr, errors.New("third")})

		ensure(err.Error()).Equals(w(`
my prefix:
 - first
 - nested prefix:
    - nested first
    - nested second
 - third
`))
	})

	ensure.Run("with double nested group", func(ensure ensurepkg.Ensure) {
		doubleGrpErr := stringerr.NewGroup("double nested prefix", []error{errors.New("double nested")})
		grpErr := stringerr.NewGroup("nested prefix", []error{errors.New("nested first"), doubleGrpErr, errors.New("nested third")})
		err := stringerr.NewGroup("my prefix", []error{errors.New("first"), grpErr, errors.New("third")})

		ensure(err.Error()).Equals(w(`
my prefix:
 - first
 - nested prefix:
    - nested first
    - double nested prefix:
       - double nested
    - nested third
 - third
`))
	})

	ensure.Run("with nested block", func(ensure ensurepkg.Ensure) {
		blockErr := stringerr.NewBlock("nested prefix", []error{errors.New("nested first"), errors.New("nested second")}, "nested footer")
		err := stringerr.NewGroup("my prefix", []error{errors.New("first"), blockErr, errors.New("third")})

		ensure(err.Error()).Equals(w(`
my prefix:
 - first
 - nested prefix:
    - nested first
    - nested second
   nested footer
 - third
`))
	})

	ensure.Run("with double nested block", func(ensure ensurepkg.Ensure) {
		doubleBlockErr := stringerr.NewBlock("double nested prefix", []error{errors.New("double nested")}, "double nested footer")
		blockErr := stringerr.NewBlock("nested prefix", []error{errors.New("nested first"), doubleBlockErr, errors.New("nested third")}, "nested footer")
		err := stringerr.NewGroup("my prefix", []error{errors.New("first"), blockErr, errors.New("third")})

		ensure(err.Error()).Equals(w(`
my prefix:
 - first
 - nested prefix:
    - nested first
    - double nested prefix:
       - double nested
      double nested footer
    - nested third
   nested footer
 - third
`))
	})
}

func TestNewGroupIs(t *testing.T) {
	ensure := ensure.New(t)

	ensure.Run("when equal", func(ensure ensurepkg.Ensure) {
		err := stringerr.NewGroup("my prefix", []error{errors.New("first"), errors.New("second")})
		err2 := stringerr.NewGroup("my prefix", []error{errors.New("first"), errors.New("second")})
		ensure(errors.Is(err, err2)).IsTrue()
	})

	ensure.Run("when prefix not equal", func(ensure ensurepkg.Ensure) {
		err := stringerr.NewGroup("my prefix", []error{errors.New("first"), errors.New("second")})
		err2 := stringerr.NewGroup("my prefix!", []error{errors.New("first"), errors.New("second")})
		ensure(errors.Is(err, err2)).IsFalse()
	})

	ensure.Run("when grouped error not equal", func(ensure ensurepkg.Ensure) {
		err := stringerr.NewGroup("my prefix", []error{errors.New("first"), errors.New("second")})
		err2 := stringerr.NewGroup("my prefix", []error{errors.New("first one"), errors.New("second")})
		ensure(errors.Is(err, err2)).IsFalse()
	})

	ensure.Run("when comparing against nil", func(ensure ensurepkg.Ensure) {
		err := stringerr.NewGroup("my prefix", []error{errors.New("first"), errors.New("second")})
		ensure(errors.Is(err, nil)).IsFalse()
	})
}

func TestNewBlock(t *testing.T) {
	ensure := ensure.New(t)

	ensure.Run("with no nested errors", func(ensure ensurepkg.Ensure) {
		err := stringerr.NewBlock("my prefix", []error{errors.New("first"), errors.New("second")}, "my footer")

		ensure(err.Error()).Equals(w(`
my prefix:
 - first
 - second
my footer
`))
	})

	ensure.Run("with nested group", func(ensure ensurepkg.Ensure) {
		grpErr := stringerr.NewGroup("nested prefix", []error{errors.New("nested first"), errors.New("nested second")})
		err := stringerr.NewBlock("my prefix", []error{errors.New("first"), grpErr, errors.New("third")}, "my footer")

		ensure(err.Error()).Equals(w(`
my prefix:
 - first
 - nested prefix:
    - nested first
    - nested second
 - third
my footer
`))
	})

	ensure.Run("with double nested group", func(ensure ensurepkg.Ensure) {
		doubleGrpErr := stringerr.NewGroup("double nested prefix", []error{errors.New("double nested")})
		grpErr := stringerr.NewGroup("nested prefix", []error{errors.New("nested first"), doubleGrpErr, errors.New("nested third")})
		err := stringerr.NewBlock("my prefix", []error{errors.New("first"), grpErr, errors.New("third")}, "my footer")

		ensure(err.Error()).Equals(w(`
my prefix:
 - first
 - nested prefix:
    - nested first
    - double nested prefix:
       - double nested
    - nested third
 - third
my footer
`))
	})

	ensure.Run("with nested block", func(ensure ensurepkg.Ensure) {
		blockErr := stringerr.NewBlock("nested prefix", []error{errors.New("nested first"), errors.New("nested second")}, "nested footer")
		err := stringerr.NewBlock("my prefix", []error{errors.New("first"), blockErr, errors.New("third")}, "my footer")

		ensure(err.Error()).Equals(w(`
my prefix:
 - first
 - nested prefix:
    - nested first
    - nested second
   nested footer
 - third
my footer
`))
	})

	ensure.Run("with double nested block", func(ensure ensurepkg.Ensure) {
		doubleBlockErr := stringerr.NewBlock("double nested prefix", []error{errors.New("double nested")}, "double nested footer")
		blockErr := stringerr.NewBlock("nested prefix", []error{errors.New("nested first"), doubleBlockErr, errors.New("nested third")}, "nested footer")
		err := stringerr.NewBlock("my prefix", []error{errors.New("first"), blockErr, errors.New("third")}, "my footer")

		ensure(err.Error()).Equals(w(`
my prefix:
 - first
 - nested prefix:
    - nested first
    - double nested prefix:
       - double nested
      double nested footer
    - nested third
   nested footer
 - third
my footer
`))
	})
}

func TestNewBlockIs(t *testing.T) {
	ensure := ensure.New(t)

	ensure.Run("when equal", func(ensure ensurepkg.Ensure) {
		err := stringerr.NewBlock("my prefix", []error{errors.New("first"), errors.New("second")}, "my footer")
		err2 := stringerr.NewBlock("my prefix", []error{errors.New("first"), errors.New("second")}, "my footer")
		ensure(errors.Is(err, err2)).IsTrue()
	})

	ensure.Run("when prefix not equal", func(ensure ensurepkg.Ensure) {
		err := stringerr.NewBlock("my prefix", []error{errors.New("first"), errors.New("second")}, "my footer")
		err2 := stringerr.NewBlock("my prefix!", []error{errors.New("first"), errors.New("second")}, "my footer")
		ensure(errors.Is(err, err2)).IsFalse()
	})

	ensure.Run("when grouped error not equal", func(ensure ensurepkg.Ensure) {
		err := stringerr.NewBlock("my prefix", []error{errors.New("first"), errors.New("second")}, "my footer")
		err2 := stringerr.NewBlock("my prefix", []error{errors.New("first one"), errors.New("second")}, "my footer")
		ensure(errors.Is(err, err2)).IsFalse()
	})

	ensure.Run("when footer not equal", func(ensure ensurepkg.Ensure) {
		err := stringerr.NewBlock("my prefix", []error{errors.New("first"), errors.New("second")}, "my footer")
		err2 := stringerr.NewBlock("my prefix", []error{errors.New("first"), errors.New("second")}, "my footer!")
		ensure(errors.Is(err, err2)).IsFalse()
	})

	ensure.Run("when comparing against nil", func(ensure ensurepkg.Ensure) {
		err := stringerr.NewBlock("my prefix", []error{errors.New("first"), errors.New("second")}, "my footer")
		ensure(errors.Is(err, nil)).IsFalse()
	})
}

func w(s string) string {
	return strings.TrimSpace(s)
}
