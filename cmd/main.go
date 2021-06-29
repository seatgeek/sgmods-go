package main

import (
	"github.com/seatgeek/sgmods-go/analyzers/wrap_error"

	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	multichecker.Main(wrap_error.WrapErrorAnalyzer)
}
