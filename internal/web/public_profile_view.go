package web

import "portfolio/internal/portfolio"

// PublicProfileView é o DTO para renderizar a página pública de perfil
// Combina dados do Profile com dados do User (nome e imagem)
type PublicProfileView struct {
	FirstName    string
	LastName     string
	ProfileImage string
	Profile      *portfolio.Profile
}
