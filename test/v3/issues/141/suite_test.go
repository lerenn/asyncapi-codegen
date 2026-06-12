//go:generate go run ../../../../cmd/asyncapi-codegen -p issue141 -i ./asyncapi.yaml -o ./asyncapi.gen.go

package issue141

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestSuite(t *testing.T) {
	suite.Run(t, NewSuite())
}

type Suite struct {
	suite.Suite
}

func NewSuite() *Suite {
	return &Suite{}
}

// TestExportedDeclarationsHaveDoc ensures every exported declaration in the
// generated code has a doc comment, so the generated API is fully documented
// (issue #141).
func (suite *Suite) TestExportedDeclarationsHaveDoc() {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "asyncapi.gen.go", nil, parser.ParseComments)
	suite.Require().NoError(err)

	for _, decl := range f.Decls {
		suite.checkDecl(decl)
	}
}

func (suite *Suite) checkDecl(decl ast.Decl) {
	switch d := decl.(type) {
	case *ast.FuncDecl:
		if d.Name.IsExported() {
			suite.Truef(d.Doc != nil, "exported function %q has no doc comment", d.Name.Name)
		}
	case *ast.GenDecl:
		// A doc comment on the group documents every spec in it.
		if d.Doc != nil {
			return
		}
		for _, spec := range d.Specs {
			suite.checkSpec(spec)
		}
	}
}

func (suite *Suite) checkSpec(spec ast.Spec) {
	switch s := spec.(type) {
	case *ast.TypeSpec:
		if s.Name.IsExported() {
			suite.Truef(s.Doc != nil, "exported type %q has no doc comment", s.Name.Name)
		}
	case *ast.ValueSpec:
		for _, name := range s.Names {
			if name.IsExported() {
				suite.Truef(s.Doc != nil, "exported value %q has no doc comment", name.Name)
			}
		}
	}
}
