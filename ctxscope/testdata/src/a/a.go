package testdata

import (
	"context"
	"net/http"
)

type datastoreInterface struct{}

func (datastoreInterface) RunInTransaction(c context.Context, f func(tc context.Context) error, opts interface{}) error {
	return nil
}

func (datastoreInterface) Get(c context.Context, key, dst interface{}) error {
	return nil
}

func (datastoreInterface) Put(c context.Context, key, src interface{}) (interface{}, error) {
	return nil, nil
}

func (datastoreInterface) Delete(c context.Context, key interface{}) error {
	return nil
}

var datastore = datastoreInterface{}

var ctx, _ = context.WithCancel(context.TODO())

func updateWithTxCtx(c context.Context) {
	datastore.RunInTransaction(c, func(tc context.Context) error {
		if err := datastore.Get(tc, nil, nil); err != nil {
			return err
		}
		_, err := datastore.Put(tc, nil, nil)
		return err
	}, nil)
}

func updateWithCtx(c context.Context) {
	datastore.RunInTransaction(c, func(tc context.Context) error {
		if err := datastore.Get(c, nil, nil); err != nil { // want `passing outer scope context`
			return err
		}
		_, err := datastore.Put(c, nil, nil) // want `passing outer scope context "c" to datastore.Put()`
		return err
	}, nil)
}

func updateWithMyCtx(c context.Context) {
	type MyCtx struct{ context.Context }

	get := func(c *MyCtx) error {
		return datastore.Get(c, nil, nil)
	}
	put := func(c *MyCtx) error {
		_, err := datastore.Put(c, nil, nil)
		return err
	}

	ctx := &MyCtx{c}
	datastore.RunInTransaction(ctx, func(tc context.Context) error {
		if err := get(ctx); err != nil { // want `passing outer scope context "ctx" to get()`
			return err
		}
		return put(ctx) // want `passing outer scope context "ctx" to put()`
	}, nil)
}

func updateWithUnregisteredCtx(c context.Context) {
	type Ctx struct{ context.Context }

	get := func(c *Ctx) error {
		return datastore.Get(c, nil, nil)
	}
	put := func(c *Ctx) error {
		_, err := datastore.Put(c, nil, nil)
		return err
	}

	ctx := &Ctx{c}
	datastore.RunInTransaction(ctx, func(tc context.Context) error {
		if err := get(ctx); err != nil {
			return err
		}
		return put(ctx)
	}, nil)
}

func useCtxInClosure(c context.Context) {
	func() {
		datastore.Delete(c, nil) // want `passing outer scope context "c" to datastore.Delete()`
	}()
}

func useCtxInDeferredClosure(c context.Context) {
	defer func() {
		datastore.Delete(c, nil)
	}()
}

func handler(w http.ResponseWriter, r *http.Request) {
	f := func(ctx context.Context) {}

	f(context.Background())
	f(r.Context())
}

func middleware(next http.Handler) http.Handler {
	f := func(ctx context.Context) context.Context { return ctx }
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(r.Context()))
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "key", "val")))
		next.ServeHTTP(w, r.WithContext(f(r.Context())))
	})
}
