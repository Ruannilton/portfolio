package search

import (
	"log"
	"portfolio/internal/config"

	"github.com/meilisearch/meilisearch-go"
)

type ProfileSearchDTO struct {
	ID                string   `json:"id"`
	Headline          string   `json:"headline"`
	Bio               string   `json:"bio"`
	Skills            []string `json:"skills"`
	Seniority         string   `json:"seniority"`
	YearsOfExp        int      `json:"years_of_experience"`
	Location          string   `json:"location"`
	RemoteOnly        bool     `json:"remote_only"`
	OpenToWork        bool     `json:"open_to_work"`
	ContractType      string   `json:"contract_type"`
	Currency          string   `json:"currency"`
	SalaryExpectation float64  `json:"salary_expectation"`
	ProjectTags       []string `json:"project_tags"`
}

type SearchService interface {
	ConfigureIndex() error
	IndexProfile(profile ProfileSearchDTO) error
	DeleteProfile(id string) error
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

	// 2. Define Atributos Filtráveis (Facets / Filtros)
	// Isso permite queries como: "skills = 'Go' AND remote_only = true"
	filterableAttributes := []interface{}{
		"id",
		"skills",
		"seniority",
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
