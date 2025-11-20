// Package mockgen generates mocks for the provided package interfaces.
package mockgen

import (
	"bytes"
	"text/template"

	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ensurefile"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/uniqpkg"
)

// Generator generates mocks for the provided package interfaces.
type Generator interface {
	GenerateMocks(pkgs []*ifacereader.Package, imports *uniqpkg.UniquePackagePaths, config *ensurefile.MockConfig) ([]*PackageMock, error)
}

// MockGen generates mocks for the provided package interfaces.
type MockGen struct {
	tmpl *template.Template
}

var _ Generator = &MockGen{}

// PackageMock contains the generated mocks for each package.
type PackageMock struct {
	Package      *ifacereader.Package
	FileContents string
}

// New creates a MockGen instance.
func New() (*MockGen, error) {
	tmpl, err := template.New("").Funcs(templateFuncs).Parse(packageTemplate)
	if err != nil {
		return nil, err // Shouldn't be possible, unless there's a syntax error in the internal template
	}

	return &MockGen{
		tmpl: tmpl,
	}, nil
}

// GenerateMocks generates mocks for the provided packages, using their respective imports.
func (g *MockGen) GenerateMocks(pkgs []*ifacereader.Package, imports *uniqpkg.UniquePackagePaths, config *ensurefile.MockConfig) ([]*PackageMock, error) {
	mocks := make([]*PackageMock, 0, len(pkgs))

	for _, pkg := range pkgs {
		mock, err := g.generateMock(pkg, imports.ForPackage(pkg.Path), config)
		if err != nil {
			return nil, err
		}

		mocks = append(mocks, mock)
	}

	return mocks, nil
}

func (g *MockGen) generateMock(pkg *ifacereader.Package, importsPkg *uniqpkg.Package, config *ensurefile.MockConfig) (*PackageMock, error) {
	reflectImport := importsPkg.AddImport("reflect", "reflect")
	goMockImport := importsPkg.AddImport("go.uber.org/mock/gomock", "gomock")

	var prettyImport *uniqpkg.ImportDetails
	if !config.DisableEnhancedMatcherFailures {
		prettyImport = importsPkg.AddImport("github.com/kr/pretty", "pretty")
	}

	params := &templateParams{
		Package: pkg,
		Imports: importsPkg.Imports(),

		ReflectPackageName: reflectImport.Name,
		GoMockPackageName:  goMockImport.Name,

		EnableEnhancedMatcherFailures: !config.DisableEnhancedMatcherFailures,
	}

	if params.EnableEnhancedMatcherFailures {
		params.PrettyPackageName = prettyImport.Name
	}

	var writer bytes.Buffer
	if err := g.tmpl.Execute(&writer, params); err != nil {
		return nil, err // Shouldn't be possible, since the parameters are controlled within this package
	}

	return &PackageMock{
		Package:      pkg,
		FileContents: writer.String(),
	}, nil
}
