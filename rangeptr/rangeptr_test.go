package rangeptr_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/daisuzu/gsc/rangeptr"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, rangeptr.Analyzer, "a")
}
