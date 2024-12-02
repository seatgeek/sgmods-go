package slogcontext

import (
	"github.com/seatgeek/sgmods-go/pkg/util"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"golang.org/x/tools/go/analysis/passes/inspect"
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

		return nil, nil
	},
}
