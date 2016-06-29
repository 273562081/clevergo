package clevergo

type SkipMiddlewares map[string]bool

func NewSkipMiddlewares(middlewares ...string) SkipMiddlewares {
	skipMiddlewares := make(SkipMiddlewares, 0)
	for i := 0; i < len(middlewares); i++ {
		skipMiddlewares[middlewares[i]] = true
	}
	return skipMiddlewares
}

type MethodSkipMiddlewares map[string]SkipMiddlewares

func NewMethodSkipMiddlewares(method string, middlewares SkipMiddlewares) MethodSkipMiddlewares {
	return MethodSkipMiddlewares{
		method: middlewares,
	}
}

// Middleware Interface.
type Middleware interface {
	Handle(next Handler) Handler // handle request.
}
