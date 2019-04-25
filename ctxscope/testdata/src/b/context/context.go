package context

type Context interface {
}

type ContextBuilder interface {
	Context() Context
}

type contextImpl struct {
}

func New() Context {
	return &contextImpl{}
}

func FromMessage(message map[string]interface{}) Context {
	builder := Builder(New())
	return builder.Context()
}

func Builder(ctx Context) ContextBuilder {
	builder := &contextBuilder{ctx: &contextImpl{}}
	return builder
}

type contextBuilder struct {
	ctx *contextImpl
}

func (b *contextBuilder) Context() Context {
	return b.ctx
}
