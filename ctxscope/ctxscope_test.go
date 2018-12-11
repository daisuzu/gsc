package ctxscope_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/daisuzu/gsc/ctxscope"
)

func init() {
	ctxscope.Analyzer.Flags.Set("target-context", "MyCtx")
}

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, ctxscope.Analyzer, "a")
}
