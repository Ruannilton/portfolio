package web

import (
	"bytes"
	"context"
	"errors"
	"log"
	"net/http"
	"portfolio/internal/jwt"
	"portfolio/internal/portfolio"
	"portfolio/web"

	"github.com/gorilla/mux"
)

func (m *WebModule) publicProfileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	profileID := vars["profile_id"]
	ctx := r.Context()
	m.webService.RenderPublicProfilePage(ctx, w, profileID)
}

func (m *WebModule) portfolioPrintHtmlHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	profileID := vars["profile_id"]
	ctx := r.Context()
	m.webService.RenderPortfolioHtmlPrint(ctx, w, profileID)
}

func (m *WebModule) portfolioPrintMarkdownHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	profileID := vars["profile_id"]	
	ctx := r.Context()
	m.webService.RenderPortfolioMarkdown(ctx, w, profileID)
}

func (module *WebService) RenderPublicProfilePage(ctx context.Context, w http.ResponseWriter, profileID string) {
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

	profileOwner, err := module.authService.GetUserByID(ctx, profile.UserID)
	if err != nil {
		log.Printf("RenderPublicProfilePage error fetching user: %v", err)
		http.Error(w, "Falha ao carregar dados do usuário", http.StatusInternalServerError)
		return
	}

	// Verifica se há usuário logado para popular top_bar
	loggedUser := jwt.GetUserCurrentUser(ctx)

	viewData := PageViewData{
		PageTitle:      profileOwner.FirstName + " " + profileOwner.LastName,
		OwnerFirstName: profileOwner.FirstName,
		OwnerLastName:  profileOwner.LastName,
		ProfileExists:  true,
	}

	// Popula dados do usuário logado se existir
	if loggedUser != nil && loggedUser.ID != "" {
		viewData.Authenticated = true
		viewData.LoggedUserFirstName = loggedUser.FirstName
		viewData.LoggedUserLastName = loggedUser.LastName
		if loggedUser.ProfileImageURL != nil {
			viewData.LoggedUserProfileImage = *loggedUser.ProfileImageURL
		}
	}

	if profileOwner.ProfileImage != nil {
		viewData.OwnerProfileImage = *profileOwner.ProfileImage
	}

	viewData.FromProfile(profile)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl, err := web.ParseTemplateHtml("pages/show_profile.html", "top_bar.html", "portfolio_view.html")
	if err != nil {
		log.Printf("Error parsing show_profile template: %v", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "base", viewData); err != nil {
		log.Printf("Error executing show_profile template: %v", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
	buf.WriteTo(w)
}

func (module *WebService) RenderPortfolioHtmlPrint(ctx context.Context, w http.ResponseWriter, profileID string) {
	profile, err := module.portfolioService.GetProfile(ctx, profileID)
	if err != nil {
		if errors.Is(err, portfolio.ErrProfileNotFound) {
			http.Error(w, "Portfolio não encontrado", http.StatusNotFound)
			return
		}
		log.Printf("RenderPortfolioPrint error: %v", err)
		http.Error(w, "Falha ao carregar portfolio", http.StatusInternalServerError)
		return
	}

	profileOwner, err := module.authService.GetUserByID(ctx, profile.UserID)
	if err != nil {
		log.Printf("RenderPortfolioPrint error fetching user: %v", err)
		http.Error(w, "Falha ao carregar dados do usuário", http.StatusInternalServerError)
		return
	}

	viewData := PageViewData{
		PageTitle:      profileOwner.FirstName + " " + profileOwner.LastName,
		OwnerFirstName: profileOwner.FirstName,
		OwnerLastName:  profileOwner.LastName,
		ProfileExists:  true,
	}

	if profileOwner.ProfileImage != nil {
		viewData.OwnerProfileImage = *profileOwner.ProfileImage
	}

	viewData.FromProfile(profile)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl, err := web.ParseTemplateFragmentHtml("print/print_portfolio.html")
	if err != nil {
		log.Printf("Error parsing print_portfolio template: %v", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "portfolio_print", viewData); err != nil {
		log.Printf("Error executing print_portfolio template: %v", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
	buf.WriteTo(w)
}

func (module *WebService) RenderPortfolioMarkdown(ctx context.Context, w http.ResponseWriter, profileID string) {
	// TODO
}