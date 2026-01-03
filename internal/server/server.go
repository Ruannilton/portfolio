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
)

type Application struct {
	config     config.Config
	db         database.DbService
	authModule *auth.AuthModule
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

	app := &Application{
		config:     *cfg,
		db:         db,
		authModule: authModule,
	}
	return app
}

func (s *Application) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", s.healthHandler)
	mux.Handle("/", s.authModule.RegisterAuthRoutes())

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
