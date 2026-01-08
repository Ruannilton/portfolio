package web

import (
	"context"
	"net/http"
	"portfolio/internal/auth"
	"portfolio/internal/jwt"
	"portfolio/internal/portfolio"
	"portfolio/web"
	"strings"

	"github.com/gorilla/mux"
)

type WebModule struct {
	authService      *auth.AuthService
	jwtService       *jwt.JWTService
	portfolioService *portfolio.PortfolioService
	webService       *WebService
}

func NewWebModule(authService *auth.AuthService, jwtService *jwt.JWTService, portfolioService *portfolio.PortfolioService) *WebModule {
	return &WebModule{
		authService:      authService,
		jwtService:       jwtService,
		portfolioService: portfolioService,
		webService:       NewWebService(authService, portfolioService),
	}
}

func (m *WebModule) SetupFrontEnd(router *mux.Router) {
	fileServer := http.FileServer(web.GetStaticAssets())
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fileServer)).Methods("GET")

	// Página de Login/Cadastro
	router.HandleFunc("/login", m.loginHandler).Methods("GET")

	// Página do App (protegida)
	router.HandleFunc("/app", m.requireAuth(m.appHandler)).Methods("GET")

	// Rotas HTML do Portfolio (HTMX)
	router.HandleFunc("/portfolio/me/html", m.requireAuth(m.portfolioHTMLHandler)).Methods("GET")
	router.HandleFunc("/portfolio/html", m.requireAuth(m.createPortfolioHTMLHandler)).Methods("POST")
	router.HandleFunc("/portfolio/html", m.requireAuth(m.updatePortfolioHTMLHandler)).Methods("PUT")

	// Rota pública para impressão/PDF do portfolio
	router.HandleFunc("/portfolio/{user_id}/print", m.portfolioPrintHandler).Methods("GET")
}

func (m *WebModule) loginHandler(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("access_token")

	if err == nil && cookie.Value != "" {
		_, tokenErr := jwt.GetUserIDFromToken(cookie.Value, m.jwtService)
		if tokenErr == nil {
			http.Redirect(w, r, "/app", http.StatusFound)
			return
		}
	}

	m.webService.RenderLoginPage(w)
}

func (m *WebModule) appHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	m.webService.RenderAppPage(ctx, w)
}

func (m *WebModule) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		autenticatedUser := m.getAutenticatedUserFromRequest(r)

		if autenticatedUser == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		ctx := context.WithValue(r.Context(), jwt.AutenticatedUserKey, *autenticatedUser)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func getJwtTokenFromRequest(r *http.Request) string {
	cookie, err := r.Cookie("access_token")
	if err == nil {
		return cookie.Value
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	tokenString := parts[1]
	return tokenString
}

func (m *WebModule) getAutenticatedUserFromRequest(r *http.Request) *jwt.AutenticatedUser {
	token := getJwtTokenFromRequest(r)

	if token == "" {
		return nil
	}

	userID, err := jwt.GetUserIDFromToken(token, m.jwtService)
	if err != nil || userID == nil || *userID == "" {
		return nil
	}

	userEmail, emailErr := jwt.GetUserEmailFromToken(token, m.jwtService)
	if emailErr != nil || userEmail == nil || *userEmail == "" {
		return nil
	}

	autenticatedUser := jwt.AutenticatedUser{
		UserID:    *userID,
		UserEmail: *userEmail,
	}
	return &autenticatedUser
}

// ==================== Portfolio HTML Handlers ====================

func (m *WebModule) portfolioHTMLHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	m.webService.RenderPortfolioHTML(ctx, w)
}

func (m *WebModule) createPortfolioHTMLHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	m.webService.CreateAndRenderPortfolioHTML(ctx, w, r)
}

func (m *WebModule) updatePortfolioHTMLHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	m.webService.UpdateAndRenderPortfolioHTML(ctx, w, r)
}

func (m *WebModule) portfolioPrintHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]
	ctx := r.Context()
	m.webService.RenderPortfolioPrint(ctx, w, userID)
}
