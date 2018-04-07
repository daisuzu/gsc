package main

import (
	"os"

	"honnef.co/go/tools/lint/lintutil"

	"github.com/daisuzu/gsc/checker"
)

type contexts []string

func (v *contexts) Set(s string) error {
	*v = append(*v, s)
	return nil
}

func (v *contexts) Get() interface{} { return *v }
func (v *contexts) String() string   { return "<contexts>" }

func main() {
	var exitNonZero bool
	var ctxs contexts

	fs := lintutil.FlagSet("gsc")
	fs.BoolVar(&exitNonZero,
		"exit-non-zero", true, "Exit non-zero if any problems were found")
	fs.Var(&ctxs,
		"target-context", "Additional target type for CtxScope other than the standard library's context")
	fs.Parse(os.Args[1:])

	opts := []checker.Option{}
	if len(ctxs) > 0 {
		opts = append(opts, checker.WithAdditionalContexts(ctxs...))
	}

	cfg := lintutil.CheckerConfig{
		Checker:     checker.New(opts...),
		ExitNonZero: exitNonZero,
	}
	lintutil.ProcessFlagSet([]lintutil.CheckerConfig{cfg}, fs)
}
