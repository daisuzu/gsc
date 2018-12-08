package main

import (
	"golang.org/x/tools/go/analysis/multichecker"

	"github.com/daisuzu/gsc/ctxscope"
	"github.com/daisuzu/gsc/rangeptr"
)

func main() {
	multichecker.Main(
		ctxscope.Analyzer,
		rangeptr.Analyzer,
	)
}
