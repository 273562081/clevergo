package clevergo

type User interface {
	IsGuest() bool
	ID() int64                      // User ID.
	Name() string                   // User name.
	Email() string                  // User email.
	Values() map[string]interface{} // User extra data.
}

type JWTUser struct {
	values map[string]interface{}
}

func NewJwtUser(r Request) *JWTUser {
	return &JWTUser{}
}

func (u *JWTUser) IsGuest() bool {
	return false
}
