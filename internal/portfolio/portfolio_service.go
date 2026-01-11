package portfolio

import (
	"context"
	"errors"
	"portfolio/internal/auth"
	"portfolio/internal/search"
	"time"
)

type PortfolioService struct {
	repo     ProfileRepository
	search   search.SearchService
	userRepo auth.UserRepository
}

type SaveProfileInput struct {
	Headline          string       `json:"headline"`
	Bio               string       `json:"bio"`
	Seniority         Seniority    `json:"seniority"`
	YearsOfExp        int          `json:"years_of_experience"`
	OpenToWork        bool         `json:"open_to_work"`
	SalaryExpectation float64      `json:"salary_expectation"`
	Currency          string       `json:"currency"`
	ContractType      string       `json:"contract_type"`
	Location          LocationType `json:"location"`
	RemoteOnly        bool         `json:"remote_only"`
	Skills            []string     `json:"skills"`
	SocialLinks       SocialLinks  `json:"social_links"`
	Experiences       Experiences  `json:"experiences"`
	Projects          Projects     `json:"projects"`
	Educations        Educations   `json:"educations"`
}

func NewPortfolioService(repo ProfileRepository, search search.SearchService, userRepo auth.UserRepository) *PortfolioService {
	return &PortfolioService{repo: repo, search: search, userRepo: userRepo}
}

func (s *PortfolioService) GetMyProfile(ctx context.Context, userID string) (*Profile, error) {
	return s.repo.FindByUserID(ctx, userID)
}

func (s *PortfolioService) GetProfile(ctx context.Context, profileID string) (*Profile, error) {
	return s.repo.Find(ctx, profileID)
}

func (s *PortfolioService) ListProfiles(ctx context.Context, profileIDs []string) ([]*Profile, error) {
	return s.repo.List(ctx, profileIDs)
}

func (s *PortfolioService) CreateProfile(ctx context.Context, userID string, input SaveProfileInput) (*Profile, error) {
	// Verifica se já existe
	existing, err := s.repo.FindByUserID(ctx, userID)
	if err != nil && !errors.Is(err, ErrProfileNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, ErrProfileAlreadyExists
	}

	profile := NewProfile(userID)
	s.mapInputToProfile(profile, input)

	if err := s.repo.Create(ctx, profile); err != nil {
		return nil, err
	}

	go s.sendToIndexing(profile, userID)
	return profile, nil
}

func (s *PortfolioService) UpdateProfile(ctx context.Context, userID string, input SaveProfileInput) (*Profile, error) {
	profile, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	s.mapInputToProfile(profile, input)
	profile.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, profile); err != nil {
		return nil, err
	}
	go s.sendToIndexing(profile, userID)
	return profile, nil
}

func (s *PortfolioService) PatchProfile(ctx context.Context, userID string, input PatchProfileDTO) (*Profile, error) {
	profile, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	profile.Update(input)

	if err := s.repo.Update(ctx, profile); err != nil {
		return nil, err
	}
	go s.sendToIndexing(profile, userID)
	return profile, nil
}

func (s *PortfolioService) DeleteProfile(ctx context.Context, userID string) error {
	go s.search.DeleteProfile(userID)
	return s.repo.Delete(ctx, userID)
}

// Helper para mapear DTO -> Entity
func (s *PortfolioService) mapInputToProfile(p *Profile, input SaveProfileInput) {
	p.Headline = input.Headline
	p.Bio = input.Bio
	p.Seniority = input.Seniority
	p.YearsOfExp = input.YearsOfExp
	p.OpenToWork = input.OpenToWork
	p.SalaryExpectation = input.SalaryExpectation
	p.Currency = input.Currency
	p.ContractType = input.ContractType
	p.Location = input.Location
	p.RemoteOnly = input.RemoteOnly
	p.Skills = StringArray(input.Skills)
	p.SocialLinks = input.SocialLinks
	p.Experiences = input.Experiences
	p.Projects = input.Projects
	p.Educations = input.Educations
}

func (s *PortfolioService) sendToIndexing(p *Profile, userId string) {

	user, err := s.userRepo.Find(context.Background(), userId)
	if err != nil {
		return
	}
	skills := make([]string, len(p.Skills))
	copy(skills, p.Skills)

	for _, project := range p.Projects {
		skills = append(skills, project.Tags...)
	}
	userName := user.FirstName + " " + user.LastName
	profileImg := ""
	if user.ProfileImage != nil {
		profileImg = *user.ProfileImage
	}
	dto := search.ProfileSearchDTO{
		ProfileId:         p.ID,
		Headline:          p.Headline,
		Bio:               p.Bio,
		Skills:            skills,
		Seniority:         p.Seniority.Int(),
		YearsOfExp:        p.YearsOfExp,
		Location:          p.Location.Int(),
		OpenToWork:        p.OpenToWork,
		ContractType:      p.ContractType,
		Currency:          p.Currency,
		SalaryExpectation: p.SalaryExpectation,
		UserName:          userName,
		UserProfileImage:  profileImg,
		RemoteOnly:        p.RemoteOnly,
	}
	s.search.IndexProfile(dto)
}
