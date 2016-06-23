package middleware

import (
	"github.com/clevergo/clevergo"
	"github.com/clevergo/jwt"
	"strings"
)

var (
	JWTMiddlewareID = "JWTMiddleware"
	urlKey          = "_jwt"
	formKey         = "_jwt"
)

type JWTMiddleware struct {
	clevergo.BaseMiddleware
	urlKey  string
	formKey string
}

func NewJWTMiddleware() *JWTMiddleware {
	return &JWTMiddleware{
		urlKey:  urlKey,
		formKey: formKey,
	}
}

func (jm *JWTMiddleware) Handle(ctx *clevergo.Context) {
	if _, canSkip := ctx.SkipMiddlewares[JWTMiddlewareID]; canSkip {
		jm.Next().Handle(ctx)
		return
	}

	// Try to get JWT raw token from URL query string.
	rawToken := ctx.Request.FormValue(jm.urlKey)
	if len(rawToken) < 0 {
		// Try to get JWT raw token from POST FORM.
		rawToken = ctx.Request.PostFormValue(jm.formKey)
		if len(rawToken) < 0 {
			// Try to get JWT raw token from Header.
			if ah := ctx.Request.Header.Get("Authorization"); ah != "" {
				// Should be a bearer token
				if len(ah) > 6 && strings.ToUpper(ah[0:7]) == "BEARER " {
					rawToken = ah[7:]
				}
			}
		}
	}

	// Check raw token is valid.
	if len(rawToken) == 0 {
		ctx.Response.Unauthorized()
		return
	}

	// Get JWT by raw token
	token, err := jwt.NewTokenByRaw(ctx.JWT(), rawToken)
	if err != nil {
		ctx.Response.Unauthorized(err.Error())
		return
	}

	// Validate JWT.
	if err = token.Validate(); err != nil {
		ctx.Response.Unauthorized(err.Error())
		return
	}

	ctx.Token = token
	// Validate successfully.
	jm.Next().Handle(ctx)
}
