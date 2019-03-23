package testdata

import (
	"context"
	"net/http"
)

func testUpdateWithTxCtx(c context.Context) {
	datastore.RunInTransaction(c, func(tc context.Context) error {
		if err := datastore.Get(tc, nil, nil); err != nil {
			return err
		}
		_, err := datastore.Put(tc, nil, nil)
		return err
	}, nil)
}

func testUpdateWithCtx(c context.Context) {
	datastore.RunInTransaction(c, func(tc context.Context) error {
		if err := datastore.Get(c, nil, nil); err != nil {
			return err
		}
		_, err := datastore.Put(c, nil, nil)
		return err
	}, nil)
}

func testUpdateWithMyCtx(c context.Context) {
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
		if err := get(ctx); err != nil {
			return err
		}
		return put(ctx)
	}, nil)
}

func testUpdateWithUnregisteredCtx(c context.Context) {
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

func testUseCtxInClosure(c context.Context) {
	func() {
		datastore.Delete(c, nil)
	}()
}

func testUseCtxInDeferredClosure(c context.Context) {
	defer func() {
		datastore.Delete(c, nil)
	}()
}

func testMiddleware(next http.Handler) http.Handler {
	f := func(ctx context.Context) context.Context { return ctx }
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(r.Context()))
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "key", "val")))
		next.ServeHTTP(w, r.WithContext(f(r.Context())))
	})
}
