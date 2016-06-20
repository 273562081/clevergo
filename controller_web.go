package clevergo

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/clevergo/log"
	"github.com/clevergo/session"
	"github.com/hoisie/mustache"
	"path"
)

type WebController struct {
	EnableLayout bool
	Action       Action // current action's info.
	Context      *Context
}

func (wc *WebController) App() *Application {
	return wc.Context.App
}

func (wc *WebController) Session() *session.Session {
	return wc.Context.Session
}

func (wc *WebController) Log() *log.Log {
	return wc.Context.Log
}

func (wc *WebController) Request() *Request {
	return wc.Context.Request
}

func (wc *WebController) Response() *Response {
	return wc.Context.Response
}

func (wc *WebController) Info() *ControllerInfo {
	return wc.Action.Controller()
}

func (wc *WebController) Init(action Action, ctx *Context) {
	wc.EnableLayout = true
	wc.Action = action
	wc.Context = ctx

	wc.getSession()
}

// Get session.
func (wc *WebController) getSession() {
	if (wc.Action.App().sessionStore != nil) && (wc.Context.Session == nil) {
		err := wc.Context.GetSession()
		if err != nil {
			panic(err)
		}
	}
}

// Save session.
func (wc *WebController) saveSession() {
	if wc.Context.Session != nil {
		err := wc.Action.App().sessionStore.Save(wc.Context.Response.writer, wc.Context.Session)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func (wc *WebController) BeforeAction() bool {
	return true
}

func (wc *WebController) BeforeResponse() {
	wc.saveSession()
}

func (wc *WebController) Actions() map[string]WebActionRoute {
	return map[string]WebActionRoute{}
}

func (wc *WebController) Render(context ...interface{}) {
	wc.RenderFile("", context...)
}

func (wc *WebController) RenderData(data string, context ...interface{}) {
	wc.Context.Response.SetHtmlHeader()
	wc.Context.Response.body = mustache.Render(data, context...)
}

// @param name the view file name.
func (wc *WebController) RenderFile(name string, context ...interface{}) {
	wc.Context.Response.SetHtmlHeader()

	file := wc.getViewFile(name)

	if wc.EnableLayout {
		wc.Context.Response.body = mustache.RenderFileInLayout(file, wc.Action.Controller().layout, context...)
	} else {
		wc.Context.Response.body = mustache.RenderFile(file, context...)
	}
}

func (wc *WebController) RenderPartial(context ...interface{}) {
	wc.RenderPartialFile("", context...)
}

func (wc *WebController) RenderPartialFile(name string, context ...interface{}) {
	file := wc.getViewFile(name)

	wc.Context.Response.body = mustache.RenderFile(file, context...)
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

func (wc *WebController) RenderText(text string) {
	wc.Context.Response.SetHtmlHeader()
	wc.Context.Response.body = text
}

// the v will be responsed directly if type of v is string.
func (wc *WebController) RenderXml(v interface{}, header string) {
	wc.Context.Response.SetXmlHeader()

	if value, ok := v.(string); ok {
		wc.Context.Response.body = value
	} else {
		byteXML, err := xml.MarshalIndent(v, "", `   `)
		if err != nil {
			//wc.Response.InternalServerError(err.Error())
			return
		}

		if len(header) == 0 {
			header = xml.Header
		}

		wc.Context.Response.body = header + string(byteXML)
	}
}

func (wc *WebController) getViewFile(name string) string {
	if len(name) == 0 {
		name = wc.Action.PrettyName() + Configuration.viewSuffix
	} else {
		name = name + Configuration.viewSuffix
	}
	return path.Join(wc.Action.Controller().viewsPath, name)
}

// Get view's path.
// The view's path is relative to the current controller's path.
// For example, controller's path: /app/controllers, the layout's path: /apps/views.
func (wc *WebController) ViewPath() string {
	return "views"
}

// Get layout name.
// Returns false will disable layout.
// Returns true and layout name to enable layout.
// The layout's path is relative to the current controller's views's path.
// For example, views's path: /app/views, the layout's path: /apps/views/layouts/main.html.
func (wc *WebController) Layout() (bool, string) {
	return true, "main.html"
}

func (wc *WebController) SkipMiddlewares() map[string]SkipMiddlewares {
	return map[string]SkipMiddlewares{}
}
