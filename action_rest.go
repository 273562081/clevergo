package clevergo

import (
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"reflect"
	"strings"
)

type RestAction struct {
	app        *Application           // resource's application.
	route      string                 // resource's route.
	methods    map[string]*RestMethod // resource's methods.
	controller *ControllerInfo        // resource's controller.
	handler    httprouter.Handle      // resource's handle.
}

type RestMethod struct {
	Name            string
	Index           int
	skipMiddlewares map[string]bool // the middleware those can be skipped.
}

func (ra *RestAction) Controller() *ControllerInfo {
	return ra.controller
}

func (ra *RestAction) App() *Application {
	return ra.app
}

func (ra *RestAction) PrettyName() string {
	return ""
}

func NewRestAction(app *Application, route string) *RestAction {
	ai := &RestAction{
		app:        app,
		route:      route,
		methods:    make(map[string]*RestMethod, 0),
		controller: nil,
		handler:    nil,
	}

	return ai
}
func (ra *RestAction) AddMethod(method *RestMethod) error {
	if ('A' > method.Name[0]) || (method.Name[0] > 'Z') {
		return errors.New("The action's name is invalid: , it's first charater must be a capital letter." + method.Name)
	}
	ra.methods[strings.ToUpper(method.Name)] = method
	return nil
}

func (ra *RestAction) Handle(ctx *Context) {
	// Create controller's reflect value.
	cv := reflect.New(ra.controller.t)

	// Invoke controller's Init() method.
	initMethod := cv.MethodByName("Init")
	initMethod.Call([]reflect.Value{
		reflect.ValueOf(ra),
		reflect.ValueOf(ctx),
	})

	var values []reflect.Value

	// Invoke controller's BeforeAction() method.
	beforeActionMethod := cv.MethodByName("BeforeAction")
	values = beforeActionMethod.Call([]reflect.Value{})
	// The request will be terminated instantly, if BeforeAction() returns false.
	if value, ok := values[0].Interface().(bool); !ok || !value {
		return
	}

	// Invoke controller's action.
	var methodIndex int
	// Get request method name.
	method := ctx.Request.SimulateMethod(Configuration.actionMethod)
	if len(method) > 0 {
		method = strings.ToUpper(method)
		if mv, ok := ra.methods[method]; ok {
			methodIndex = mv.Index
		} else {
			ctx.Response.MethodNotAllowed()
			return
		}
	} else {
		methodIndex = ra.methods[ctx.Request.Method].Index
	}

	actionMethod := cv.Method(methodIndex) // MethodByIndex is faster than MethodByName.
	// actionMethod := cv.MethodByName(a.fullName)
	actionMethod.Call([]reflect.Value{})

	// Invoke controller's BeforeResponse() method.
	beforeResponseMethod := cv.MethodByName("BeforeResponse")
	beforeResponseMethod.Call([]reflect.Value{})

	return
}

func GenerateRestActionHandler(ra *RestAction) httprouter.Handle {
	handler := getActionHandler(ra)

	return func(rw http.ResponseWriter, r *http.Request, params httprouter.Params) {
		ctx := NewContext(ra.app, rw, r, params)
		ctx.SkipMiddlewares = ra.methods[ctx.Request.Method].skipMiddlewares

		defer ctx.Flush()

		if Configuration.enableLog {
			ctx.Log = ra.app.logger.NewLog()
			defer ctx.Log.Flush()
		}

		handler.Handle(ctx)

		return
	}
}
