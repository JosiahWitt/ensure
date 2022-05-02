package complexconstraints

import (
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader/fixtures/generics/complexconstraints/externaltype"
	"golang.org/x/exp/constraints"
)

type Thing[T, V any] struct {
	A T
	B V
}

type Constraint interface {
	constraints.Ordered | ~map[string]string
}

type Thingable[
	T Constraint,
	V interface{ ~string },
	Composite *Thing[T, V],
	Unused constraints.Complex, // Unused, but should pull in the constraints package
] interface {
	Crazyness(in1 *Thing[T, externaltype.MyType], in2 Composite) T // externaltype package should be pulled in
}
