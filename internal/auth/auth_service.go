package auth

import (
	"context"
	"errors"
	"log"
	"net/http"
	"portfolio/internal/config"
	"portfolio/internal/jwt"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
)

type AuthService struct {
	repo       UserRepository
	jwtService *jwt.JWTService
}

type RegisterLocalUserInput struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ForgotPasswordInput struct {
	Email string `json:"email"`
}

type ResetPasswordInput struct {
	Token       string `json:"token"`
	NewPassword string `json:"newPassword"`
}

type GetUserResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

func NewAuthService(cfg *config.Config, repo UserRepository, jwtService *jwt.JWTService) *AuthService {
	// Config already carregada em `config.LoadConfig()` e variáveis de ambiente
	// são fornecidas pelo Docker via `env_file`; não devemos panicar se não
	// existir um arquivo .env no filesystem.

	googleClientID := cfg.GoogleClientID
	googleClientSecret := cfg.GoogleClientSecret

	githubClientID := cfg.GithubClientID
	githubClientSecret := cfg.GithubClientSecret

	key := cfg.SessionKey
	maxAge := cfg.SessionMaxAge
	isProduction := cfg.IsProduction


	// Garante que a chave tenha 32 bytes para AES-256
	keyBytes := []byte(key)
	if len(keyBytes) < 32 {
		padded := make([]byte, 32)
		copy(padded, keyBytes)
		keyBytes = padded
	} else if len(keyBytes) > 32 {
		keyBytes = keyBytes[:32]
	}

	store := sessions.NewCookieStore(keyBytes)
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = isProduction
	store.Options.SameSite = http.SameSiteLaxMode

	gothic.Store = store
	redirectUrl := cfg.AppRedirectURL

	log.Println("Github client: ", githubClientID)
	log.Println("Github secret: ", githubClientSecret)
	log.Println("Redirect URL: ", redirectUrl+"/auth/github/callback")

	goth.UseProviders(
		google.New(googleClientID, googleClientSecret, redirectUrl+"/auth/google/callback", "email", "profile"),
		github.New(githubClientID, githubClientSecret, redirectUrl+"/auth/github/callback", "read:user", "user:email"),
	)

	return &AuthService{
		repo:       repo,
		jwtService: jwtService,
	}
}

func (s *AuthService) RegisterLocalUser(ctx context.Context, input RegisterLocalUserInput) error {
	existent, err := s.repo.FindByEmail(ctx, input.Email)

	if err != nil && err != ErrUserNotFound {
		return err
	}

	if existent != nil {
		return ErrEmailAlreadyInUse
	}

	user, createErr := NewLocalUser(input.FirstName, input.LastName, input.Email, input.Password, nil)

	if createErr != nil {
		return createErr
	}

	if repoErr := s.repo.Create(ctx, user); repoErr != nil {
		return repoErr
	}

	return nil
}

func (uc *AuthService) LoginLocal(ctx context.Context, input LoginInput) (*jwt.TokenResponse, error) {
	user, err := uc.repo.FindByEmail(ctx, input.Email)
	if err != nil || user == nil {
		return nil, errors.New("invalid credentials")
	}

	if !user.ValidatePassword(input.Password) {
		return nil, errors.New("invalid credentials")
	}

	profImageUrl := ""
	if user.ProfileImage != nil {
		profImageUrl = *user.ProfileImage
	}

	inputToken := &jwt.GenerateTokenInput{
		UserID:          user.ID,
		UserEmail:       user.Email,
		UerName:         user.FirstName + " " + user.LastName,
		ProfileImageURL: profImageUrl,
	}
	token, err := uc.jwtService.GenerateToken(inputToken)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (uc *AuthService) ForgotPassword(ctx context.Context, input ForgotPasswordInput) error {
	user, err := uc.repo.FindByEmail(ctx, input.Email)
	if user == nil || err != nil {
		return nil
	}

	token := uuid.New().String()
	user.ResetToken = &token

	if err := uc.repo.Save(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) ResetPassword(ctx context.Context, input ResetPasswordInput) error {
	user, err := s.repo.FindByResetToken(ctx, input.Token)
	if err != nil || user == nil {
		return errors.New("invalid or expired reset token")
	}

	if err := user.SetPassword(input.NewPassword); err != nil {
		return err
	}

	if err := s.repo.Save(ctx, user); err != nil {
		return err
	}
	return nil
}

func (s *AuthService) CompleteOAuthLogin(ctx context.Context, gothUser goth.User) (*jwt.TokenResponse, error) {

	// GitHub pode retornar apenas o Name. Se FirstName estiver vazio, tentamos extrair do Name.
	if gothUser.FirstName == "" && gothUser.Name != "" {
		parts := strings.SplitN(gothUser.Name, " ", 2)
		gothUser.FirstName = parts[0]
		if len(parts) > 1 {
			gothUser.LastName = parts[1]
		}
	}

	user, err := s.repo.FindByEmail(ctx, gothUser.Email)

	if err != nil && !errors.Is(err, ErrUserNotFound) {
		return nil, err
	}

	if user != nil {
		user.Provider = gothUser.Provider
		user.ProviderID = &gothUser.UserID
		user.ProfileImage = &gothUser.AvatarURL

		if gothUser.Provider == "github" {
             // O token vem aqui. Recomendo criptografar antes de salvar (ver nota abaixo)
            user.SetGithubAccessToken(gothUser.AccessToken) 
        }

		if saveErr := s.repo.Save(ctx, user); saveErr != nil {
			return nil, saveErr
		}
	} else {

		user = NewProviderUser(
			gothUser.FirstName,
			gothUser.LastName,
			gothUser.Email,
			gothUser.Provider,
			gothUser.UserID,
			&gothUser.AvatarURL,
		)

		if gothUser.Provider == "github" {
             // O token vem aqui. Recomendo criptografar antes de salvar (ver nota abaixo)
            user.SetGithubAccessToken(gothUser.AccessToken) 
        }

		if createErr := s.repo.Create(ctx, user); createErr != nil {
			return nil, createErr
		}
	}

	userProfile := ""
	if user.ProfileImage != nil {
		userProfile = *user.ProfileImage
	}

	inputToken := &jwt.GenerateTokenInput{
		UserID:          user.ID,
		UserEmail:       user.Email,
		UerName:         user.FirstName + " " + user.LastName,
		ProfileImageURL: userProfile,
	}
	token, tokenErr := s.jwtService.GenerateToken(inputToken)
	if tokenErr != nil {
		return nil, tokenErr
	}

	return token, nil
}

func (s *AuthService) Logout(res http.ResponseWriter, req *http.Request) error {
	err := gothic.Logout(res, req)
	if err != nil {
		return err
	}
	return nil
}

func (s *AuthService) GetUserFromContext(ctx context.Context) (*User, error) {

	loggedUser := jwt.GetUserCurrentUser(ctx)

	user, err := s.repo.FindByEmail(ctx, loggedUser.Email)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByID busca um usuário pelo ID
func (s *AuthService) GetUserByID(ctx context.Context, userID string) (*User, error) {
	return s.repo.Find(ctx, userID)
}
