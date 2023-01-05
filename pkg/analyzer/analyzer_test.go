package analyzer_test

import (
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/AntonBraer/constantcheck/pkg/analyzer"
)

func TestAll(t *testing.T) {
	t.Parallel()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get wd: %s", err)
	}

	testdata := filepath.Join(filepath.Dir(filepath.Dir(wd)), "/pkg/dataset")

	a := analyzer.NewAnalyzer()

	analysistest.Run(t, testdata, a)
}
