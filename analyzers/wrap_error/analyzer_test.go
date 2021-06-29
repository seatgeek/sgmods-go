package wrap_error

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestWrapErrorAnalyzer(t *testing.T) {
	analysistest.RunWithSuggestedFixes(t, analysistest.TestData(), WrapErrorAnalyzer)
}
