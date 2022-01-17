package example1

import "reflect"

var PackagePath = reflect.TypeOf(String("")).PkgPath()

type String string

type Message struct {
	ID   string
	Body string
}
