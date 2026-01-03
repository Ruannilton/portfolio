package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"portfolio/internal/jwt"

	"github.com/markbates/goth/gothic"
)

type AuthModule struct {
	authService  *AuthService
	tokenService *jwt.TokenService
}

func NewAuthModule(authService *AuthService, tokenService *jwt.TokenService) *AuthModule {
	return &AuthModule{
		authService:  authService,
		tokenService: tokenService,
	}
}

func (module *AuthModule) RegisterAuthRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /auth/{provider}", module.beginOAuthHandler)
	mux.HandleFunc("GET /auth/{provider}/callback", module.oAuthCallbackHandler)
	mux.HandleFunc("GET /auth/logout", module.logoutHandler)
	mux.HandleFunc("POST /auth/register", module.registerLocalUser)
	mux.HandleFunc("POST /auth/login", module.loginLocalUser)
	mux.HandleFunc("POST /auth/forgot-password", module.forgotPassword)
	mux.HandleFunc("POST /auth/reset-password", module.resetPassword)

	mux.HandleFunc("GET /auth/me", jwt.AuthMiddleware(module.tokenService, module.me))

	return mux
}

func (module *AuthModule) beginOAuthHandler(w http.ResponseWriter, r *http.Request) {
	provider := r.PathValue("provider")
	// Adiciona provider como query param, que é como o Gothic espera
	q := r.URL.Query()
	q.Add("provider", provider)
	r.URL.RawQuery = q.Encode()

	gothic.BeginAuthHandler(w, r)
}

func (module *AuthModule) oAuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	provider := r.PathValue("provider")
	// Adiciona provider como query param, que é como o Gothic espera
	q := r.URL.Query()
	q.Add("provider", provider)
	r.URL.RawQuery = q.Encode()

	// Completa autenticação OAuth e obtém dados do usuário
	gothUser, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		log.Printf("OAuth error: %v", err)
		http.Error(w, "OAuth authentication failed", http.StatusUnauthorized)
		return
	}

	// Processa login/registro e gera tokens
	tokenResponse, err := module.authService.CompleteOAuthLogin(r.Context(), gothUser)
	if err != nil {
		log.Printf("CompleteOAuthLogin error: %v", err)
		http.Error(w, "Failed to complete authentication", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tokenResponse); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

func (module *AuthModule) logoutHandler(w http.ResponseWriter, r *http.Request) {
	if err := module.authService.Logout(w, r); err != nil {
		log.Printf("Logout error: %v", err)
		http.Error(w, "Failed to logout", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	resp := map[string]string{"message": "Successfully logged out"}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

func (module *AuthModule) registerLocalUser(w http.ResponseWriter, r *http.Request) {
	var request RegisterLocalUserInput

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := module.authService.RegisterLocalUser(r.Context(), request)
	if err != nil {
		log.Printf("RegisterLocalUser error: %v", err)
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (module *AuthModule) loginLocalUser(w http.ResponseWriter, r *http.Request) {
	var request LoginInput
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	tokenResponse, err := module.authService.LoginLocal(r.Context(), request)
	if err != nil {
		log.Printf("LoginLocalUser error: %v", err)
		http.Error(w, "Failed to login user", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tokenResponse); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func (module *AuthModule) forgotPassword(w http.ResponseWriter, r *http.Request) {
	var request ForgotPasswordInput
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	err := module.authService.ForgotPassword(r.Context(), request)
	if err != nil {
		log.Printf("ForgotPassword error: %v", err)
		http.Error(w, "Failed to process forgot password", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (module *AuthModule) resetPassword(w http.ResponseWriter, r *http.Request) {
	var request ResetPasswordInput
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	err := module.authService.ResetPassword(r.Context(), request)
	if err != nil {
		log.Printf("ResetPassword error: %v", err)
		http.Error(w, "Failed to reset password", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (module *AuthModule) me(w http.ResponseWriter, r *http.Request) {
	user, err := module.authService.GetUserFromContext(r.Context())

	if err != nil {
		log.Printf("GetUserFromContext error: %v", err)
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	response := &GetUserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
