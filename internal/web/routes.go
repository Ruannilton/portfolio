package web

import (
	"net/http"
	"portfolio/internal/auth"
	"portfolio/internal/jwt"
	"portfolio/internal/portfolio"
	"portfolio/internal/search"
	"portfolio/web"

	"github.com/gorilla/mux"
)

type WebModule struct {
	authService      *auth.AuthService
	jwtService       *jwt.JWTService
	portfolioService *portfolio.PortfolioService
	webService       *WebService
}

func NewWebModule(authService *auth.AuthService, jwtService *jwt.JWTService, portfolioService *portfolio.PortfolioService, searchService search.SearchService) *WebModule {
	return &WebModule{
		authService:      authService,
		jwtService:       jwtService,
		portfolioService: portfolioService,

		webService: NewWebService(authService, portfolioService, searchService),
	}
}

func (m *WebModule) SetupFrontEnd(router *mux.Router) {
	fileServer := http.FileServer(web.GetStaticAssets())
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fileServer)).Methods("GET")
	// Página Raiz
	router.HandleFunc("/", m.rootPageEndpoint).Methods("GET")

	// Página de Login/Cadastro
	router.HandleFunc("/app/login", m.loginPageEndpoint).Methods("GET")

	// Pagina de perfil
	router.HandleFunc("/app/profile", m.requireAuth(m.profilePageEndpoint)).Methods("GET")
	router.HandleFunc("/app/profile", m.requireAuth(m.createProfileEndpoint)).Methods("POST")
	router.HandleFunc("/app/profile", m.requireAuth(m.updateProfileEndpoint)).Methods("PUT")

	// Página de Busca
	router.HandleFunc("/app/search", m.optionalAuth(m.searchPageEndpoint)).Methods("GET")
	router.HandleFunc("/app/search/results", m.searchResultHandler).Methods("GET")

	// Página pública de visualização de perfil
	router.HandleFunc("/app/profile/{profile_id}", m.optionalAuth(m.publicProfileHandler)).Methods("GET")
	router.HandleFunc("/app/profile/{profile_id}/print", m.portfolioPrintHandler).Methods("GET")

}

func (m *WebModule) rootPageEndpoint(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/app/login", http.StatusSeeOther)
}