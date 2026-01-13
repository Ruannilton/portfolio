package web

import "portfolio/internal/portfolio"

type PageViewData struct {
	Authenticated bool
	PageTitle     string

	LoggedUserFirstName    string
	LoggedUserLastName     string
	LoggedUserProfileImage string

	OwnerFirstName    string
	OwnerLastName     string
	OwnerProfileImage string
	OwnerId           string

	ProfileExists bool

	ProfileID         string
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
	p.ProfileID = profile.ID
	p.OwnerId = profile.UserID
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
