//go:build go1.18
// +build go1.18

package ifacereader

import (
	"go/types"
)

func (r *internalPackageReader) parseTypeParams(namedType *types.Named) []*TypeParam {
	typeParamList := namedType.TypeParams()
	typeParamCount := typeParamList.Len()

	if typeParamCount == 0 {
		return nil
	}

	typeParams := make([]*TypeParam, 0, typeParamCount)
	for i := 0; i < typeParamCount; i++ {
		typeParam := typeParamList.At(i)

		typeParams = append(typeParams, &TypeParam{
			Name: typeParam.Obj().Name(),
			Type: types.TypeString(typeParam.Constraint(), func(p *types.Package) string {
				return r.pkgNameGen.GeneratePackageName(r.pkg, p)
			}),
		})
	}

	return typeParams
}
