package web

import (
	"fmt"
	"net/http"
	"portfolio/internal/portfolio"
	"portfolio/internal/search"
	"strings"
)

type ProfileSearchRequest struct {
	KeyWords     *[]string                 `json:"key_words,omitempty"`
	Skills       *[]string                 `json:"skills,omitempty"`
	Seniority    *[]portfolio.Seniority    `json:"seniority,omitempty"`
	Location     *[]portfolio.LocationType `json:"location,omitempty"`
	RemoteOnly   *bool                     `json:"remote_only,omitempty"`
	OpenToWork   *bool                     `json:"open_to_work,omitempty"`
	ContractType *[]string                 `json:"contract_type,omitempty"`
	MinYearsExp  *int                      `json:"min_years_of_experience,omitempty"`
	MaxYearsExp  *int                      `json:"max_years_of_experience,omitempty"`
	MinSalary    *float64                  `json:"min_salary,omitempty"`
	MaxSalary    *float64                  `json:"max_salary,omitempty"`
}

func (p *ProfileSearchRequest) ToProfileBuilder() *search.ProfileSearchQueryBuilder {
	builder := search.NewProfileSearchQueryBuilder()
	if p.KeyWords != nil {
		builder.WithKeyWords(*p.KeyWords...)
	}
	if p.Skills != nil {
		builder.WithSkills(*p.Skills...)
	}
	if p.Seniority != nil {
		var seniorityInts []int
		for _, s := range *p.Seniority {
			seniorityInts = append(seniorityInts, s.Int())
		}
		builder.WithSeniority(seniorityInts...)
	}
	if p.Location != nil {
		var locationInts []int
		for _, l := range *p.Location {
			locationInts = append(locationInts, l.Int())
		}
		builder.WithLocation(locationInts...)
	}
	if p.RemoteOnly != nil {
		builder.WithRemoteOnly(*p.RemoteOnly)
	}
	if p.OpenToWork != nil {
		builder.WithOpenToWork(*p.OpenToWork)
	}
	if p.ContractType != nil {
		builder.WithContractType(*p.ContractType...)
	}
	if p.MinYearsExp != nil {
		builder.WithMinYearsOfExperience(p.MinYearsExp)
	}
	if p.MaxYearsExp != nil {
		builder.WithMaxYearsOfExperience(p.MaxYearsExp)
	}
	if p.MinSalary != nil {
		builder.WithMinSalaryRange(p.MinSalary)
	}
	if p.MaxSalary != nil {
		builder.WithMaxSalaryRange(p.MaxSalary)
	}
	return builder
}

func extractSearchForm(r *http.Request) search.ProfileSearchQueryBuilder {
	var searchDto ProfileSearchRequest
	err := r.ParseForm()
	if err != nil {
		return *searchDto.ToProfileBuilder()
	}
	if keywords, exists := r.Form["keywords"]; exists && len(keywords) > 0 && keywords[0] != "" {
		kwList := strings.Split(keywords[0], ",")
		searchDto.KeyWords = &kwList
	}
	if skills, exists := r.Form["skills"]; exists && len(skills) > 0 && skills[0] != "" {
		skillsList := strings.Split(skills[0], ",")
		searchDto.Skills = &skillsList
	}
	// Seniority - múltiplos checkboxes vêm como array
	if seniority, exists := r.Form["seniority"]; exists && len(seniority) > 0 {
		var seniorityList []portfolio.Seniority
		for _, s := range seniority {
			if s != "" {
				seniorityList = append(seniorityList, portfolio.Seniority(s))
			}
		}
		if len(seniorityList) > 0 {
			searchDto.Seniority = &seniorityList
		}
	}
	// Location - múltiplos checkboxes vêm como array
	if location, exists := r.Form["location"]; exists && len(location) > 0 {
		var locationList []portfolio.LocationType
		for _, l := range location {
			if l != "" {
				locationList = append(locationList, portfolio.LocationType(l))
			}
		}
		if len(locationList) > 0 {
			searchDto.Location = &locationList
		}
	}
	if remoteOnly, exists := r.Form["remote_only"]; exists && len(remoteOnly) > 0 {
		ro := remoteOnly[0] == "true"
		searchDto.RemoteOnly = &ro
	}
	if openToWork, exists := r.Form["open_to_work"]; exists && len(openToWork) > 0 {
		otw := openToWork[0] == "true"
		searchDto.OpenToWork = &otw
	}
	// Contract type - múltiplos checkboxes vêm como array
	if contractType, exists := r.Form["contract_type"]; exists && len(contractType) > 0 {
		var ctList []string
		for _, ct := range contractType {
			if ct != "" {
				ctList = append(ctList, ct)
			}
		}
		if len(ctList) > 0 {
			searchDto.ContractType = &ctList
		}
	}
	if minYearsExp, exists := r.Form["min_years_of_experience"]; exists && len(minYearsExp) > 0 && minYearsExp[0] != "" {
		var minYE int
		_, err := fmt.Sscanf(minYearsExp[0], "%d", &minYE)
		if err == nil {
			searchDto.MinYearsExp = &minYE
		}
	}
	if maxYearsExp, exists := r.Form["max_years_of_experience"]; exists && len(maxYearsExp) > 0 && maxYearsExp[0] != "" {
		var maxYE int
		_, err := fmt.Sscanf(maxYearsExp[0], "%d", &maxYE)
		if err == nil {
			searchDto.MaxYearsExp = &maxYE
		}
	}
	if minSalary, exists := r.Form["min_salary"]; exists && len(minSalary) > 0 && minSalary[0] != "" {
		var minS float64
		_, err := fmt.Sscanf(minSalary[0], "%f", &minS)
		if err == nil {
			searchDto.MinSalary = &minS
		}
	}
	if maxSalary, exists := r.Form["max_salary"]; exists && len(maxSalary) > 0 && maxSalary[0] != "" {
		var maxS float64
		_, err := fmt.Sscanf(maxSalary[0], "%f", &maxS)
		if err == nil {
			searchDto.MaxSalary = &maxS
		}
	}

	return *searchDto.ToProfileBuilder()
}
