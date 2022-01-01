package nodepsiface

type NoIO interface {
	Method1()
}

type MultipleMethods interface {
	Method1(a string) string
	Method2(a string, b float64) (string, error)
}
