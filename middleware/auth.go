package middleware

import (
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session")
		if err != nil {
			http.Redirect(w, r, "/page401", http.StatusSeeOther)
			return
		}

		if c.Value == "" {
			http.Redirect(w, r, "/page401", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func AuthMiddlewareV2(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session")
		if err != nil {
			http.Redirect(w, r, "/page401", http.StatusSeeOther)
			return
		}

		if c.Value == "" {
			http.Redirect(w, r, "/page401", http.StatusSeeOther)
			return
		}
		next(w, r)
	})
}
