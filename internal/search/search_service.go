package search

import (
	"log"
	"portfolio/internal/config"

	"github.com/meilisearch/meilisearch-go"
)

type ProfileSearchDTO struct {
	ProfileId         string   `json:"profileId"`
	UserName          string   `json:"username"`
	UserProfileImage  string   `json:"userProfileImage"`
	Headline          string   `json:"headline"`
	Bio               string   `json:"bio"`
	Skills            []string `json:"skills"`
	Seniority         int      `json:"seniority"`
	YearsOfExp        int      `json:"yearsOfExperience"`
	Location          int      `json:"location"`
	OpenToWork        bool     `json:"openToWork"`
	ContractType      string   `json:"contractType"`
	Currency          string   `json:"currency"`
	SalaryExpectation float64  `json:"salaryExpectation"`
	RemoteOnly        bool     `json:"remoteOnly"`
}

type ProfileSearchResponseDTO struct {
	ProfileId        string   `json:"profileId"`
	UserName         string   `json:"username"`
	UserProfileImage string   `json:"userProfileImage"`
	Headline         string   `json:"headline"`
	Seniority        int      `json:"seniority"`
	Skills           []string `json:"skills"`
	Location         int      `json:"location"`
}

type ProfileSearchResponse struct {
	TotalHits int64
	Hits      []ProfileSearchResponseDTO
}

type SearchResponse struct {
	ProfileId string
	UserId    string
}

type SearchService interface {
	ConfigureIndex() error
	IndexProfile(profile ProfileSearchDTO) error
	DeleteProfile(id string) error
	SearchProfiles(query *ProfileSearchQueryBuilder) (ProfileSearchResponse, error)
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
	// Cria o índice com primary key definida (se não existir, será criado)
	_, err := s.client.CreateIndex(&meilisearch.IndexConfig{
		Uid:        s.indexName,
		PrimaryKey: "profileId",
	})
	if err != nil {
		// Ignora erro se o índice já existir
		log.Printf("Índice pode já existir: %v", err)
	}

	index := s.client.Index(s.indexName)

	// Define Atributos Filtráveis (Facets / Filtros)
	// Isso permite queries como: "skills = 'Go' AND remote_only = true"
	filterableAttributes := []interface{}{
		"skills",
		"seniority",
		"yearsOfExperience",
		"salaryExpectation",
		"location",
		"remoteOnly",
		"openToWork",
		"contractType",
	}
	_, err = index.UpdateFilterableAttributes(&filterableAttributes)
	if err != nil {
		log.Printf("Erro ao configurar filtros do Meili: %v", err)
		return err
	}

	sortableAttributes := []string{
		"yearsOfExperience",
		"salaryExpectation",
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
		log.Printf("Erro ao enviar perfil %s para indexação: %v", profile.ProfileId, err)
		return err
	}

	log.Printf("Perfil %s enviado para fila de indexação. Task ID: %d", profile.ProfileId, task.TaskUID)
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

func (s *meiliService) SearchProfiles(query *ProfileSearchQueryBuilder) (ProfileSearchResponse, error) {
	searchQuery := ""
	var filter string

	if query != nil {
		searchQuery = query.BuildQuery()
		filter = query.BuildFilter()
	}

	searchRequest := &meilisearch.SearchRequest{
		AttributesToRetrieve: []string{"profileId", "username", "userProfileImage", "headline", "seniority", "skills", "location"},
		Limit:                1000,
	}

	if filter != "" {
		searchRequest.Filter = filter
	}

	searchRes, err := s.client.Index(s.indexName).Search(searchQuery, searchRequest)
	if err != nil {
		log.Printf("Erro ao buscar perfis: %v", err)
		return ProfileSearchResponse{}, err
	}

	var results []ProfileSearchResponseDTO
	if err := searchRes.Hits.DecodeInto(&results); err != nil {
		log.Printf("Erro ao decodificar resultados: %v", err)
		return ProfileSearchResponse{}, err
	}

	hitCount := searchRes.EstimatedTotalHits
	response := ProfileSearchResponse{
		TotalHits: hitCount,
		Hits:      results,
	}

	log.Printf("Busca retornou %d perfis", len(results))
	return response, nil
}
