package jwt

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)


const AutenticatedUserKey string = "autenticatedUser"
type AutenticatedUser struct {
	UserID    string
	UserEmail string
}

func (service *JWTService) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		
		token := getJwtTokenFromRequest(r)
		if token == "" {
			http.Error(w, "Authorization token required", http.StatusUnauthorized)
			return
		}
	
		userID, err := GetUserIDFromToken(token,service)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		userEmail, err := GetUserEmailFromToken(token,service)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		autenticatedUser :=  AutenticatedUser{
			UserID:    *userID,
			UserEmail: *userEmail,
		}

		ctx := context.WithValue(r.Context(), AutenticatedUserKey, autenticatedUser)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func GetUserCurrentUser(ctx context.Context) (*AutenticatedUser) {
	user, _ := ctx.Value(AutenticatedUserKey).(AutenticatedUser)
	return &user
}

func GetUserIDFromToken(tokenString string, service *JWTService) (*string, error) {
	token, err := parseToken(tokenString,service.secretKey)
	if err != nil {
		return nil, err
	}
	return getClaimAsString(*token, "sub")
}

func  GetUserEmailFromToken(tokenString string,service *JWTService) (*string, error) {
	token, err := parseToken(tokenString,service.secretKey)
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

func parseToken(tokenString string, secretKey []byte) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return secretKey, nil
	})
	return token, err
}

func getJwtTokenFromRequest(r *http.Request) string {
		
		cookie, err := r.Cookie("access_token")
		if err == nil {
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
