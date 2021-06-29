package wrap_error

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"strings"

	"github.com/seatgeek/sgmods-go/analyzers"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const errorsPkgAlias = "errorsPkg"

var WrapErrorAnalyzer = &analysis.Analyzer{
	Name:     "wrap_error",
	Doc:      "check that new errors wrap context from existing errors in the call stack",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run: func(pass *analysis.Pass) (interface{}, error) {
		inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

		nodeFilter := []ast.Node{
			(*ast.IfStmt)(nil),
			(*ast.File)(nil),
			(*ast.GenDecl)(nil),
		}

		// was our ast.File modified during traversal? (eg, was an unwrapped error wrapped?)
		modified := false
		// is there already a package named "errors" imported in our ast.File?
		errorsNameConflict := false
		inspector.Nodes(nodeFilter, func(node ast.Node, push bool) bool {
			action := analyzers.Action(push)

			switch node := node.(type) {

			case *ast.File:
				switch action {
				case analyzers.Visit:
					for _, imp := range node.Imports {
						path := strings.Trim(imp.Path.Value, `"`)
						if strings.HasSuffix(path, "errors") && path != "github.com/pkg/errors" {
							errorsNameConflict = true
						}
					}
					modified = false
					return true
				case analyzers.Leave:
					if !modified {
						return true
					}
					importGenDecl, ok := findImportGenDecl(node)
					if !ok {
						panic("Expected to find an import generic declaration node")
					}
					pos := importGenDecl.Pos()
					end := importGenDecl.End()

					importSpec := &ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"github.com/pkg/errors"`,
						},
					}
					if errorsNameConflict {
						importSpec.Name = ast.NewIdent(errorsPkgAlias)
					}
					importGenDecl.Specs = append(importGenDecl.Specs, importSpec)

					pass.Report(analysis.Diagnostic{
						Pos:     pos,
						Message: `adding "github.com/pkg/errors" import`,
						SuggestedFixes: []analysis.SuggestedFix{
							{
								Message: `adding "github.com/pkg/errors" import`,
								TextEdits: []analysis.TextEdit{
									{
										Pos:     pos,
										End:     end,
										NewText: []byte(render(importGenDecl, pass.Fset)),
									},
								},
							},
						},
					})

				}

			case *ast.IfStmt:
				switch action {
				case analyzers.Visit:
					if !isErrNeqNull(node) {
						return true
					}

					// check for a return in our if statement body
					for _, stmt := range node.Body.List {
						returnStmt, ok := stmt.(*ast.ReturnStmt)
						if !ok {
							continue
						}

						// assume error is returned last
						if len(returnStmt.Results) == 0 {
							continue
						}
						lastResult := returnStmt.Results[len(returnStmt.Results)-1]

						resultValue, ok := lastIdent(lastResult)
						if !ok {
							continue
						}
						if !strings.HasPrefix(resultValue.Name, "Err") {
							continue
						}

						suggested := &ast.ReturnStmt{
							Return:  returnStmt.Return,
							Results: make([]ast.Expr, len(returnStmt.Results)),
						}
						copy(suggested.Results, returnStmt.Results)
						errorsSelector := "errors"
						if errorsNameConflict {
							errorsSelector = errorsPkgAlias
						}
						suggested.Results[len(suggested.Results)-1] = &ast.CallExpr{
							Fun: ast.NewIdent(errorsSelector + ".Wrap"),
							Args: []ast.Expr{
								lastResult,
								&ast.SelectorExpr{
									X:   ast.NewIdent("err"),
									Sel: ast.NewIdent("Error()"),
								},
							},
						}

						old := render(returnStmt, pass.Fset)
						new := render(suggested, pass.Fset)

						pass.Report(analysis.Diagnostic{
							Pos:     returnStmt.Pos(),
							Message: fmt.Sprintf("unwrapped error found '%s'", old),
							SuggestedFixes: []analysis.SuggestedFix{
								{
									Message: fmt.Sprintf("should replace '%s' with '%s'", old, new),
									TextEdits: []analysis.TextEdit{
										{
											Pos:     returnStmt.Pos(),
											End:     returnStmt.End(),
											NewText: []byte(new),
										},
									},
								},
							},
						})

						modified = true
					}
					return true
				}
			}
			return true
		})

		return nil, nil
	},
}

func isErrNeqNull(ifStatement *ast.IfStmt) bool {
	switch expr := ifStatement.Cond.(type) {
	case *ast.BinaryExpr:
		x, ok := expr.X.(*ast.Ident)
		if !ok {
			return false
		}
		y, ok := expr.Y.(*ast.Ident)
		if !ok {
			return false
		}

		return (x.Name == "err" &&
			expr.Op == token.NEQ &&
			y.Name == "nil")
	default:
		return false
	}
}

// MyStruct -> MyStruct
// helpers.MyStruct -> MyStruct
func lastIdent(expr ast.Expr) (*ast.Ident, bool) {
	switch expr := expr.(type) {
	case *ast.SelectorExpr:
		return lastIdent(expr.Sel)
	case *ast.Ident:
		return expr, true
	default:
		return nil, false
	}
}

func findImportGenDecl(node *ast.File) (*ast.GenDecl, bool) {
	for _, decl := range node.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.IMPORT {
			return genDecl, true
		}
	}

	return nil, false
}

func render(node interface{}, fset *token.FileSet) string {
	buf := bytes.Buffer{}
	printer.Fprint(&buf, fset, node)
	return buf.String()
}
