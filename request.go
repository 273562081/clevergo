package clevergo

import (
	"net/http"
	"strings"
)

type Request struct {
	*http.Request
}

func NewRequest(r *http.Request) *Request {
	return &Request{r}
}

// Get simulation method.
func (r *Request) SimulateMethod(name string) string {
	method := r.FormValue(name)
	if len(method) == 0 {
		method = r.PostFormValue(method)
	}
	return method
}

// Returns a boolean indicating whether r is a GET request.
func (r *Request) IsGet() bool {
	return strings.EqualFold("GET", r.Method)
}

// Returns a boolean indicating whether r is a POST request.
func (r *Request) IsPost() bool {
	return strings.EqualFold("POST", r.Method)
}

// Returns a boolean indicating whether r is a PUT request.
func (r *Request) IsPut() bool {
	return strings.EqualFold("PUT", r.Method)
}

// Returns a boolean indicating whether r is a HEAD request.
func (r *Request) IsHead() bool {
	return strings.EqualFold("HEAD", r.Method)
}

// Returns a boolean indicating whether r is a DELETE request.
func (r *Request) IsDelete() bool {
	return strings.EqualFold("DELETE", r.Method)
}

// Returns a boolean indicating whether r is a PATCH request.
func (r *Request) IsPatch() bool {
	return strings.EqualFold("PATCH", r.Method)
}

// Returns a boolean indicating whether r is a OPTIONS request.
func (r *Request) IsOptions() bool {
	return strings.EqualFold("OPTIONS", r.Method)
}

// Returns a boolean indicating whether r is a AJAX request.
func (r *Request) IsAjax() bool {
	header := r.Header.Get("X-Requested-With")
	return strings.EqualFold("XMLHttpRequest", header)
}
