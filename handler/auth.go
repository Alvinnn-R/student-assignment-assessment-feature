package handler

import (
	"html/template"
	"net/http"
	"session-17/service"
	"strconv"
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
	if err := h.Templates.ExecuteTemplate(w, "login", nil); err != nil {
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

	user, err := h.AuthService.Login(email, password)
	if err != nil {
		h.Templates.ExecuteTemplate(w, "login", map[string]interface{}{
			"Error": err.Error(),
		})
		return
	}

	// cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "lumos-" + strconv.Itoa(user.ID),
		Path:     "/",
		HttpOnly: true,
	})

	http.Redirect(w, r, "/user/home", http.StatusSeeOther)
}

func (h *AuthHandler) LogoutView(w http.ResponseWriter, r *http.Request) {
	if err := h.Templates.ExecuteTemplate(w, "logout", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
