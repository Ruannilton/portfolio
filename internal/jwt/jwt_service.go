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
}

type TokenService struct {
	secretKey []byte
	issuer    string
}

func NewJWTService(cfg *config.Config) *TokenService {
	return &TokenService{
		secretKey: []byte(cfg.JWTSecretKey),
		issuer:    cfg.JWTIssuer,
	}
}

func (s *TokenService) GenerateToken(input *GenerateTokenInput) (*TokenResponse, error) {
	accessDuration := time.Minute * 15
	refreshDuration := time.Hour * 24 * 7 // 7 dias

	now := time.Now()
	accessExp := now.Add(accessDuration)
	refreshExp := now.Add(refreshDuration)

	accessClaims := jwt.MapClaims{
		"sub":   input.UserID,
		"email": input.UserEmail,
		"type":  "access",
		"iss":   s.issuer,
		"exp":   accessExp.Unix(),
	}
	accessTokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(s.secretKey)
	if err != nil {
		return nil, err
	}

	refreshClaims := jwt.MapClaims{
		"sub":  input.UserID,
		"type": "refresh",
		"iss":  s.issuer,
		"exp":  refreshExp.Unix(),
	}
	refreshTokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(s.secretKey)
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

func (s *TokenService) GetUserIDFromToken(tokenString string) (*string, error) {
	token, err := parseToken(tokenString, s)
	if err != nil {
		return nil, err
	}
	return getClaimAsString(*token, "sub")
}

func (s *TokenService) GetUserEmailFromToken(tokenString string) (*string, error) {
	token, err := parseToken(tokenString, s)
	if err != nil {
		return nil, err
	}
	return getClaimAsString(*token, "email")
}

func getClaimAsString(token jwt.Token, key string) (*string, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	claim, ok := claims[key].(string)
	if !ok {
		return nil, jwt.ErrTokenInvalidSubject
	}

	return &claim, nil
}

func parseToken(tokenString string, s *TokenService) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return s.secretKey, nil
	})
	return token, err
}
