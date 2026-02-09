package example1

import "reflect"

var PackagePath = reflect.TypeFor[String]().PkgPath()

type String string

type Message struct {
	ID   string
	Body string
}
