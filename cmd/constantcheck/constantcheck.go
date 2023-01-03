package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"constantcheck/pkg/analyzer"
)

func main() {
	singlechecker.Main(analyzer.NewAnalyzer())
}
