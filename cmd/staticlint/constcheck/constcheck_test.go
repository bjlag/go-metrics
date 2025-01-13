package constcheck_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/bjlag/go-metrics/cmd/staticlint/constcheck"
)

func TestAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), constcheck.ConstCheckAnalyzer, "./...")
}
