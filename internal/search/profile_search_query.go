package search

import (
	"fmt"
	"strings"
)

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
func (b *ProfileSearchQueryBuilder) WithMinYearsOfExperience(min *int) *ProfileSearchQueryBuilder {

	b.minYearsExp = min
	return b
}

func (b *ProfileSearchQueryBuilder) WithMaxYearsOfExperience(max *int) *ProfileSearchQueryBuilder {

	b.maxYearsExp = max
	return b
}

// WithSalaryRange sets the salary range filter
func (b *ProfileSearchQueryBuilder) WithMinSalaryRange(min *float64) *ProfileSearchQueryBuilder {
	b.minSalary = min
	return b
}

func (b *ProfileSearchQueryBuilder) WithMaxSalaryRange(max *float64) *ProfileSearchQueryBuilder {
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
		filters = append(filters, fmt.Sprintf("remoteOnly = %t", *b.remoteOnly))
	}

	// Open To Work
	if b.openToWork != nil {
		filters = append(filters, fmt.Sprintf("openToWork = %t", *b.openToWork))
	}

	// Contract Type (OR between values)
	if len(b.contractType) > 0 {
		ctFilters := make([]string, len(b.contractType))
		for i, ct := range b.contractType {
			ctFilters[i] = fmt.Sprintf("contractType = '%s'", ct)
		}
		filters = append(filters, "("+strings.Join(ctFilters, " OR ")+")")
	}

	// Years of Experience (range)
	if b.minYearsExp != nil && b.maxYearsExp != nil {
		filters = append(filters, fmt.Sprintf("yearsOfExperience %d TO %d", *b.minYearsExp, *b.maxYearsExp))
	} else if b.minYearsExp != nil {
		filters = append(filters, fmt.Sprintf("yearsOfExperience >= %d", *b.minYearsExp))
	} else if b.maxYearsExp != nil {
		filters = append(filters, fmt.Sprintf("yearsOfExperience <= %d", *b.maxYearsExp))
	}

	// Salary Range
	if b.minSalary != nil && b.maxSalary != nil {
		filters = append(filters, fmt.Sprintf("salaryExpectation %.2f TO %.2f", *b.minSalary, *b.maxSalary))
	} else if b.minSalary != nil {
		filters = append(filters, fmt.Sprintf("salaryExpectation >= %.2f", *b.minSalary))
	} else if b.maxSalary != nil {
		filters = append(filters, fmt.Sprintf("salaryExpectation <= %.2f", *b.maxSalary))
	}

	return strings.Join(filters, " AND ")
}
