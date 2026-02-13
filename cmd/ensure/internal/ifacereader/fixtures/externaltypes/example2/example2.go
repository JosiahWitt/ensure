package example2

import "reflect"

var PackagePath = reflect.TypeFor[Float64]().PkgPath()

type Float64 float64

type User struct {
	ID   string
	Name string
}
