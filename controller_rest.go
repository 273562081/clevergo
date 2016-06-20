package clevergo

import (
	"encoding/json"
	"encoding/xml"
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
func (rc *RestController) RenderJson(v interface{}) {
	rc.Context.Response.SetJsonHeader()

	if value, ok := v.(string); ok {
		rc.Context.Response.body = value
	} else {
		json, err := json.Marshal(v)
		if err != nil {
			// wc.Response.InternalServerError(err.Error())
			return
		}
		rc.Context.Response.body += string(json)
	}
}

// the v will be responsed directly if type of v is string.
func (rc *RestController) RenderJsonp(v interface{}, callback string) {
	rc.Context.Response.SetJsonpHeader()

	if value, ok := v.(string); ok {
		rc.Context.Response.body = value
	} else {

		json, err := json.Marshal(v)
		if err != nil {
			// wc.Response.InternalServerError(err.Error())
			return
		}
		rc.Context.Response.body += callback + "(" + string(json) + ")"
	}
}

// the v will be responsed directly if type of v is string.
func (rc *RestController) RenderXml(v interface{}, header string) {
	rc.Context.Response.SetXmlHeader()

	if value, ok := v.(string); ok {
		rc.Context.Response.body = value
	} else {
		byteXML, err := xml.MarshalIndent(v, "", `   `)
		if err != nil {
			//wc.Response.InternalServerError(err.Error())
			return
		}

		if len(header) == 0 {
			header = xml.Header
		}

		rc.Context.Response.body = header + string(byteXML)
	}
}
