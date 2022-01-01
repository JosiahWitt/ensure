package example2

import "reflect"

var PackagePath = reflect.TypeOf(Float64(0)).PkgPath()

type Float64 float64

type User struct {
	ID   string
	Name string
}
