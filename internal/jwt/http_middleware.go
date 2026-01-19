package jwt

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

const AutenticatedUserKey string = "autenticatedUser"

type AutenticatedUser struct {
	ID              string
	Email           string
	FirstName       string
	LastName        string
	ProfileImageURL *string
}

func (service *JWTService) RequiredAutenticationMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		autenticatedUser := GetAutenticatedUserFromRequest(r, service)
		if autenticatedUser == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		log.Printf("[RequiredAutenticationMiddleware] User %s authenticated successfully", autenticatedUser.ID)
		ctx := context.WithValue(r.Context(), AutenticatedUserKey, *autenticatedUser)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (service *JWTService) OptionalAutenticationMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		autenticatedUser := GetAutenticatedUserFromRequest(r, service)
		if autenticatedUser == nil {
			next.ServeHTTP(w, r)
			return
		}
		ctx := context.WithValue(r.Context(), AutenticatedUserKey, *autenticatedUser)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}


func GetUserCurrentUser(ctx context.Context) *AutenticatedUser {
	user, _ := ctx.Value(AutenticatedUserKey).(AutenticatedUser)
	return &user
}

func GetUserIDFromToken(tokenString string, service *JWTService) (*string, error) {
	token, err := parseToken(tokenString, service.secretKey)
	if err != nil {
		return nil, err
	}
	return getClaimAsString(*token, "sub")
}

func GetUserEmailFromToken(tokenString string, service *JWTService) (*string, error) {
	token, err := parseToken(tokenString, service.secretKey)
	if err != nil {
		return nil, err
	}
	return getClaimAsString(*token, "email")
}

func GetUserNameFromToken(tokenString string, service *JWTService) (*string, error) {
	token, err := parseToken(tokenString, service.secretKey)
	if err != nil {
		return nil, err
	}
	return getClaimAsString(*token, "name")
}

func GetUserProfileFromToken(tokenString string, service *JWTService) (*string, error) {
	token, err := parseToken(tokenString, service.secretKey)
	if err != nil {
		return nil, err
	}
	return getClaimAsString(*token, "profileImageURL")
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

func parseToken(tokenString string, secretKey []byte) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return secretKey, nil
	})
	return token, err
}

func GetJwtTokenFromRequest(r *http.Request) string {

	cookie, err := r.Cookie("access_token")
	if err == nil {
		log.Printf("Found JWT token in cookie")
		return cookie.Value
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	tokenString := parts[1]
	return tokenString
}

func GetAutenticatedUserFromRequest(r *http.Request, jwtService *JWTService) *AutenticatedUser {
	token := GetJwtTokenFromRequest(r)
	
	if token == "" {
		return nil
	}

	userID, err := GetUserIDFromToken(token, jwtService)
	if err != nil || userID == nil || *userID == "" {
		return nil
	}

	userEmail, emailErr := GetUserEmailFromToken(token, jwtService)
	if emailErr != nil || userEmail == nil || *userEmail == "" {
		return nil
	}
	name, nameErr := GetUserNameFromToken(token, jwtService)
	if nameErr != nil || name == nil || *name == "" {
		return nil
	}
	profileImageURL, profileErr := GetUserProfileFromToken(token, jwtService)
	if profileErr != nil {
		return nil
	}
	if *profileImageURL == "" {
		profileImageURL = nil
	}

	// Split name into FirstName and LastName
	firstName := *name
	lastName := ""
	if parts := strings.SplitN(*name, " ", 2); len(parts) == 2 {
		firstName = parts[0]
		lastName = parts[1]
	}

	autenticatedUser := AutenticatedUser{
		ID:              *userID,
		Email:           *userEmail,
		FirstName:       firstName,
		LastName:        lastName,
		ProfileImageURL: profileImageURL,
	}
	return &autenticatedUser
}
