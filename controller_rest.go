package clevergo

import "github.com/clevergo/log"

type RestController struct {
	Action  *Action // current action's info.
	Context *Context
}

func (rc *RestController) Init(action *Action, ctx *Context) {
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
	return rc.Action.controller
}