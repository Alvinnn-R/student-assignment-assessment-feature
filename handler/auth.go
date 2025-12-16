package handler

import (
	"html/template"
	"net/http"
	"session-17/service"
)

type AuthHandler struct {
	AuthService service.AuthService
	Templates   *template.Template
}

func NewAuthHandler(authHendler service.AuthService, templates *template.Template) AuthHandler {
	return AuthHandler{
		AuthService: authHendler,
		Templates:   templates,
	}
}

func (h *AuthHandler) LoginView(w http.ResponseWriter, r *http.Request) {
	if err := h.Templates.ExecuteTemplate(w, "home", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	_, err := h.AuthService.Login(email, password)
	if err != nil {
		h.Templates.ExecuteTemplate(w, "login", map[string]interface{}{
			"Error": err.Error(),
		})
		return
	}

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}
