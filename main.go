package main

import (
	"os"

	"honnef.co/go/tools/lint/lintutil"

	"github.com/daisuzu/gsc/checker"
)

func main() {
	var exitNonZero bool

	fs := lintutil.FlagSet("gsc")
	fs.BoolVar(&exitNonZero,
		"exit-non-zero", true, "Exit non-zero if any problems were found")
	fs.Parse(os.Args[1:])

	cfg := lintutil.CheckerConfig{
		Checker:     checker.New(),
		ExitNonZero: exitNonZero,
	}
	lintutil.ProcessFlagSet([]lintutil.CheckerConfig{cfg}, fs)
}
