package auth

import (
	"context"
	"database/sql"
	"errors"
	"strings"
)

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	Find(ctx context.Context, userID string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByResetToken(ctx context.Context, token string) (*User, error)
	FindByProviderID(ctx context.Context, provider, providerID string) (*User, error)
	Save(ctx context.Context, user *User) error
}

type userRepo struct {
	db *sql.DB // postgres database connection
}

// Create implements [UserRepository].
func (u *userRepo) Create(ctx context.Context, user *User) error {
	user.Email = strings.ToLower(user.Email)
	query := `
		INSERT INTO users (id, first_name, last_name, email, password_hash, provider, provider_id, reset_token, created_at, profile_image, github_access_token)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := u.db.ExecContext(ctx, query,
		user.ID,
		user.FirstName,
		user.LastName,
		user.Email,
		user.PasswordHash,
		user.Provider,
		user.ProviderID,
		user.ResetToken,
		user.CreatedAt,
		user.ProfileImage,
		user.GithubAcessToken,
	)
	return err
}

// FindByEmail implements [UserRepository].
func (u *userRepo) FindByEmail(ctx context.Context, email string) (*User, error) {
	email = strings.ToLower(email)
	query := `
		SELECT id, first_name, last_name, email, password_hash, provider, provider_id, reset_token, created_at, profile_image, github_access_token
		FROM users
		WHERE email = $1
		LIMIT 1
	`
	user := &User{}
	err := u.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.PasswordHash,
		&user.Provider,
		&user.ProviderID,
		&user.ResetToken,
		&user.CreatedAt,
		&user.ProfileImage,
		&user.GithubAcessToken,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (u *userRepo) Find(ctx context.Context, id string) (*User, error) {
	
	query := `
		SELECT id, first_name, last_name, email, password_hash, provider, provider_id, reset_token, created_at, profile_image, github_access_token
		FROM users
		WHERE id = $1
		LIMIT 1
	`
	user := &User{}
	err := u.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.PasswordHash,
		&user.Provider,
		&user.ProviderID,
		&user.ResetToken,
		&user.CreatedAt,
		&user.ProfileImage,
		&user.GithubAcessToken,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

// Save implements [UserRepository].
func (u *userRepo) Save(ctx context.Context, user *User) error {
	user.Email = strings.ToLower(user.Email)
	query := `
		UPDATE users
		SET first_name = $1, last_name = $2, email = $3, password_hash = $4, 
		    provider = $5, provider_id = $6, reset_token = $7, profile_image = $8, github_access_token = $9
		WHERE id = $10
	`
	result, err := u.db.ExecContext(ctx, query,
		user.FirstName,
		user.LastName,
		user.Email,
		user.PasswordHash,
		user.Provider,
		user.ProviderID,
		user.ResetToken,
		user.ProfileImage,
		user.GithubAcessToken,
		user.ID,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (u *userRepo) FindByResetToken(ctx context.Context, token string) (*User, error) {
	query := `
		SELECT id, first_name, last_name, email, password_hash, provider, provider_id, reset_token, created_at, profile_image, github_access_token
		FROM users
		WHERE reset_token = $1
		LIMIT 1
	`
	user := &User{}
	err := u.db.QueryRowContext(ctx, query, token).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.PasswordHash,
		&user.Provider,
		&user.ProviderID,
		&user.ResetToken,
		&user.CreatedAt,
		&user.ProfileImage,
		&user.GithubAcessToken,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

// FindByProviderID implements [UserRepository].
func (u *userRepo) FindByProviderID(ctx context.Context, provider, providerID string) (*User, error) {
	query := `
		SELECT id, first_name, last_name, email, password_hash, provider, provider_id, reset_token, created_at, profile_image, github_access_token
		FROM users
		WHERE provider = $1 AND provider_id = $2
		LIMIT 1
	`
	user := &User{}
	err := u.db.QueryRowContext(ctx, query, provider, providerID).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.PasswordHash,
		&user.Provider,
		&user.ProviderID,
		&user.ResetToken,
		&user.CreatedAt,
		&user.ProfileImage,
		&user.GithubAcessToken,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepo{db: db}
}
