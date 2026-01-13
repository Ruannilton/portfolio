package web

import "portfolio/internal/portfolio"

// PageViewData é a struct unificada para renderizar páginas com templates
// Combina dados de autenticação, usuário logado, perfil visualizado e portfolio
type PageViewData struct {
	// === AUTENTICAÇÃO/SESSÃO ===
	Authenticated bool   // Se usuário está autenticado
	PageTitle     string // Título da página

	// === DADOS DO USUÁRIO LOGADO (para top_bar) ===
	// Prefixo "LoggedUser" indica claramente que é o usuário da sessão
	LoggedUserFirstName    string // Primeiro nome do usuário logado
	LoggedUserLastName     string // Sobrenome do usuário logado
	LoggedUserProfileImage string // URL da imagem do usuário logado

	// === DADOS DO DONO DO PORTFOLIO (para show_profile) ===
	// Prefixo "Owner" indica que é o dono do portfolio sendo visualizado
	OwnerFirstName    string // Primeiro nome do dono do perfil
	OwnerLastName     string // Sobrenome do dono do perfil
	OwnerProfileImage string // URL da imagem do dono do perfil

	// === ESTADO DO PORTFOLIO ===
	ProfileExists bool // Se o portfolio existe

	// Dados do portfolio (embeded para acesso direto nos templates)
	ID                string
	Headline          string
	Bio               string
	Seniority         portfolio.Seniority
	YearsOfExp        int
	OpenToWork        bool
	SalaryExpectation float64
	Currency          string
	ContractType      string
	Location          portfolio.LocationType
	RemoteOnly        bool
	Skills            portfolio.StringArray
	SocialLinks       portfolio.SocialLinks
	Experiences       portfolio.Experiences
	Projects          portfolio.Projects
	Educations        portfolio.Educations
}

// FromProfile popula os campos do portfolio a partir de um Profile
func (p *PageViewData) FromProfile(profile *portfolio.Profile) {
	if profile == nil {
		return
	}
	p.ID = profile.ID
	p.Headline = profile.Headline
	p.Bio = profile.Bio
	p.Seniority = profile.Seniority
	p.YearsOfExp = profile.YearsOfExp
	p.OpenToWork = profile.OpenToWork
	p.SalaryExpectation = profile.SalaryExpectation
	p.Currency = profile.Currency
	p.ContractType = profile.ContractType
	p.Location = profile.Location
	p.RemoteOnly = profile.RemoteOnly
	p.Skills = profile.Skills
	p.SocialLinks = profile.SocialLinks
	p.Experiences = profile.Experiences
	p.Projects = profile.Projects
	p.Educations = profile.Educations
}

// PublicProfileView é mantido para compatibilidade (deprecated)
// Use PageViewData para novos templates
type PublicProfileView struct {
	FirstName    string
	LastName     string
	ProfileImage string
	Profile      *portfolio.Profile
}
