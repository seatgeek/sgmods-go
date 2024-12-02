package util

import "go/types"

// from: https://cs.opensource.google/go/x/tools/+/master:go/analysis/passes/internal/analysisutil/util.go;drc=e7bd2274d184f7579a8c7a1a12d8ad0351aeee8a;l=109
func Imports(pkg *types.Package, path string) bool {
	for _, imp := range pkg.Imports() {
		if imp.Path() == path {
			return true
		}
	}
	return false
}
