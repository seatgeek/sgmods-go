package slogcontext

import (
	"fmt"
	"go/ast"

	"github.com/seatgeek/sgmods-go/pkg/util"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const slogPackage = "log/slog"

var SlogContextAnalyzer = &analysis.Analyzer{
	Name:     "slogcontext",
	Doc:      "check that context is passed to all slog calls",
	Requires: []*analysis.Analyzer{inspect.Analyzer, ctrlflow.Analyzer},
	Run: func(pass *analysis.Pass) (interface{}, error) {
		if !util.Imports(pass.Pkg, slogPackage) {
			return nil, nil
		}

		nodeFilter := []ast.Node{
			(*ast.CallExpr)(nil),
		}

		inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
		inspector.WithStack(nodeFilter, func(node ast.Node, push bool, stack []ast.Node) bool {
			callExpr, ok := node.(*ast.CallExpr)
			if !ok {
				panic(fmt.Sprintf("unexpected node type %T", node))
			}

			slogCall, ok := determineSlogCall(callExpr)
			if !ok {
				return false
			}

			switch slogCall {
			case "Debug", "Info", "Warn", "Error":
				break
			default:
				return false
			}

			availableCtx := availableContext(stack)
			newCallExpr := *callExpr
			newCallExpr.Fun.(*ast.SelectorExpr).Sel.Name += "Context"
			if availableCtx == "" {
				newCallExpr.Args = append([]ast.Expr{ast.NewIdent("context.TODO()")}, newCallExpr.Args...)
			} else {
				newCallExpr.Args = append([]ast.Expr{ast.NewIdent(availableCtx)}, newCallExpr.Args...)
			}
			newText := util.Render(&newCallExpr, pass.Fset)

			pass.Report(analysis.Diagnostic{
				Pos:     node.Pos(),
				End:     node.End(),
				Message: "context not passed to slog call",
				SuggestedFixes: []analysis.SuggestedFix{
					{
						Message: "Add context to slog call",
						TextEdits: []analysis.TextEdit{
							{
								Pos:     node.Pos(),
								End:     node.End(),
								NewText: []byte(newText),
							},
						},
					},
				},
			})

			return false
		})

		return nil, nil
	},
}

func availableContext(stack []ast.Node) string {
	fn := containingFunc(stack)
	if fn == nil || fn.Type.Params.NumFields() == 0 {
		return ""
	}

	firstArg := fn.Type.Params.List[0]
	if selectorMatches(firstArg.Type, "context", "Context") {
		return firstArg.Names[0].Name
	}

	return ""
}

func containingFunc(stack []ast.Node) *ast.FuncDecl {
	for i := 0; i < len(stack); i++ {
		if fn, ok := stack[i].(*ast.FuncDecl); ok {
			return fn
		}
	}

	return nil
}

// returns the slog call name, false if not an slog call
func determineSlogCall(callExpr *ast.CallExpr) (string, bool) {
	selector, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return "", false
	}

	x, ok := selector.X.(*ast.Ident)
	if !ok {
		return "", false
	}

	if x.Name != "slog" {
		return "", false
	}

	return selector.Sel.Name, true
}

func selectorMatches(node ast.Node, x string, y string) bool {
	selector, ok := node.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	xIdent, ok := selector.X.(*ast.Ident)
	if !ok {
		return false
	}

	if xIdent.Name != x {
		return false
	}

	if selector.Sel.Name != y {
		return false
	}

	return true
}
