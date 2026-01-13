package web

import (
	"io"
	"log"
	"net/http"
	"portfolio/internal/jwt"
	"portfolio/web"
)

func (m *WebModule) loginPageEndpoint(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("access_token")

	if err == nil && cookie.Value != "" {
		_, tokenErr := jwt.GetUserIDFromToken(cookie.Value, m.jwtService)
		if tokenErr == nil {
			http.Redirect(w, r, "/app/profile", http.StatusFound)
			return
		}
	}

	RenderLoginPage(w)
}

func  RenderLoginPage(w io.Writer) error {

	tmpl, err := web.ParseTemplate("pages/login.html")
	if err != nil {
		log.Printf("Error parsing login template: %v", err)
		return err
	}
	tmpl.ExecuteTemplate(w, "base", nil)
	return nil
}
