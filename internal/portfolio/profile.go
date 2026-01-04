package portfolio

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// Enums
type Seniority string

type LocationType string

const (
	Junior    Seniority = "JUNIOR"
	MidLevel  Seniority = "MID_LEVEL"
	Senior    Seniority = "SENIOR"
	Lead      Seniority = "LEAD"
	Principal Seniority = "PRINCIPAL"
)

const (
	LocationOnSite LocationType = "ON_SITE"
	LocationRemote LocationType = "REMOTE"
	LocationHybrid LocationType = "HYBRID"
	LocationAny    LocationType = "ANY"
)

type Profile struct {
	ID                string       `json:"id"`
	UserID            string       `json:"userId"`
	Headline          string       `json:"headline"`
	Bio               string       `json:"bio"`
	Seniority         Seniority    `json:"seniority"`
	YearsOfExp        int          `json:"yearsOfExperience"`
	OpenToWork        bool         `json:"openToWork"`
	SalaryExpectation float64      `json:"salaryExpectation"`
	Currency          string       `json:"currency"`
	ContractType      string       `json:"contractType"`
	Location          LocationType `json:"location"`
	RemoteOnly        bool         `json:"remoteOnly"`
	Skills            StringArray  `json:"skills"`
	SocialLinks       SocialLinks  `json:"socialLinks"`
	Experiences       Experiences  `json:"experiences"`
	Projects          Projects     `json:"projects"`
	Educations        Educations   `json:"educations"`
	CreatedAt         time.Time    `json:"createdAt"`
	UpdatedAt         time.Time    `json:"updatedAt"`
}

type PatchProfileDTO struct {
	Headline          *string       `json:"headline,omitempty"`
	Bio               *string       `json:"bio,omitempty"`
	Seniority         *Seniority    `json:"seniority,omitempty"`
	YearsOfExp        *int          `json:"yearsOfExperience,omitempty"`
	OpenToWork        *bool         `json:"openToWork,omitempty"`
	SalaryExpectation *float64      `json:"salaryExpectation,omitempty"`
	Currency          *string       `json:"currency,omitempty"`
	ContractType      *string       `json:"contractType,omitempty"`
	Location          *LocationType `json:"location,omitempty"`
	RemoteOnly        *bool         `json:"remoteOnly,omitempty"`
	Skills            *StringArray  `json:"skills,omitempty"`
	SocialLinks       *SocialLinks  `json:"socialLinks,omitempty"`
	Experiences       *Experiences  `json:"experiences,omitempty"`
	Projects          *Projects     `json:"projects,omitempty"`
	Educations        *Educations   `json:"educations,omitempty"`
}

// --- Sub-structs e Tipos para JSONB ---

type StringArray []string

type SocialLinks struct {
	LinkedIn string `json:"linkedin,omitempty"`
	GitHub   string `json:"github,omitempty"`
	Website  string `json:"website,omitempty"`
}

type Experience struct {
	Company     string     `json:"company"`
	Role        string     `json:"role"`
	StartDate   time.Time  `json:"startDate"`
	EndDate     *time.Time `json:"endDate"`
	Description string     `json:"description"`
	TechStack   []string   `json:"techStack"`
}
type Experiences []Experience

type Project struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	RepoURL     string   `json:"repoUrl"`
	LiveURL     string   `json:"liveUrl"`
	Tags        []string `json:"tags"`
}
type Projects []Project

type Education struct {
	Institution string     `json:"institution"`
	Degree      string     `json:"degree"`
	Field       string     `json:"field"`
	StartDate   time.Time  `json:"startDate"`
	EndDate     *time.Time `json:"endDate"`
}
type Educations []Education

// --- Implementação de Valuer/Scanner para JSONB ---

func (a StringArray) Value() (driver.Value, error) {
	return json.Marshal(a)
}
func (a *StringArray) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}

func (s SocialLinks) Value() (driver.Value, error) {
	return json.Marshal(s)
}
func (s *SocialLinks) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &s)
}

func (e Experiences) Value() (driver.Value, error) {
	return json.Marshal(e)
}
func (e *Experiences) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &e)
}

func (p Projects) Value() (driver.Value, error) {
	return json.Marshal(p)
}
func (p *Projects) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &p)
}

func (ed Educations) Value() (driver.Value, error) {
	return json.Marshal(ed)
}
func (ed *Educations) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &ed)
}

// Factory
func NewProfile(userID string) *Profile {
	return &Profile{
		ID:          uuid.New().String(),
		UserID:      userID,
		Skills:      make(StringArray, 0),
		Experiences: make(Experiences, 0),
		Projects:    make(Projects, 0),
		Educations:  make(Educations, 0),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func (p *Profile) Update(dto PatchProfileDTO) {
	updateIfNotNil(&p.Headline, dto.Headline)
	updateIfNotNil(&p.Bio, dto.Bio)
	updateIfNotNil(&p.Seniority, dto.Seniority)
	updateIfNotNil(&p.YearsOfExp, dto.YearsOfExp)
	updateIfNotNil(&p.OpenToWork, dto.OpenToWork)
	updateIfNotNil(&p.SalaryExpectation, dto.SalaryExpectation)
	updateIfNotNil(&p.Currency, dto.Currency)
	updateIfNotNil(&p.ContractType, dto.ContractType)
	updateIfNotNil(&p.Location, dto.Location)
	updateIfNotNil(&p.RemoteOnly, dto.RemoteOnly)
	updateIfNotNil(&p.Skills, dto.Skills)
	updateIfNotNil(&p.SocialLinks, dto.SocialLinks)
	updateIfNotNil(&p.Experiences, dto.Experiences)
	updateIfNotNil(&p.Projects, dto.Projects)
	updateIfNotNil(&p.Educations, dto.Educations)

	p.UpdatedAt = time.Now()
}

func updateIfNotNil[T any](target *T, source *T) {
	if source != nil {
		*target = *source
	}
}
