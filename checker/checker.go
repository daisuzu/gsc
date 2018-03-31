package checker

import (
	"honnef.co/go/tools/lint"
)

type options struct {
	contextNames []string
}

var defaultOptions = options{
	contextNames: []string{"context.Context"},
}

// Option sets checker options.
type Option func(*options)

// WithAdditionalContexts returns a Option which sets the type name of the
// Context userd in CheckCtxScope.
func WithAdditionalContexts(ctxs ...string) Option {
	return func(o *options) {
		o.contextNames = append(o.contextNames, ctxs...)
	}
}

// New returns a new lint.Checker.
func New(opt ...Option) lint.Checker {
	opts := defaultOptions
	for _, o := range opt {
		o(&opts)
	}
	return &checker{opts: opts}
}

type checker struct {
	opts options
}

func (*checker) Name() string            { return "gsc" }
func (*checker) Prefix() string          { return "GSC" }
func (*checker) Init(prog *lint.Program) {}

func (c *checker) Funcs() map[string]lint.Func {
	return map[string]lint.Func{
		"CtxScope": c.CheckCtxScope,
	}
}
