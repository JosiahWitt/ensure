package singletype

type Identifier[T any] interface {
	Identity(in T) (out T)
}
