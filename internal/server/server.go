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
	"github.com/gorilla/mux"
)

type Application struct {
	config     config.Config
	db         database.DbService
	authModule *auth.AuthModule
	porfolioModule *portfolio.PortfolioModule
}

func NewApplication() *Application {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// db
	db := database.New(cfg)

	// jwt
	jwtService := jwt.NewJWTService(cfg)

	// auth
	userRepository := auth.NewUserRepository(db.GetDB())
	authService := auth.NewAuthService(cfg, userRepository, jwtService)
	authModule := auth.NewAuthModule(authService, jwtService)
	
	//portfolio
	portfolioRepository := portfolio.NewProfileRepository(db.GetDB())
	portfolioService := portfolio.NewPortfolioService(portfolioRepository)
	porfolioModule := portfolio.NewPortfolioModule(portfolioService, jwtService)

	app := &Application{
		config:     *cfg,
		db:         db,
		authModule: authModule,
		porfolioModule: porfolioModule,
	}
	return app
}

func (s *Application) RegisterRoutes() http.Handler {
	mux := mux.NewRouter()

	mux.HandleFunc("/health", s.healthHandler)
	mux.PathPrefix("/auth").Handler(http.StripPrefix("/auth",s.authModule.RegisterAuthRoutes()))
	mux.PathPrefix("/portfolio").Handler(http.StripPrefix("/portfolio",s.porfolioModule.RegisterRoutes()))

	return s.corsMiddleware(mux)
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
