package search

import (
	"fmt"
	"log"
	"portfolio/internal/config"
	"strings"

	"github.com/meilisearch/meilisearch-go"
)

type ProfileSearchDTO struct {
	ID                string   `json:"id"`
	Headline          string   `json:"headline"`
	Bio               string   `json:"bio"`
	Skills            []string `json:"skills"`
	Seniority         int      `json:"seniority"`
	YearsOfExp        int      `json:"years_of_experience"`
	Location          int      `json:"location"`
	OpenToWork        bool     `json:"open_to_work"`
	ContractType      string   `json:"contract_type"`
	Currency          string   `json:"currency"`
	SalaryExpectation float64  `json:"salary_expectation"`
}

type ProfileSearchQueryBuilder struct {
	keyWords     []string
	skills       []string
	seniority    []int
	location     []int
	remoteOnly   *bool
	openToWork   *bool
	contractType []string
	minYearsExp  *int
	maxYearsExp  *int
	minSalary    *float64
	maxSalary    *float64
}

// NewProfileSearchQueryBuilder creates a new query builder
func NewProfileSearchQueryBuilder() *ProfileSearchQueryBuilder {
	return &ProfileSearchQueryBuilder{}
}

// WithKeyWords sets the keywords for full-text search
func (b *ProfileSearchQueryBuilder) WithKeyWords(keywords ...string) *ProfileSearchQueryBuilder {
	b.keyWords = keywords
	return b
}

// WithSkills filters by skills (OR condition)
func (b *ProfileSearchQueryBuilder) WithSkills(skills ...string) *ProfileSearchQueryBuilder {
	b.skills = skills
	return b
}

// WithSeniority filters by seniority levels (OR condition)
func (b *ProfileSearchQueryBuilder) WithSeniority(seniority ...int) *ProfileSearchQueryBuilder {
	b.seniority = seniority
	return b
}

// WithLocation filters by location (OR condition)
func (b *ProfileSearchQueryBuilder) WithLocation(location ...int) *ProfileSearchQueryBuilder {
	b.location = location
	return b
}

// WithRemoteOnly filters by remote_only flag
func (b *ProfileSearchQueryBuilder) WithRemoteOnly(remoteOnly bool) *ProfileSearchQueryBuilder {
	b.remoteOnly = &remoteOnly
	return b
}

// WithOpenToWork filters by open_to_work flag
func (b *ProfileSearchQueryBuilder) WithOpenToWork(openToWork bool) *ProfileSearchQueryBuilder {
	b.openToWork = &openToWork
	return b
}

// WithContractType filters by contract types (OR condition)
func (b *ProfileSearchQueryBuilder) WithContractType(contractType ...string) *ProfileSearchQueryBuilder {
	b.contractType = contractType
	return b
}

// WithYearsOfExperience sets the range for years of experience
func (b *ProfileSearchQueryBuilder) WithYearsOfExperience(min, max *int) *ProfileSearchQueryBuilder {
	b.minYearsExp = min
	b.maxYearsExp = max
	return b
}

// WithSalaryRange sets the salary range filter
func (b *ProfileSearchQueryBuilder) WithSalaryRange(min, max *float64) *ProfileSearchQueryBuilder {
	b.minSalary = min
	b.maxSalary = max
	return b
}

// BuildQuery returns the full-text search query string
func (b *ProfileSearchQueryBuilder) BuildQuery() string {
	return strings.Join(b.keyWords, " ")
}

// BuildFilter constructs the Meilisearch filter string
func (b *ProfileSearchQueryBuilder) BuildFilter() string {
	var filters []string

	// Skills (OR between skills)
	if len(b.skills) > 0 {
		skillFilters := make([]string, len(b.skills))
		for i, skill := range b.skills {
			skillFilters[i] = fmt.Sprintf("skills = '%s'", skill)
		}
		filters = append(filters, "("+strings.Join(skillFilters, " OR ")+")")
	}

	// Seniority (OR between values)
	if len(b.seniority) > 0 {
		senFilters := make([]string, len(b.seniority))
		for i, sen := range b.seniority {
			senFilters[i] = fmt.Sprintf("seniority = %d", sen)
		}
		filters = append(filters, "("+strings.Join(senFilters, " OR ")+")")
	}

	// Location (OR between values)
	if len(b.location) > 0 {
		locFilters := make([]string, len(b.location))
		for i, loc := range b.location {
			locFilters[i] = fmt.Sprintf("location = %d", loc)
		}
		filters = append(filters, "("+strings.Join(locFilters, " OR ")+")")
	}

	// Remote Only
	if b.remoteOnly != nil {
		filters = append(filters, fmt.Sprintf("remote_only = %t", *b.remoteOnly))
	}

	// Open To Work
	if b.openToWork != nil {
		filters = append(filters, fmt.Sprintf("open_to_work = %t", *b.openToWork))
	}

	// Contract Type (OR between values)
	if len(b.contractType) > 0 {
		ctFilters := make([]string, len(b.contractType))
		for i, ct := range b.contractType {
			ctFilters[i] = fmt.Sprintf("contract_type = '%s'", ct)
		}
		filters = append(filters, "("+strings.Join(ctFilters, " OR ")+")")
	}

	// Years of Experience (range)
	if b.minYearsExp != nil && b.maxYearsExp != nil {
		filters = append(filters, fmt.Sprintf("years_of_experience %d TO %d", *b.minYearsExp, *b.maxYearsExp))
	} else if b.minYearsExp != nil {
		filters = append(filters, fmt.Sprintf("years_of_experience >= %d", *b.minYearsExp))
	} else if b.maxYearsExp != nil {
		filters = append(filters, fmt.Sprintf("years_of_experience <= %d", *b.maxYearsExp))
	}

	// Salary Range
	if b.minSalary != nil && b.maxSalary != nil {
		filters = append(filters, fmt.Sprintf("salary_expectation %f TO %f", *b.minSalary, *b.maxSalary))
	} else if b.minSalary != nil {
		filters = append(filters, fmt.Sprintf("salary_expectation >= %f", *b.minSalary))
	} else if b.maxSalary != nil {
		filters = append(filters, fmt.Sprintf("salary_expectation <= %f", *b.maxSalary))
	}

	return strings.Join(filters, " AND ")
}

type SearchService interface {
	ConfigureIndex() error
	IndexProfile(profile ProfileSearchDTO) error
	DeleteProfile(id string) error
	SearchProfiles(query *ProfileSearchQueryBuilder) ([]string, error)
}

type meiliService struct {
	client    meilisearch.ServiceManager
	indexName string
}

func NewSearchService(cfg *config.Config) SearchService {
	client := meilisearch.New(cfg.MeiliHost, meilisearch.WithAPIKey(cfg.MeiliMasterKey))

	return &meiliService{
		client:    client,
		indexName: "profiles",
	}
}

func (s *meiliService) ConfigureIndex() error {
	index := s.client.Index(s.indexName)

	// Define Atributos Filtráveis (Facets / Filtros)
	// Isso permite queries como: "skills = 'Go' AND remote_only = true"
	filterableAttributes := []interface{}{
		"id",
		"skills",
		"seniority",
		"years_of_experience",
		"salary_expectation",
		"location",
		"remote_only",
		"open_to_work",
		"contract_type",
		"project_tags",
	}
	_, err := index.UpdateFilterableAttributes(&filterableAttributes)
	if err != nil {
		log.Printf("Erro ao configurar filtros do Meili: %v", err)
		return err
	}

	sortableAttributes := []string{
		"years_of_experience",
		"salary_expectation",
	}
	_, err = index.UpdateSortableAttributes(&sortableAttributes)
	if err != nil {
		log.Printf("Erro ao configurar ordenação do Meili: %v", err)
		return err
	}

	log.Println("Meilisearch configurado com sucesso!")
	return nil
}

func (s *meiliService) IndexProfile(profile ProfileSearchDTO) error {

	task, err := s.client.Index(s.indexName).UpdateDocuments([]ProfileSearchDTO{profile}, nil)

	if err != nil {
		log.Printf("Erro ao enviar perfil %s para indexação: %v", profile.ID, err)
		return err
	}

	log.Printf("Perfil %s enviado para fila de indexação. Task ID: %d", profile.ID, task.TaskUID)
	return nil
}

func (s *meiliService) DeleteProfile(id string) error {
	task, err := s.client.Index(s.indexName).DeleteDocument(id, nil)

	if err != nil {
		log.Printf("Erro ao solicitar remoção do perfil %s: %v", id, err)
		return err
	}

	log.Printf("Remoção do perfil %s solicitada. Task ID: %d", id, task.TaskUID)
	return nil
}

func (s *meiliService) SearchProfiles(query *ProfileSearchQueryBuilder) ([]string, error) {
	searchQuery := ""
	var filter string

	if query != nil {
		searchQuery = query.BuildQuery()
		filter = query.BuildFilter()
	}

	searchRequest := &meilisearch.SearchRequest{
		AttributesToRetrieve: []string{"id"},
		Limit:                1000,
	}

	if filter != "" {
		searchRequest.Filter = filter
	}

	searchRes, err := s.client.Index(s.indexName).Search(searchQuery, searchRequest)
	if err != nil {
		log.Printf("Erro ao buscar perfis: %v", err)
		return nil, err
	}

	// Extrair apenas os IDs dos resultados
	type idResult struct {
		ID string `json:"id"`
	}

	var results []idResult
	if err := searchRes.Hits.DecodeInto(&results); err != nil {
		log.Printf("Erro ao decodificar resultados: %v", err)
		return nil, err
	}

	ids := make([]string, len(results))
	for i, r := range results {
		ids[i] = r.ID
	}

	log.Printf("Busca retornou %d perfis", len(ids))
	return ids, nil
}
