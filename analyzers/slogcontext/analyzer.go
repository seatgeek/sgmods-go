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

			pass.Report(analysis.Diagnostic{
				Pos:     node.Pos(),
				End:     node.End(),
				Message: "context not passed to slog call",
			})

			return false
		})

		return nil, nil
	},
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
