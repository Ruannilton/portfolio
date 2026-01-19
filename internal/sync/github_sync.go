package sync

import (
	"context"
	"fmt"
	"log"
	"sort"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type LanguageNode struct {
	Size int `graphql:"size"`
	Node struct {
		Name  string `graphql:"name"`
		Color string  `graphql:"color"`
	} `graphql:"node"`
}

type Repository struct {
	Name            string `graphql:"name"`
	DatabaseId      int `graphql:"databaseId"`
	Description     string `graphql:"description"`
	Url             string `graphql:"url"`
	StargazerCount int `graphql:"stargazerCount"`
	IsFork          bool `graphql:"isFork"`
	Languages       struct {
		Edges []LanguageNode  `graphql:"edges"`
	} `graphql:"languages(first: 5, orderBy: {field: SIZE, direction: DESC})"`
	RepositoryTopics struct {
		Nodes []struct {
			Topic struct {
				Name string  `graphql:"name"`
			}  `graphql:"topic"`
		} `graphql:"nodes"`
	} `graphql:"repositoryTopics(first: 5)"`
}

type GithubQuery struct {
	Viewer struct {
		Login        string `graphql:"login"`
		Bio          string `graphql:"bio"`
		AvatarUrl    string `graphql:"avatarUrl"`
		Company      string `graphql:"company"`
		Location     string `graphql:"location"`
		Repositories struct {
			Nodes []Repository
		} `graphql:"repositories(first: 100, privacy: PUBLIC, isFork: false, ownerAffiliations: OWNER, orderBy: {field: PUSHED_AT, direction: DESC})"`
	}
}

type TechRadarStats struct {
	Language   string
	Color      string
	Percentage float64
	TotalBytes int
}

type GithubRepository struct {
	Name        string
	Description string
	Url         string
	Languages   []string
	Topics      []string
	ProviderId string
}

type GithubProfile struct {
	Bio          string `json:"bio"`
	TechRadar    []TechRadarStats `json:"techRadar"`
	Repositories []GithubRepository `json:"repositories"`
}

func SyncGithubData(ctx context.Context, token string) (*GithubProfile, error) {

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(ctx, src)

	client := githubv4.NewClient(httpClient)

	var query GithubQuery

	// Executa a Query
	err := client.Query(ctx, &query, nil)
	if err != nil {
		log.Fatalf("Erro na query GraphQL: %v", err)
		return nil, err
	}

	techRadar := calculateTechRadar(query)
	repositories := extractRepositories(query)
	profile := &GithubProfile{
		Bio:          query.Viewer.Bio,
		TechRadar:    techRadar,
		Repositories: repositories,
	}
	return profile, nil
}

func calculateTechRadar(query GithubQuery) []TechRadarStats {
	langByteCount := make(map[string]int)
	langColors := make(map[string]string)
	totalBytesAllLangs := 0

	for _, repo := range query.Viewer.Repositories.Nodes {

		for _, lang := range repo.Languages.Edges {
			name := lang.Node.Name
			size := lang.Size
			color := lang.Node.Color

			langByteCount[name] += size
			langColors[name] = color
			totalBytesAllLangs += size
		}
	}

	var radar []TechRadarStats

	for name, size := range langByteCount {
		percentage := (float64(size) / float64(totalBytesAllLangs)) * 100

		if percentage < 1.0 {
			continue
		}

		radar = append(radar, TechRadarStats{
			Language:   name,
			Color:      langColors[name],
			Percentage: percentage,
			TotalBytes: size,
		})
	}

	// Ordena do mais usado para o menos usado
	sort.Slice(radar, func(i, j int) bool {
		return radar[i].Percentage > radar[j].Percentage
	})
	return radar
}

func extractRepositories(query GithubQuery) []GithubRepository {
	var repositories []GithubRepository

	for _, repo := range query.Viewer.Repositories.Nodes {
		var languages []string
		var topics []string

		for _, lang := range repo.Languages.Edges {
			languages = append(languages, lang.Node.Name)
		}

		for _, topic := range repo.RepositoryTopics.Nodes {
			topics = append(topics, topic.Topic.Name)
		}
		repositories = append(repositories, GithubRepository{
			Name:        repo.Name,
			Description: repo.Description,
			Url:         repo.Url,
			Languages:   languages,
			Topics:      topics,
			ProviderId:  fmt.Sprintf("%d", repo.DatabaseId),
		})
	}
	return repositories
}
