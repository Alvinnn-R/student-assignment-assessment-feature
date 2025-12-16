package handler

import (
	"html/template"
	"session-17/service"
)

type Handler struct {
	HandlerAuth AuthHandler
	HandlerMenu MenuHandler
}

func NewHandler(service service.Service, templates *template.Template) Handler {
	return Handler{
		HandlerAuth: NewAuthHandler(service.AuthService, templates),
		HandlerMenu: NewMenuHandler(templates),
	}
}
