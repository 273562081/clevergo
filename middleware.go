package clevergo

// Middleware Interface.
type Middleware interface {
	Handle(*Context)               // handle request.
	Next() Middleware              // next middleware.
	SetNext(Middleware) Middleware // set next middleware.
	SetFinal(Middleware)           // set final middleware.
	Final() Middleware             // final middleware
}

// BaseMiddleware.
type BaseMiddleware struct {
	next  Middleware
	final Middleware
}

func (bm *BaseMiddleware) Handle(ctx *Context) {
	bm.next.Handle(ctx)
}
func (bm *BaseMiddleware) Next() Middleware {
	return bm.next
}

func (bm *BaseMiddleware) SetNext(next Middleware) Middleware {
	bm.next = next
	return next
}

func (bm *BaseMiddleware) Final() Middleware {
	return bm.final
}

func (bm *BaseMiddleware) SetFinal(final Middleware) {
	bm.final = final
}

// Middleware example.
// Add headers info.
type ExampleMiddleware struct {
	BaseMiddleware
	Key   string
	Value string
}

func NewExampleMiddleware(key, value string) *ExampleMiddleware {
	return &ExampleMiddleware{
		Key:   key,
		Value: value,
	}
}

func (em *ExampleMiddleware) Handle(ctx *Context) {
	ctx.Response.Header().Add(em.Key, em.Value)
	em.next.Handle(ctx)
}
