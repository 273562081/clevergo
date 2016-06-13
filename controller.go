package clevergo

import (
	"reflect"
)

type Controller interface {
	Init(action *Action, ctx *Context)
	BeforeAction() bool
	BeforeResponse()
}

type ControllerInfo struct {
	fullName   string
	name       string
	prettyName string
	t          reflect.Type
	pkgPath    string
	layout     string
	viewsPath  string
}

func (ci *ControllerInfo) FullName() string {
	return ci.fullName
}

func (ci *ControllerInfo) Name() string {
	return ci.name
}

func (ci *ControllerInfo) PrettyName() string {
	return ci.prettyName
}

func (ci *ControllerInfo) PkgPath() string {
	return ci.pkgPath
}

func (ci *ControllerInfo) Layout() string {
	return ci.layout
}

func (ci *ControllerInfo) ViewsPath() string {
	return ci.viewsPath
}
