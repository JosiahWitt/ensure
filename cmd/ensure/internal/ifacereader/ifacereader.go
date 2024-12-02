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

const unexpectedErrorPrefix = "Unexpected error; please open an issue on GitHub (https://github.com/JosiahWitt/ensure/issues)"

var (
	ErrNoInterfaces     = erk.New(ErkInvalidInput{}, "No interfaces provided for path '{{.path}}'. Please provide which interfaces to mock.")
	ErrDuplicatePath    = erk.New(ErkInvalidInput{}, "Duplicate entry for path '{{.path}}'. Please combine both entries and list multiple interfaces instead.")
	ErrLoadingPackages  = erk.New(ErkInvalidInput{}, "Unable to load all packages: {{.err}}")
	ErrReadingPackage   = erk.New(ErkInvalidInput{}, "Error reading package '{{.path}}'")
	ErrInvalidInterface = erk.New(ErkInvalidInput{}, "Interface '{{.interface}}' not found in package: {{.package}}")
	ErrNotInterface     = erk.New(ErkInvalidInput{}, "Type '{{.interface}}' is not an interface in package '{{.package}}', it's a '{{.type}}'")

	ErrPathMismatch           = erk.New(ErkInternal{}, unexpectedErrorPrefix+": Could not find package details for path: {{.path}}")
	ErrLeftoverPackageDetails = erk.New(ErkInternal{}, unexpectedErrorPrefix+": Unexpected leftover package details")
	ErrInterfaceTypeNotNamed  = erk.New(ErkInternal{}, unexpectedErrorPrefix+": interface type for '{{.interface}}' was not *types.Named, it was: {{type .type}}")
	ErrFuncUnderlyingType     = erk.New(ErkInternal{},
		unexpectedErrorPrefix+": *types.Func underlying type was not *types.Signature, it was: {{type .underlyingType}}",
	)
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
	Name       string
	Path       string
	Interfaces []*Interface
}

// Interface includes the details of parsing the interface in the package.
type Interface struct {
	Name       string
	Methods    []*Method
	TypeParams []*TypeParam
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
	Type         string
	Variadic     bool
}

// TypeParam contains details about a Go 1.18+ generic type parameter.
type TypeParam struct {
	Name string
	Type string
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
			return nil, buildPackageReadError(pkg)
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

		namedIface, ok := rawIface.Type().(*types.Named)
		if !ok {
			// Not sure if this is possible
			return nil, erk.WithParams(ErrInterfaceTypeNotNamed, erk.Params{
				"interface": ifaceName,
				"package":   pkgDetail.Path,
				"type":      rawIface.Type(),
			})
		}

		builtIface.TypeParams = r.parseTypeParams(namedIface)
		ifaces = append(ifaces, builtIface)
	}

	return &Package{
		Name:       pkg.Name,
		Path:       pkgDetail.Path,
		Interfaces: ifaces,
	}, nil
}

func (r *internalPackageReader) buildIface(ifaceName string, iface *types.Interface) (*Interface, error) {
	methods := make([]*Method, 0, iface.NumMethods())

	for i := range iface.NumMethods() {
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
	signature, ok := method.Type().Underlying().(*types.Signature)
	if !ok {
		return nil, erk.WithParams(ErrFuncUnderlyingType, erk.Params{"underlyingType": method.Type().Underlying()}) // Not sure if this is possible
	}

	inputs := make([]*Tuple, 0, signature.Params().Len())
	for i := range signature.Params().Len() {
		param := signature.Params().At(i)
		inputs = append(inputs, r.buildTuple(param.Name(), param.Type()))
	}

	outputs := make([]*Tuple, 0, signature.Results().Len())
	for i := range signature.Results().Len() {
		result := signature.Results().At(i)
		outputs = append(outputs, r.buildTuple(result.Name(), result.Type()))
	}

	// If the signature is variadic, mark the last input type variadic
	if signature.Variadic() {
		inputs[len(inputs)-1].Variadic = true
	}

	return &Method{
		Name:    method.Name(),
		Inputs:  inputs,
		Outputs: outputs,
	}, nil
}

func (r *internalPackageReader) buildTuple(variableName string, rawType types.Type) *Tuple {
	return &Tuple{
		VariableName: variableName,

		Type: types.TypeString(rawType, func(p *types.Package) string {
			return r.pkgNameGen.GeneratePackageName(r.pkg, p)
		}),
	}
}

func buildPackageReadError(pkg *packages.Package) error {
	err := erg.NewAs(erk.WithParams(ErrReadingPackage, erk.Params{"path": pkg.PkgPath}))

	for _, pkgErr := range pkg.Errors {
		err = erg.Append(err, pkgErr)
	}

	return err
}

// InterfaceNames extracts each of the interface names for which to generate mocks.
func (pkg *Package) InterfaceNames() []string {
	names := make([]string, 0, len(pkg.Interfaces))

	for _, iface := range pkg.Interfaces {
		names = append(names, iface.Name)
	}

	return names
}
