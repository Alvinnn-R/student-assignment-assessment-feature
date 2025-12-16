package router

import (
	"net/http"
	"session-17/handler"
	"session-17/middleware"

	"github.com/go-chi/chi/v5"
)

func NewRouter(handler handler.Handler) *chi.Mux {
	r := chi.NewRouter()

	//authentication
	r.Get("/login", handler.HandlerAuth.LoginView)
	r.Post("/login", handler.HandlerAuth.Login)
	r.Post("/logout", handler.HandlerAuth.Logout)

	//menu
	r.Route("/user", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Get("/home", handler.HandlerMenu.HomeView)
		r.Get("/assignments", handler.AssignmentHandler.List)
		r.Get("/submit-form", handler.AssignmentHandler.SubmitView)
		r.Get("/success-submit", handler.AssignmentHandler.SuccessSubmit)
		r.Post("/submit-assignment", handler.AssignmentHandler.SubmitAssignment)
		r.Get("/grade", handler.HandlerMenu.GradeView)
		r.Get("/logout", handler.HandlerAuth.LogoutView)
	})
	r.Get("/page401", handler.HandlerMenu.PageUnauthorized)

	fs := http.FileServer(http.Dir("public"))
	r.Handle("/public/*", http.StripPrefix("/public/", fs))

	return r
}
