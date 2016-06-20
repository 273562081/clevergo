package clevergo

import "github.com/clevergo/jwt"

// User interface
type User interface {
	IsGuest() bool
	ID() int64                      // User ID.
	Name() string                   // User name.
	Email() string                  // User email.
	Values() map[string]interface{} // User extra data.
}

type JWTUser struct {
	token  *jwt.Token
	values map[string]interface{}
}

func (u *JWTUser) IsGuest() bool {
	return false
}

func (u *JWTUser) ID() int64 {
	return 0
}

func (u *JWTUser) Name() string {
	return ""
}

func (u *JWTUser) Email() string {
	return ""
}

func (u *JWTUser) Values() map[string]interface{} {
	return map[string]interface{}{}
}
