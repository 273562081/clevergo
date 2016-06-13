package clevergo

type Action interface {
	Handle(*Context)
	Controller() *ControllerInfo
	App() *Application
	PrettyName() string
}