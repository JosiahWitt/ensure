package reflectensure_test

import (
	"reflect"
	"testing"

	"github.com/JosiahWitt/ensure"
	"github.com/JosiahWitt/ensure/ensurepkg" //lint:ignore SA1019 To ensure compatibility
	"github.com/JosiahWitt/ensure/ensuring"
	"github.com/JosiahWitt/ensure/internal/reflectensure"
)

func TestIsEnsuringE(t *testing.T) {
	ensure := ensure.New(t)

	ensure.Run("when provided ensuring.E", func(ensure ensuring.E) {
		t := reflect.TypeFor[ensuring.E]()
		ensure(reflectensure.IsEnsuringE(t)).IsTrue()
	})

	ensure.Run("when provided pointer to ensuring.E", func(ensure ensuring.E) {
		t := reflect.TypeFor[*ensuring.E]()
		ensure(reflectensure.IsEnsuringE(t)).IsFalse()
	})

	ensure.Run("when provided ensurepkg.Ensure", func(ensure ensuring.E) {
		t := reflect.TypeFor[ensurepkg.Ensure]() //lint:ignore SA1019 To ensure compatibility
		ensure(reflectensure.IsEnsuringE(t)).IsTrue()
	})

	ensure.Run("when provided another type implementing ensuring.E", func(ensure ensuring.E) {
		type E ensuring.E
		t := reflect.TypeFor[E]()
		ensure(reflectensure.IsEnsuringE(t)).IsFalse()
	})

	ensure.Run("when provided another type named E", func(ensure ensuring.E) {
		type E func(any) *ensuring.Chain
		t := reflect.TypeFor[E]()
		ensure(reflectensure.IsEnsuringE(t)).IsFalse()
	})

	ensure.Run("when provided another type in ensuring", func(ensure ensuring.E) {
		t := reflect.TypeFor[ensuring.Chain]()
		ensure(reflectensure.IsEnsuringE(t)).IsFalse()
	})
}
