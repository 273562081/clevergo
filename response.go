package clevergo

import "net/http"

type Response struct {
	writer http.ResponseWriter
	status int
	body   string
	cancel bool
}

func NewResponse(rw http.ResponseWriter) *Response {
	return &Response{
		writer: rw,
		status: http.StatusOK,
		body:   "",
		cancel: false,
	}
}

func (r *Response) Status() int {
	return r.status
}

func (r *Response) SetStatus(status int) {
	r.status = status
}

func (r *Response) Header() http.Header {
	return r.writer.Header()
}

func (r *Response) SetHtmlHeader() {
	r.writer.Header().Add("Content-Type", "text/html; charset=utf-8")
}

func (r *Response) SetJsonHeader() {
	r.writer.Header().Add("Content-Type", "application/json; charset=utf-8")
}

func (r *Response) SetJsonpHeader() {
	r.writer.Header().Add("Content-Type", "application/javascript; charset=utf-8")
}

func (r *Response) SetXmlHeader() {
	r.writer.Header().Add("Content-Type", "application/xml; charset=utf-8")
}

func (r *Response) Body() string {
	return r.body
}

func (r *Response) SetBody(body string) {
	r.body = body
}

func (r *Response) Unauthorized(args ...string) {
	r.SetStatus(http.StatusUnauthorized)
	if len(args) > 0 {
		r.SetBody(args[0])
	}
	r.SetBody(http.StatusText(http.StatusUnauthorized))
}

func (r *Response) MethodNotAllowed(args ...string) {
	r.SetStatus(http.StatusMethodNotAllowed)
	if len(args) > 0 {
		r.SetBody(args[0])
	}
	r.SetBody(http.StatusText(http.StatusMethodNotAllowed))
}
