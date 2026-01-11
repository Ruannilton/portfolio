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
	"portfolio/internal/search"
	"portfolio/web"
)

type WebService struct {
	authService      *auth.AuthService
	portfolioService *portfolio.PortfolioService
	searchService    search.SearchService
}

func (module *WebService) RenderSearchPage(w http.ResponseWriter) error {
	tmpl, err := web.ParseTemplate("pages/search_page.html", "profile_search_query_builder_form.html")
	if err != nil {
		log.Printf("Error parsing search page template: %v", err)
		return err
	}
	return tmpl.ExecuteTemplate(w, "base", nil)
}

func (module *WebService) RenderPortfolioSearchResults(ctx context.Context, w http.ResponseWriter, query search.ProfileSearchQueryBuilder) error {
	searchResult, err := module.searchService.SearchProfiles(&query)

	if err != nil {
		log.Printf("RenderPortfolioSearchResults error: %v", err)
		http.Error(w, "Failed to search profiles", http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl, err := web.ParseTemplateFragment("components/profile_search_response_card.html")
	if err != nil {
		log.Printf("Error parsing search results template: %v", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return err
	}

	return tmpl.ExecuteTemplate(w, "profile_search_results", searchResult)
}

func NewWebService(authService *auth.AuthService, portfolioService *portfolio.PortfolioService, searchService search.SearchService) *WebService {
	return &WebService{
		authService:      authService,
		portfolioService: portfolioService,
		searchService:    searchService,
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

// RenderPublicProfilePage renderiza a página pública de visualização de um perfil (sem autenticação)
func (module *WebService) RenderPublicProfilePage(ctx context.Context, w http.ResponseWriter, profileID string) {
	// Busca o perfil pelo ID
	profile, err := module.portfolioService.GetProfile(ctx, profileID)
	if err != nil {
		if errors.Is(err, portfolio.ErrProfileNotFound) {
			http.Error(w, "Perfil não encontrado", http.StatusNotFound)
			return
		}
		log.Printf("RenderPublicProfilePage error: %v", err)
		http.Error(w, "Falha ao carregar perfil", http.StatusInternalServerError)
		return
	}

	// Busca os dados do usuário dono do perfil
	user, err := module.authService.GetUserByID(ctx, profile.UserID)
	if err != nil {
		log.Printf("RenderPublicProfilePage error fetching user: %v", err)
		http.Error(w, "Falha ao carregar dados do usuário", http.StatusInternalServerError)
		return
	}

	// Monta o DTO para a view
	profileImage := ""
	if user.ProfileImage != nil {
		profileImage = *user.ProfileImage
	}

	viewData := PublicProfileView{
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		ProfileImage: profileImage,
		Profile:      profile,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl, err := web.ParseTemplate("pages/public_profile.html", "read_only_portfolio_content.html")
	if err != nil {
		log.Printf("Error parsing public_profile template: %v", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, "base", viewData)
}

// ==================== Portfolio HTML Rendering ====================

func (module *WebService) RenderPortfolioHTML(ctx context.Context, w http.ResponseWriter, userID string) {

}

func (module *WebService) RenderEditablePortfolioHTML(ctx context.Context, w http.ResponseWriter) {
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
