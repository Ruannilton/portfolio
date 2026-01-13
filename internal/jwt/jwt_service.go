package jwt

import (
	"portfolio/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // Segundos até expirar
	TokenType    string `json:"token_type"` // Geralmente "Bearer"
}

type GenerateTokenInput struct {
	UserID    string
	UserEmail string
	UerName  string
	ProfileImageURL string
}

type JWTService struct {
		secretKey []byte
	issuer    string
}

var (

)

func NewJWTService(cfg *config.Config) JWTService {
	return JWTService{
	secretKey  : []byte(cfg.JWTSecretKey),
	issuer : cfg.JWTIssuer,
	}
}

func (service *JWTService) GenerateToken(input *GenerateTokenInput) (*TokenResponse, error) {
	accessDuration := time.Minute * 15
	refreshDuration := time.Hour * 24 * 7 // 7 dias

	now := time.Now()
	accessExp := now.Add(accessDuration)
	refreshExp := now.Add(refreshDuration)

	accessClaims := jwt.MapClaims{
		"sub":   input.UserID,
		"email": input.UserEmail,
		"type":  "access",
		"iss":   service.issuer,
		"exp":   accessExp.Unix(),
		"name":  input.UerName,
		"profileImageURL": input.ProfileImageURL,
	}
	accessTokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(service.secretKey)
	if err != nil {
		return nil, err
	}

	refreshClaims := jwt.MapClaims{
		"sub":  input.UserID,
		"type": "refresh",
		"iss":  service.issuer,
		"exp":  refreshExp.Unix(),
		"name":  input.UerName,
		"profileImageURL": input.ProfileImageURL,
	}
	refreshTokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(service.secretKey)
	if err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresIn:    int64(accessDuration.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

