package clevergo

type Action interface {
	Handle(*Context)
	Controller() *ControllerInfo
	App() *Application
	PrettyName() string
}

func getActionHandler(a Action) Handler {
	finalHandler := HandlerFunc(a.Handle)

	middlewareLen := len(a.App().middlewares)
	if middlewareLen > 0 {
		handler := a.App().middlewares[middlewareLen-1].Handle(finalHandler)
		for i := middlewareLen - 2; i >= 0; i-- {
			handler = a.App().middlewares[i].Handle(handler)
		}
		return handler
	}

	return finalHandler
}
