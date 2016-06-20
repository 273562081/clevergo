package clevergo

import (
	"fmt"
	"github.com/clevergo/jwt"
	"github.com/clevergo/log"
	"github.com/clevergo/session"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"strings"
)

type Context struct {
	app             *Application                // Application.
	Response        *Response                   // Response.
	Request         *Request                    // Request.
	Params          Params                      // URL params.
	Session         *session.Session            // Session.
	Log             *log.Log                    // Log.
	Values          map[interface{}]interface{} // Share data of middlewares.
	Token           *jwt.Token                  // JWT(JSON WEB TOKEN).
	SkipMiddlewares SkipMiddlewares             // List of middlewares those can be skip.
}

func NewContext(app *Application, rw http.ResponseWriter, r *http.Request, params httprouter.Params) *Context {
	r.ParseForm()
	return &Context{
		app:             app,
		Response:        NewResponse(rw),
		Request:         NewRequest(r),
		Params:          NewParams(params),
		Session:         nil,
		Values:          make(map[interface{}]interface{}, 0),
		SkipMiddlewares: nil,
	}
}

func (ctx *Context) JWT() *jwt.JWT {
	return ctx.app.jwt
}

func (ctx *Context) GetSession() error {
	var err error
	ctx.Session, err = ctx.app.sessionStore.Get(ctx.Request, Configuration.sessionName)
	if err != nil {
		ctx.Session, err = ctx.app.sessionStore.New(Configuration.sessionName)
	}
	return err
}

func (ctx *Context) Flush() {
	// send response status and headers.
	ctx.Response.writer.WriteHeader(ctx.Response.status)

	// send response body.
	fmt.Fprint(ctx.Response.writer, ctx.Response.body)
}

type Params struct {
	httprouter.Params
}

func NewParams(params httprouter.Params) Params {
	return Params{params}
}

func (ps Params) String(name string) string {
	return ps.ByName(name)
}

// Returns param's integer value by name.
// If error reached, returns zero and error.
func (ps Params) Int(name string) (int, error) {
	value, err := strconv.Atoi(ps.ByName(name))
	if err != nil {
		return 0, err
	}
	return value, nil
}

// Returns param's boolean value by name.
// Returns true if the param's value is equal to "true"(case insensitive) or nonzero,
// Otherwise returns false.
func (ps Params) Bool(name string) bool {
	value := ps.Params.ByName(name)
	if len(value) == 0 {
		return false
	}
	return strings.EqualFold(value, "true") || (value != "0")
}
