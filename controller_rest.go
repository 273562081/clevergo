package clevergo

import (
	"encoding/json"
	"github.com/clevergo/log"
)

var RestHTTPMethods = map[string]string{"GET": "", "PATCH": "", "POST": "", "PUT": "", "DELETE": ""}

type RestController struct {
	Action  Action // current action's info.
	Context *Context
}

func (rc *RestController) Init(action Action, ctx *Context) {
	rc.Action = action
	rc.Context = ctx
}

func (rc *RestController) App() *Application {
	return rc.Context.App
}

func (rc *RestController) Log() *log.Log {
	return rc.Context.Log
}

func (rc *RestController) Request() *Request {
	return rc.Context.Request
}

func (rc *RestController) Response() *Response {
	return rc.Context.Response
}

func (rc *RestController) Info() *ControllerInfo {
	return rc.Action.Controller()
}

func (rc *RestController) BeforeAction() bool {
	return true
}

func (rc *RestController) BeforeResponse() {
}

func (rc *RestController) SkipMiddlewares() map[string]SkipMiddlewares {
	return map[string]SkipMiddlewares{}
}

// the v will be responsed directly if type of v is string.
func (wc *WebController) RenderJson(v interface{}) {
	wc.Context.Response.SetJsonHeader()

	if value, ok := v.(string); ok {
		wc.Context.Response.body = value
	} else {
		json, err := json.Marshal(v)
		if err != nil {
			// wc.Response.InternalServerError(err.Error())
			return
		}
		wc.Context.Response.body += string(json)
	}
}

// the v will be responsed directly if type of v is string.
func (wc *WebController) RenderJsonp(v interface{}, callback string) {
	wc.Context.Response.SetJsonpHeader()

	if value, ok := v.(string); ok {
		wc.Context.Response.body = value
	} else {

		json, err := json.Marshal(v)
		if err != nil {
			// wc.Response.InternalServerError(err.Error())
			return
		}
		wc.Context.Response.body += callback + "(" + string(json) + ")"
	}
}
