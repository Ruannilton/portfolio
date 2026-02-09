package web

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"portfolio/internal/search"
	"portfolio/web"
)

func (m *WebModule) searchPageEndpoint(w http.ResponseWriter, r *http.Request) {
	if err := RenderSearchPage(w); err != nil {
		log.Printf("Error rendering search page: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (m *WebModule) searchResultHandler(w http.ResponseWriter, r *http.Request) {
	searchQuery := extractSearchForm(r)
	ctx := r.Context()
	if err := RenderPortfolioSearchResults(ctx, w, searchQuery, m.webService.searchService); err != nil {
		log.Printf("Error rendering search results: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func RenderSearchPage(w http.ResponseWriter) error {
	tmpl, err := web.ParseTemplateHtml("pages/search_page.html", "profile_search_query_builder_form.html")
	if err != nil {
		log.Printf("Error parsing search page template: %v", err)
		return err
	}

	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "base", nil); err != nil {
		return err
	}
	_, err = buf.WriteTo(w)
	return err
}

func RenderPortfolioSearchResults(ctx context.Context, w http.ResponseWriter, query search.ProfileSearchQueryBuilder, searchService search.SearchService) error {
	searchResult, err := searchService.SearchProfiles(&query)

	if err != nil {
		log.Printf("RenderPortfolioSearchResults error: %v", err)
		return err
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl, err := web.ParseTemplateFragmentHtml("components/profile_search_response_card.html")
	if err != nil {
		log.Printf("Error parsing search results template: %v", err)
		return err
	}

	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "profile_search_results", searchResult); err != nil {
		return err
	}
	_, err = buf.WriteTo(w)
	return err
}
