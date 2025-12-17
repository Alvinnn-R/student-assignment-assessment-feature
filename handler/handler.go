package handler

import (
	"html/template"
	"session-17/service"
)

type Handler struct {
	HandlerAuth       AuthHandler
	HandlerMenu       MenuHandler
	AssignmentHandler AssignmentHandler
	SubmissionHandler SubmissionHandler
}

func NewHandler(service service.Service, templates *template.Template) Handler {
	return Handler{
		HandlerAuth:       NewAuthHandler(service.AuthService, templates),
		HandlerMenu:       NewMenuHandler(templates),
		AssignmentHandler: NewAssignmentHandler(templates, service.AssignmentService),
		SubmissionHandler: NewSubmissionHandler(templates, service.SubmissionService),
	}
}
