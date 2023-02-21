// Package reflectensure provides a helper for identifying ensure types via reflection.
// It is used to avoid import cycles.
package reflectensure

import "reflect"

const (
	ensuringPath = "github.com/JosiahWitt/ensure/ensuring"
	ensuringE    = "E"
)

// IsEnsuringE returns true only when [ensuring.E] or any of its aliases is provided.
func IsEnsuringE(t reflect.Type) bool {
	return t.PkgPath() == ensuringPath && t.Name() == ensuringE
}
