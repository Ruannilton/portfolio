package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"portfolio/internal/auth"
	"portfolio/internal/config"
	"portfolio/internal/database"
	"portfolio/internal/jwt"
	"portfolio/internal/portfolio"
	"portfolio/internal/search"
	"portfolio/internal/web"

	"github.com/gorilla/mux"
)

type Application struct {
	config         config.Config
	db             database.DbService
	authModule     *auth.AuthModule
	porfolioModule *portfolio.PortfolioModule
	webModule      *web.WebModule
}

func NewApplication() *Application {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// meilisearch
	searchService := search.NewSearchService(cfg)

	// Tenta configurar o índice ao iniciar (opcional, mas recomendado para criar filtros)
	if err := searchService.ConfigureIndex(); err != nil {
		log.Printf("Aviso: Falha ao configurar índice do Meilisearch: %v", err)
		// Não damos Fatalf aqui para permitir que a app suba mesmo se o Meili estiver indisponível momentaneamente
	}

	// db
	db := database.New(cfg)

	// jwt
	jwtService := jwt.NewJWTService(cfg)

	// auth
	userRepository := auth.NewUserRepository(db.GetDB())
	authService := auth.NewAuthService(cfg, userRepository, &jwtService)
	authModule := auth.NewAuthModule(authService, &jwtService)

	//portfolio
	portfolioRepository := portfolio.NewProfileRepository(db.GetDB())
	portfolioService := portfolio.NewPortfolioService(portfolioRepository, searchService, userRepository)
	porfolioModule := portfolio.NewPortfolioModule(portfolioService, &jwtService)

	// web
	webModule := web.NewWebModule(authService, &jwtService, portfolioService, searchService)

	app := &Application{
		config:         *cfg,
		db:             db,
		authModule:     authModule,
		porfolioModule: porfolioModule,
		webModule:      webModule,
	}
	return app
}

func (s *Application) RegisterRoutes() http.Handler {
	router := mux.NewRouter()

	s.webModule.SetupFrontEnd(router)

	router.HandleFunc("/health", s.healthHandler)
	router.PathPrefix("/auth").Handler(http.StripPrefix("/auth", s.authModule.RegisterAuthRoutes()))
	router.PathPrefix("/portfolio").Handler(http.StripPrefix("/portfolio", s.porfolioModule.RegisterRoutes()))

	return s.corsMiddleware(router)
}

func (app *Application) BuildHttpServer() *http.Server {

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.Port),
		Handler:      app.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Printf("Application starting on port %d", app.config.Port)

	return server
}
