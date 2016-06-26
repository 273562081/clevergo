package clevergo

import (
	"fmt"
	"github.com/clevergo/cache"
	"github.com/clevergo/jwt"
	"github.com/clevergo/log"
	"github.com/clevergo/session"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"path"
	"reflect"
	"strings"
)

type Application struct {
	router          *httprouter.Router
	middlewares     []Middleware
	firstMiddleware Middleware
	actions         []*WebAction
	resources       []*RestAction
	sessionStore    session.Store
	logger          *log.Logger
	cache           *cache.RedisCache
	jwt             *jwt.JWT
	panicHandler    func(http.ResponseWriter, *http.Request, interface{})
}

func NewApplication() *Application {
	return &Application{
		router:          NewRouter(),
		middlewares:     make([]Middleware, 0),
		firstMiddleware: nil,
		actions:         make([]*WebAction, 0),
		resources:       make([]*RestAction, 0),
		sessionStore:    nil,
		logger:          nil,
		cache:           nil,
		panicHandler:    PanicHandler,
	}
}

func (a *Application) SetPanicHandler(handler func(http.ResponseWriter, *http.Request, interface{})) {
	a.panicHandler = handler
}

func (a *Application) SetMethodNotAllowedHandler(handler http.Handler) {
	a.router.MethodNotAllowed = handler
}

func (a *Application) SetNotFoundHandler(handler http.Handler) {
	a.router.NotFound = handler
}

func (a *Application) SetSessionStore(store session.Store) {
	a.sessionStore = store
}

func (a *Application) SetLogger(logger *log.Logger) {
	a.logger = logger
}

func (a *Application) SetMiddlewares(middlewares []Middleware) {
	a.middlewares = middlewares
}

func (a *Application) AddMiddleware(middleware Middleware) {
	a.middlewares = append(a.middlewares, middleware)
}

func (a *Application) RegisterWebControllers(controllers ...WebControllerInterface) {
	for i := 0; i < len(controllers); i++ {
		a.RegisterWebController(controllers[i])
	}
}

func (a *Application) RegisterWebController(c WebControllerInterface) {
	ct := reflect.TypeOf(c)
	cv := reflect.ValueOf(c)

	// Controller's info.
	ci := &ControllerInfo{
		fullName: ct.Elem().Name(),
		t:        cv.Elem().Type(),
		pkgPath:  path.Join(Configuration.srcPath, ct.Elem().PkgPath()),
		layout:   "",
	}

	ci.name = getControllerName(ct.Elem().Name())
	ci.prettyName = PrettyName(ci.name)

	// Views's path.
	ci.viewsPath = path.Join(path.Dir(ci.pkgPath), "views", ci.prettyName)

	// Get EnableLayout, see also @method Layout() of WebController.
	enableLayoutMethod := cv.MethodByName("Layout")
	if enableLayoutMethod.IsValid() {
		values := enableLayoutMethod.Call([]reflect.Value{})
		if len(values) == 2 {
			if enable, ok := values[0].Interface().(bool); ok && enable {
				if layout, ok := values[1].Interface().(string); ok && (len(layout) > 0) {
					ci.layout = path.Join(path.Dir(ci.viewsPath), "layouts", layout)
				}
			}
		}
	}

	// Get actions's route.
	actionsRoute := make(map[string]WebActionRoute, 0)
	actionsMethod := cv.MethodByName("Actions")
	if actionsMethod.IsValid() {
		values := actionsMethod.Call([]reflect.Value{})
		for i := 0; i < len(values); i++ {
			if value, ok := values[i].Interface().(WebActionRoutes); ok {
				for k, v := range value {
					actionsRoute[Configuration.actionPrefix+k+Configuration.actionSuffix] = v
				}
			}
			break
		}
	}

	// Get skip middlewares.
	skipMiddlewares := make(map[string]SkipMiddlewares, 0)
	skipMiddlewaresMethod := cv.MethodByName("SkipMiddlewares")
	if skipMiddlewaresMethod.IsValid() {
		values := skipMiddlewaresMethod.Call([]reflect.Value{})
		for i := 0; i < len(values); i++ {
			if value, ok := values[i].Interface().(map[string]SkipMiddlewares); ok {
				skipMiddlewares = value
			}
			break
		}
	}

	for i := 0; i < ct.NumMethod(); i++ {
		method := ct.Method(i)
		if v, ok := actionsRoute[method.Name]; ok {
			action, err := NewWebAction(a, v.Routes, v.Methods, method.Name, i)

			if err != nil {
				panic(err)
			}

			// Set skip middlewares.
			if middlewares, ok := skipMiddlewares[method.Name]; ok {
				action.skipMiddlewares = middlewares
			}

			action.controller = ci
			a.actions = append(a.actions, action)
		}
	}
}

func (a *Application) RegisterRestController(route string, c RestControllerInterface) {
	ct := reflect.TypeOf(c)
	cv := reflect.ValueOf(c)

	// Controller's info.
	ci := &ControllerInfo{
		fullName: ct.Elem().Name(),
		t:        cv.Elem().Type(),
		pkgPath:  path.Join(Configuration.srcPath, ct.Elem().PkgPath()),
	}

	ci.name = getControllerName(ct.Elem().Name())
	ci.prettyName = PrettyName(ci.name)

	resource := NewRestAction(a, route)
	allowedMethods := RestHTTPMethods

	// Get skip middlewares.
	skipMiddlewares := make(map[string]SkipMiddlewares, 0)
	skipMiddlewaresMethod := cv.MethodByName("SkipMiddlewares")
	if skipMiddlewaresMethod.IsValid() {
		values := skipMiddlewaresMethod.Call([]reflect.Value{})
		for i := 0; i < len(values); i++ {
			if value, ok := values[i].Interface().(map[string]SkipMiddlewares); ok {
				skipMiddlewares = value
			}
			break
		}
	}

	for i := 0; i < ct.NumMethod(); i++ {
		method := ct.Method(i)

		if _, ok := allowedMethods[strings.ToUpper(method.Name)]; ok {
			middlewares := make(SkipMiddlewares, 0)
			if v, ok := skipMiddlewares[strings.ToUpper(method.Name)]; ok {
				middlewares = v
			}

			err := resource.AddMethod(&RestMethod{
				Name:            method.Name,
				Index:           i,
				skipMiddlewares: middlewares,
			})
			if err != nil {
				panic(err)
			}
		}
	}

	if len(resource.methods) == 0 {
		fmt.Printf(`Failed to register restful controller named "%s": no valid methods.\n`, ci.name)
	} else {
		resource.controller = ci
		a.resources = append(a.resources, resource)
	}
}

func (a *Application) RegisterStaticResources(route, path string) {
	a.router.ServeFiles("/"+route+"/*filepath", http.Dir(path))
}

func (a *Application) Run() {
	// Initialize first middleware and final middleware.
	middlewaresLen := len(a.middlewares)
	if middlewaresLen > 0 {
		if middlewaresLen > 1 {
			for i := 0; i < middlewaresLen-1; i++ {
				a.middlewares[i].SetNext(a.middlewares[i+1])
			}
		}
		a.middlewares[0].SetFinal(a.middlewares[middlewaresLen-1])
		a.firstMiddleware = a.middlewares[0]
	}

	// Register web controller's action.
	for i := 0; i < len(a.actions); i++ {
		a.actions[i].handler = GenerateWebActionHandler(a.actions[i])
		for j := 0; j < len(a.actions[i].routes); j++ {
			route := a.actions[i].routes[j]
			for k := 0; k < len(a.actions[i].methods); k++ {
				fmt.Printf("Register web controller's route \"%s\" with method: %s\n", route, a.actions[i].methods[k])
				a.router.Handle(a.actions[i].methods[k], route, a.actions[i].handler)
			}
		}
	}

	// Register restful controller's action.
	for i := 0; i < len(a.resources); i++ {
		a.resources[i].handler = GenerateRestActionHandler(a.resources[i])
		for method, _ := range a.resources[i].methods {
			fmt.Printf("Register restful controller's route \"%s\" with method: %s\n", a.resources[i].route, method)
			a.router.Handle(method, a.resources[i].route, a.resources[i].handler)
		}
	}
}

type Applications map[string]*Application

func (as Applications) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get domain from host.
	host := strings.Split(r.Host, ":")
	if app, ok := as[host[0]]; ok {
		app.router.ServeHTTP(w, r)
	} else {
		defaultApp.router.ServeHTTP(w, r)
	}
}
