package multipletypes

type Thingable[T any, V any] interface {
	Identity(in T) (out T)
	Transform(in T) (out V)
}
