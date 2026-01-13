package web

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"portfolio/internal/jwt"
	"portfolio/internal/portfolio"
	"portfolio/web"
)

func (m *WebModule) profilePageEndpoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	m.webService.RenderAppPage(ctx, w)
}

func (m *WebModule) createProfileEndpoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	m.webService.CreatePortfolioFragment(ctx, w, r)
}

func (m *WebModule) updateProfileEndpoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	m.webService.UpdatePortfolioFragment(ctx, w, r)
}

func (module *WebService) RenderAppPage(ctx context.Context, w io.Writer) error {
	user, err := module.authService.GetUserFromContext(ctx)
	if err != nil {
		return err
	}

	// Prepara dados da view
	viewData := PageViewData{
		Authenticated:       true,
		PageTitle:           "Meu Portfolio",
		LoggedUserFirstName: user.FirstName,
		LoggedUserLastName:  user.LastName,
	}

	if user.ProfileImage != nil {
		viewData.LoggedUserProfileImage = *user.ProfileImage
	}

	// Tenta buscar o portfolio do usuário
	profile, err := module.portfolioService.GetMyProfile(ctx, user.ID)
	if err != nil {
		if !errors.Is(err, portfolio.ErrProfileNotFound) {
			log.Printf("RenderAppPage error fetching profile: %v", err)
			return err
		}
		// Portfolio não existe ainda
		viewData.ProfileExists = false
	} else {
		viewData.ProfileExists = true
		viewData.FromProfile(profile)
	}

	tmpl, err := web.ParseTemplate("pages/my_profile.html", "top_bar.html", "portfolio_view.html", "portfolio_editor.html")
	if err != nil {
		log.Printf("Error parsing my_profile template: %v", err)
		return err
	}
	return tmpl.ExecuteTemplate(w, "base", viewData)
}

func (module *WebService) CreatePortfolioFragment(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	user := jwt.GetUserCurrentUser(ctx)
	userID := user.ID

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

func (module *WebService) UpdatePortfolioFragment(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	user := jwt.GetUserCurrentUser(ctx)
	userID := user.ID

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
	tmpl, err := web.ParseTemplateFragment("components/portfolio_view.html")
	if err != nil {
		log.Printf("Error parsing portfolio_view template: %v", err)
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, "portfolio_view", profile)
}
