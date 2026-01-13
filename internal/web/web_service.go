package web

import (
	"portfolio/internal/auth"
	"portfolio/internal/portfolio"
	"portfolio/internal/search"
)

type WebService struct {
	authService      *auth.AuthService
	portfolioService *portfolio.PortfolioService
	searchService    search.SearchService
}


func NewWebService(authService *auth.AuthService, portfolioService *portfolio.PortfolioService, searchService search.SearchService) *WebService {
	return &WebService{
		authService:      authService,
		portfolioService: portfolioService,
		searchService:    searchService,
	}
}



