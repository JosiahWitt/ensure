//go:build !go1.18
// +build !go1.18

package ifacereader

import "go/types"

// TODO: Remove once Go 1.18 is the lowest supported version.
func (r *internalPackageReader) parseTypeParams(namedType *types.Named) []*TypeParam {
	// Prior to Go 1.18, type params weren't supported.
	return nil
}
