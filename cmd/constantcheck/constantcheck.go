package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/AntonBraer/constantcheck/pkg/analyzer"
)

func main() {
	singlechecker.Main(analyzer.NewAnalyzer())
}
