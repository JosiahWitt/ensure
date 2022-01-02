// Package ifacereader reads the specified interfaces, returning their listed methods.
package ifacereader

import (
	"go/types"

	"github.com/JosiahWitt/erk"
	"github.com/JosiahWitt/erk/erg"
	"golang.org/x/tools/go/packages"
)

type (
	ErkInvalidInput struct{ erk.DefaultKind }
	ErkInternal     struct{ erk.DefaultKind }
)

var (
	ErrNoInterfaces     = erk.New(ErkInvalidInput{}, "No interfaces provided for path '{{.path}}'. Please provide which interfaces to mock.")
	ErrDuplicatePath    = erk.New(ErkInvalidInput{}, "Duplicate entry for path '{{.path}}'. Please combine both entries and list multiple interfaces instead.")
	ErrLoadingPackages  = erk.New(ErkInvalidInput{}, "Unable to load all packages: {{.err}}")
	ErrReadingPackage   = erk.New(ErkInvalidInput{}, "Error reading package '{{.path}}'")
	ErrInvalidInterface = erk.New(ErkInvalidInput{}, "Interface '{{.interface}}' not found in package: {{.package}}")
	ErrNotInterface     = erk.New(ErkInvalidInput{}, "Type '{{.interface}}' is not an interface in package '{{.package}}', it's a '{{.type}}'")

	ErrPathMismatch           = erk.New(ErkInternal{}, "Unexpected error; please open an issue on GitHub: Could not find package details for path: {{.path}}")
	ErrLeftoverPackageDetails = erk.New(ErkInternal{}, "Unexpected error; please open an issue on GitHub: Unexpected leftover package details")
	ErrUnsupportedType        = erk.New(ErkInternal{}, "Unexpected error; please open an issue on GitHub: Unsupported type '{{type .rawType}}': {{.typeString}}")
)

// Readable reads the specified interfaces, returning their listed methods.
type Readable interface {
	ReadPackages(pkgDetails []*PackageDetails, pkgNameGen PackageNameGenerator) ([]*Package, error)
}

// InterfaceReader reads the specified interfaces, returning their listed methods.
type InterfaceReader struct{}

var _ Readable = &InterfaceReader{}

// PackageNameGenerator allows generating a package name for a type given the current package and the imported package.
type PackageNameGenerator interface {
	GeneratePackageName(scopePackage *packages.Package, importedPackage *types.Package) string
}

// PackageDetails provides the package path and interfaces to the ReadPackages method.
type PackageDetails struct {
	Path       string
	Interfaces []string
}

// Package includes the details of parsing the interfaces in the package.
type Package struct {
	Path       string
	Interfaces []*Interface
}

// Interface includes the details of parsing the interface in the package.
type Interface struct {
	Name    string
	Methods []*Method
}

// Method includes the details of a single method inside of an interface.
type Method struct {
	Name    string
	Inputs  []*Tuple
	Outputs []*Tuple
}

// Tuple includes the details of a single input or output parameter in a method signature.
type Tuple struct {
	VariableName string
	PackagePaths []string
	Type         string
}

type internalPackageReader struct {
	pkg        *packages.Package
	pkgNameGen PackageNameGenerator
}

// ReadPackages reads all the packages within the specified package and interface combinations.
func (r *InterfaceReader) ReadPackages(pkgDetails []*PackageDetails, pkgNameGen PackageNameGenerator) ([]*Package, error) {
	pkgDetailsByPath := make(map[string]*PackageDetails, len(pkgDetails))
	pkgPaths := make([]string, 0, len(pkgDetails))

	for _, pkgDetail := range pkgDetails {
		if len(pkgDetail.Interfaces) == 0 {
			return nil, erk.WithParams(ErrNoInterfaces, erk.Params{"path": pkgDetail.Path})
		}

		if _, ok := pkgDetailsByPath[pkgDetail.Path]; ok {
			return nil, erk.WithParams(ErrDuplicatePath, erk.Params{"path": pkgDetail.Path})
		}

		pkgDetailsByPath[pkgDetail.Path] = pkgDetail
		pkgPaths = append(pkgPaths, pkgDetail.Path)
	}

	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedImports | packages.NeedTypes,
	}

	rawPkgs, err := packages.Load(cfg, pkgPaths...)
	if err != nil {
		return nil, erk.WrapAs(ErrLoadingPackages, err)
	}

	pkgs := make([]*Package, 0, len(rawPkgs))
	for _, pkg := range rawPkgs {
		if len(pkg.Errors) > 0 {
			err := erg.NewAs(erk.WithParams(ErrReadingPackage, erk.Params{"path": pkg.PkgPath}))
			for _, pkgErr := range pkg.Errors {
				err = erg.Append(err, pkgErr)
			}

			return nil, err
		}

		pkgDetail, ok := pkgDetailsByPath[pkg.PkgPath]
		if !ok {
			// Not sure if this is possible
			return nil, erk.WithParams(ErrPathMismatch, erk.Params{"path": pkgDetail.Path})
		}

		pkgReader := &internalPackageReader{
			pkg:        pkg,
			pkgNameGen: pkgNameGen,
		}

		builtPkg, err := pkgReader.buildPackage(pkgDetail, pkg)
		if err != nil {
			return nil, err
		}

		pkgs = append(pkgs, builtPkg)
		delete(pkgDetailsByPath, pkg.PkgPath) // Delete so we can check if all packages were loaded
	}

	if len(pkgDetailsByPath) != 0 {
		// Don't think this is possible
		return nil, erk.WithParams(ErrLeftoverPackageDetails, erk.Params{"leftoverDetails": pkgDetailsByPath})
	}

	return pkgs, nil
}

func (r *internalPackageReader) buildPackage(pkgDetail *PackageDetails, pkg *packages.Package) (*Package, error) {
	ifaces := make([]*Interface, 0, len(pkgDetail.Interfaces))

	for _, ifaceName := range pkgDetail.Interfaces {
		rawIface := pkg.Types.Scope().Lookup(ifaceName)
		if rawIface == nil {
			return nil, erk.WithParams(ErrInvalidInterface, erk.Params{
				"interface": ifaceName,
				"package":   pkgDetail.Path,
			})
		}

		iface, ok := rawIface.Type().Underlying().(*types.Interface)
		if !ok {
			return nil, erk.WithParams(ErrNotInterface, erk.Params{
				"interface": ifaceName,
				"package":   pkgDetail.Path,
				"type":      rawIface.String(),
			})
		}

		builtIface, err := r.buildIface(ifaceName, iface)
		if err != nil {
			return nil, err
		}

		ifaces = append(ifaces, builtIface)
	}

	return &Package{
		Path:       pkgDetail.Path,
		Interfaces: ifaces,
	}, nil
}

func (r *internalPackageReader) buildIface(ifaceName string, iface *types.Interface) (*Interface, error) {
	methods := make([]*Method, 0, iface.NumMethods())

	for i := 0; i < iface.NumMethods(); i++ {
		builtMethod, err := r.buildMethod(iface.Method(i))
		if err != nil {
			return nil, err
		}

		methods = append(methods, builtMethod)
	}

	return &Interface{
		Name:    ifaceName,
		Methods: methods,
	}, nil
}

func (r *internalPackageReader) buildMethod(method *types.Func) (*Method, error) {
	signature := method.Type().Underlying().(*types.Signature)

	inputs := make([]*Tuple, 0, signature.Params().Len())
	for i := 0; i < signature.Params().Len(); i++ {
		param := signature.Params().At(i)

		builtInput, err := r.buildTuple(param.Name(), param.Type())
		if err != nil {
			return nil, err
		}

		inputs = append(inputs, builtInput)
	}

	outputs := make([]*Tuple, 0, signature.Results().Len())
	for i := 0; i < signature.Results().Len(); i++ {
		result := signature.Results().At(i)

		builtOutput, err := r.buildTuple(result.Name(), result.Type())
		if err != nil {
			return nil, err
		}

		outputs = append(outputs, builtOutput)
	}

	return &Method{
		Name:    method.Name(),
		Inputs:  inputs,
		Outputs: outputs,
	}, nil
}

func (r *internalPackageReader) buildTuple(variableName string, rawType types.Type) (*Tuple, error) {
	pkgPaths, err := r.extractPackagePaths(rawType)
	if err != nil {
		return nil, err
	}

	tuple := &Tuple{
		VariableName: variableName,
		PackagePaths: pkgPaths,

		Type: types.TypeString(rawType, func(p *types.Package) string {
			return r.pkgNameGen.GeneratePackageName(r.pkg, p)
		}),
	}

	return tuple, nil
}

func (r *internalPackageReader) extractPackagePaths(rawType types.Type) ([]string, error) {
	switch t := rawType.(type) {
	case *types.Named:
		if pkg := t.Obj().Pkg(); pkg != nil {
			return []string{pkg.Path()}, nil
		}

		return []string{}, nil

	case *types.Basic:
		return []string{}, nil

	case *types.Slice:
		return r.extractPackagePaths(t.Elem())

	case *types.Array:
		return r.extractPackagePaths(t.Elem())

	case *types.Pointer:
		return r.extractPackagePaths(t.Elem())

	case *types.Chan:
		return r.extractPackagePaths(t.Elem())

	case *types.Map:
		return r.extractMapPackagePaths(t)

	case *types.Signature:
		return r.extractSignaturePackagePaths(t)

	case *types.Interface:
		return r.extractInterfacePackagePaths(t)

	case *types.Struct:
		return r.extractStructPackagePaths(t)

	default:
		// Shouldn't be possible, unless some types are missing
		return nil, erk.WithParams(ErrUnsupportedType, erk.Params{
			"rawType":    rawType,
			"typeString": rawType.String(),
		})
	}
}

func (r *internalPackageReader) extractMapPackagePaths(t *types.Map) ([]string, error) {
	keyPaths, err := r.extractPackagePaths(t.Key())
	if err != nil {
		return nil, err
	}

	elemPaths, err := r.extractPackagePaths(t.Elem())
	if err != nil {
		return nil, err
	}

	return uniqueStrings(append(keyPaths, elemPaths...)), nil
}

func (r *internalPackageReader) extractSignaturePackagePaths(t *types.Signature) ([]string, error) {
	paths := make([]string, 0, t.Params().Len()+t.Results().Len())

	for i := 0; i < t.Params().Len(); i++ {
		paramPaths, err := r.extractPackagePaths(t.Params().At(i).Type())
		if err != nil {
			return nil, err
		}

		paths = append(paths, paramPaths...)
	}

	for i := 0; i < t.Results().Len(); i++ {
		resultPaths, err := r.extractPackagePaths(t.Results().At(i).Type())
		if err != nil {
			return nil, err
		}

		paths = append(paths, resultPaths...)
	}

	return uniqueStrings(paths), nil
}

func (r *internalPackageReader) extractInterfacePackagePaths(t *types.Interface) ([]string, error) {
	builtIface, err := r.buildIface("<unused>", t)
	if err != nil {
		return nil, err
	}

	var paths []string
	for _, method := range builtIface.Methods {
		for _, input := range method.Inputs {
			paths = append(paths, input.PackagePaths...)
		}

		for _, output := range method.Outputs {
			paths = append(paths, output.PackagePaths...)
		}
	}

	return uniqueStrings(paths), nil
}

func (r *internalPackageReader) extractStructPackagePaths(t *types.Struct) ([]string, error) {
	paths := make([]string, 0, t.NumFields())

	for i := 0; i < t.NumFields(); i++ {
		fieldPaths, err := r.extractPackagePaths(t.Field(i).Type())
		if err != nil {
			return nil, err
		}

		paths = append(paths, fieldPaths...)
	}

	return uniqueStrings(paths), nil
}

func uniqueStrings(strs []string) []string {
	uniqueStrs := make([]string, 0, len(strs))
	exists := make(map[string]bool, len(strs))

	for _, str := range strs {
		if !exists[str] {
			exists[str] = true
			uniqueStrs = append(uniqueStrs, str)
		}
	}

	return uniqueStrs
}

// InterfaceNames extracts each of the interface names for which to generate mocks.
func (pkg *Package) InterfaceNames() []string {
	names := make([]string, 0, len(pkg.Interfaces))

	for _, iface := range pkg.Interfaces {
		names = append(names, iface.Name)
	}

	return names
}
