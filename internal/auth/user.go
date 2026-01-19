package auth

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           string
	FirstName    string
	LastName     string
	Email        string
	PasswordHash *string
	Provider     string
	ProviderID   *string
	ResetToken   *string
	CreatedAt    time.Time
	ProfileImage *string
	GithubAcessToken *string
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func NewLocalUser(firstName, lastName, email, password string, profileImage *string) (*User, error) {
	hash, err := hashPassword(password)

	if err != nil {
		return nil, ErrFailedToGeneratePasswordHash
	}

	passwordHash := string(hash)
	return &User{
		ID:           uuid.New().String(),
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		PasswordHash: &passwordHash,
		Provider:     "local",
		CreatedAt:    time.Now(),
		ProfileImage: profileImage,
	}, nil
}

func NewProviderUser(firstName, lastName, email, provider, providerID string, profileImage *string) *User {
	return &User{
		ID:         uuid.New().String(),
		FirstName:  firstName,
		LastName:   lastName,
		Email:      email,
		Provider:   provider,
		ProviderID: &providerID,
		CreatedAt:  time.Now(),
		ProfileImage: profileImage,
	}
}

func (u *User) ValidatePassword(password string) bool {
	if u.PasswordHash == nil {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(*u.PasswordHash), []byte(password))
	return err == nil
}

func (u *User) SetPassword(newPassword string) error {
	hash, err := hashPassword(newPassword)
	if err != nil {
		return ErrFailedToGeneratePasswordHash
	}
	passwordHash := string(hash)
	u.PasswordHash = &passwordHash
	u.ResetToken = nil
	return nil
}

func (u *User) SetGithubAccessToken(accessToken string) error {
	// hash, err := hashPassword(accessToken)
	// if err != nil {
	// 	return ErrFailedToGenerateAccessTokenHash
	// }
	// accessTokoenHash := string(hash)
	u.GithubAcessToken = &accessToken
	return nil
}

func (u *User)  GetGithubAccessToken() *string {
	return u.GithubAcessToken
}