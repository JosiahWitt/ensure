package namedoutputs

type NamedOutputs interface {
	NamedOut() (a string, b error)
}
