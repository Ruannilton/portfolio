package web

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"portfolio/internal/auth"
	"portfolio/internal/jwt"
	"portfolio/internal/portfolio"
	"portfolio/web"
)

type WebService struct {
	authService      *auth.AuthService
	portfolioService *portfolio.PortfolioService
}

func NewWebService(authService *auth.AuthService, portfolioService *portfolio.PortfolioService) *WebService {
	return &WebService{
		authService:      authService,
		portfolioService: portfolioService,
	}
}

func (module *WebService) RenderLoginPage(w io.Writer) error {

	tmpl, err := web.ParseTemplate("pages/login.html")
	if err != nil {
		log.Printf("Error parsing login template: %v", err)
		return err
	}
	tmpl.ExecuteTemplate(w, "base", nil)
	return nil
}

func (module *WebService) RenderAppPage(ctx context.Context, w io.Writer) error {

	user, err := module.authService.GetUserFromContext(ctx)

	if err != nil {
		return err
	}

	tmpl, err := web.ParseTemplate("pages/app.html", "portfolio.html")
	if err != nil {
		log.Printf("Error parsing app template: %v", err)
		return err
	}
	tmpl.ExecuteTemplate(w, "base", user)
	return nil
}

// ==================== Portfolio HTML Rendering ====================

func (module *WebService) RenderPortfolioHTML(ctx context.Context, w http.ResponseWriter) {
	user := jwt.GetUserCurrentUser(ctx)
	userID := user.UserID

	profile, err := module.portfolioService.GetMyProfile(ctx, userID)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err != nil {
		if errors.Is(err, portfolio.ErrProfileNotFound) {
			// Renderiza template de portfolio vazio
			tmpl, tmplErr := web.ParseTemplateFragment("components/portfolio_empty.html")
			if tmplErr != nil {
				log.Printf("Error parsing portfolio_empty template: %v", tmplErr)
				http.Error(w, "Failed to render template", http.StatusInternalServerError)
				return
			}
			tmpl.ExecuteTemplate(w, "portfolio_empty", nil)
			return
		}
		log.Printf("RenderPortfolioHTML error: %v", err)
		http.Error(w, "Failed to get profile", http.StatusInternalServerError)
		return
	}

	module.renderProfileContent(w, profile)
}

func (module *WebService) CreateAndRenderPortfolioHTML(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	user := jwt.GetUserCurrentUser(ctx)
	userID := user.UserID

	var input portfolio.SaveProfileInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Printf("CreateAndRenderPortfolioHTML decode error: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	profile, err := module.portfolioService.CreateProfile(ctx, userID, input)
	if err != nil {
		if errors.Is(err, portfolio.ErrProfileAlreadyExists) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		log.Printf("CreateAndRenderPortfolioHTML error: %v", err)
		http.Error(w, "Failed to create profile", http.StatusInternalServerError)
		return
	}

	module.renderProfileContent(w, profile)
}

func (module *WebService) UpdateAndRenderPortfolioHTML(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	user := jwt.GetUserCurrentUser(ctx)
	userID := user.UserID

	var input portfolio.SaveProfileInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Printf("UpdateAndRenderPortfolioHTML decode error: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	profile, err := module.portfolioService.UpdateProfile(ctx, userID, input)
	if err != nil {
		if errors.Is(err, portfolio.ErrProfileNotFound) {
			http.Error(w, "Profile not found", http.StatusNotFound)
			return
		}
		log.Printf("UpdateAndRenderPortfolioHTML error: %v", err)
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	module.renderProfileContent(w, profile)
}

func (module *WebService) renderProfileContent(w http.ResponseWriter, profile *portfolio.Profile) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl, err := web.ParseTemplateFragment("components/portfolio_content.html")
	if err != nil {
		log.Printf("Error parsing portfolio_content template: %v", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, "portfolio_content", profile)
}

// RenderPortfolioPrint renderiza a versão para impressão/PDF do portfolio (endpoint público)
func (module *WebService) RenderPortfolioPrint(ctx context.Context, w http.ResponseWriter, userID string) {
	profile, err := module.portfolioService.GetMyProfile(ctx, userID)
	if err != nil {
		if errors.Is(err, portfolio.ErrProfileNotFound) {
			http.Error(w, "Portfolio não encontrado", http.StatusNotFound)
			return
		}
		log.Printf("RenderPortfolioPrint error: %v", err)
		http.Error(w, "Falha ao carregar portfolio", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl, err := web.ParseTemplateFragment("components/portfolio_print.html")
	if err != nil {
		log.Printf("Error parsing portfolio_print template: %v", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, "portfolio_print", profile)
}
