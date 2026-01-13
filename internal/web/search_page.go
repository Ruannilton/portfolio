package web

import (
	"context"
	"log"
	"net/http"
	"portfolio/internal/search"
	"portfolio/web"
)

func (m *WebModule) searchPageEndpoint(w http.ResponseWriter, r *http.Request) {
	RenderSearchPage(w)
}

func (m *WebModule) searchResultHandler(w http.ResponseWriter, r *http.Request) {
	searchQuery := extractSearchForm(r)
	ctx := r.Context()
	RenderPortfolioSearchResults(ctx, w, searchQuery, m.webService.searchService)
}

func RenderSearchPage(w http.ResponseWriter) error {
	tmpl, err := web.ParseTemplate("pages/search_page.html", "profile_search_query_builder_form.html")
	if err != nil {
		log.Printf("Error parsing search page template: %v", err)
		return err
	}
	return tmpl.ExecuteTemplate(w, "base", nil)
}

func RenderPortfolioSearchResults(ctx context.Context, w http.ResponseWriter, query search.ProfileSearchQueryBuilder, searchService search.SearchService) error {
	searchResult, err := searchService.SearchProfiles(&query)

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
