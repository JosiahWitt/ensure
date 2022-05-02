package externalconstraints

import (
	constraints2 "github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader/fixtures/generics/externalconstraints/constraints"
	"golang.org/x/exp/constraints"
)

type Thingable[T constraints.Ordered, V constraints2.Thing] interface {
	Identity(in T) (out T)
	Transform(in T) (out V)
}
