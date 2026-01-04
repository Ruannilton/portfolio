package portfolio

import (
	"context"
	"database/sql"
	"errors"
)

type ProfileRepository interface {
	Create(ctx context.Context, profile *Profile) error
	Update(ctx context.Context, profile *Profile) error
	FindByUserID(ctx context.Context, userID string) (*Profile, error)
	Delete(ctx context.Context, userID string) error
}

type profileRepo struct {
	db *sql.DB
}

func NewProfileRepository(db *sql.DB) ProfileRepository {
	return &profileRepo{db: db}
}

func (r *profileRepo) Create(ctx context.Context, p *Profile) error {
	query := `
		INSERT INTO profiles (
			id, user_id, headline, bio, seniority, years_of_experience, open_to_work,
			salary_expectation, currency, contract_type, location, remote_only,
			skills, social_links, experiences, projects, educations, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
	`
	_, err := r.db.ExecContext(ctx, query,
		p.ID, p.UserID, p.Headline, p.Bio, p.Seniority, p.YearsOfExp, p.OpenToWork,
		p.SalaryExpectation, p.Currency, p.ContractType, p.Location, p.RemoteOnly,
		p.Skills, p.SocialLinks, p.Experiences, p.Projects, p.Educations, p.CreatedAt, p.UpdatedAt,
	)
	return err
}

func (r *profileRepo) Update(ctx context.Context, p *Profile) error {
	query := `
		UPDATE profiles SET
			headline=$1, bio=$2, seniority=$3, years_of_experience=$4, open_to_work=$5,
			salary_expectation=$6, currency=$7, contract_type=$8, location=$9, remote_only=$10,
			skills=$11, social_links=$12, experiences=$13, projects=$14, educations=$15, updated_at=$16
		WHERE user_id = $17
	`
	result, err := r.db.ExecContext(ctx, query,
		p.Headline, p.Bio, p.Seniority, p.YearsOfExp, p.OpenToWork,
		p.SalaryExpectation, p.Currency, p.ContractType, p.Location, p.RemoteOnly,
		p.Skills, p.SocialLinks, p.Experiences, p.Projects, p.Educations, p.UpdatedAt,
		p.UserID,
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrProfileNotFound
	}
	return nil
}

func (r *profileRepo) FindByUserID(ctx context.Context, userID string) (*Profile, error) {
	query := `
		SELECT id, user_id, headline, bio, seniority, years_of_experience, open_to_work,
		       salary_expectation, currency, contract_type, location, remote_only,
		       skills, social_links, experiences, projects, educations, created_at, updated_at
		FROM profiles WHERE user_id = $1 LIMIT 1
	`
	p := &Profile{}
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&p.ID, &p.UserID, &p.Headline, &p.Bio, &p.Seniority, &p.YearsOfExp, &p.OpenToWork,
		&p.SalaryExpectation, &p.Currency, &p.ContractType, &p.Location, &p.RemoteOnly,
		&p.Skills, &p.SocialLinks, &p.Experiences, &p.Projects, &p.Educations, &p.CreatedAt, &p.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrProfileNotFound
		}
		return nil, err
	}
	return p, nil
}

func (r *profileRepo) Delete(ctx context.Context, userID string) error {
	query := `DELETE FROM profiles WHERE user_id = $1`
	result, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrProfileNotFound
	}
	return nil
}