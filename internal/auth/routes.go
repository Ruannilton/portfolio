package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"portfolio/internal/jwt"

	"github.com/gorilla/mux"
	"github.com/markbates/goth/gothic"
)

type AuthModule struct {
	authService  *AuthService
	jwtService *jwt.JWTService
}

func NewAuthModule(authService *AuthService, jwtService *jwt.JWTService ) *AuthModule {
	return &AuthModule{
		authService:  authService,
		jwtService:     jwtService,
	}
}

func (module *AuthModule) RegisterAuthRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/{provider}", module.beginOAuthHandler).Methods("GET")
	router.HandleFunc("/{provider}/callback", module.oAuthCallbackHandler).Methods("GET")
	router.HandleFunc("/logout", module.logoutHandler).Methods("GET")
	router.HandleFunc("/register", module.registerLocalUser).Methods("POST")
	router.HandleFunc("/login", module.loginLocalUser).Methods("POST")
	router.HandleFunc("/forgot-password", module.forgotPassword).Methods("POST")
	router.HandleFunc("/reset-password", module.resetPassword).Methods("POST")

	router.HandleFunc("/me", module.jwtService.AuthMiddleware(module.me)).Methods("GET")

	return router
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
		http.Redirect(w, r, "/login?error=oauth_failed", http.StatusFound)
		return
	}

	// Processa login/registro e gera tokens
	tokenResponse, err := module.authService.CompleteOAuthLogin(r.Context(), gothUser)
	if err != nil {
		log.Printf("CompleteOAuthLogin error: %v", err)
		http.Redirect(w, r, "/login?error=auth_failed", http.StatusFound)
		return
	}

	// Seta cookie com access_token
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    tokenResponse.AccessToken,
		Path:     "/",
		MaxAge:   int(tokenResponse.ExpiresIn),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	// Redireciona para a página do app
	http.Redirect(w, r, "/app", http.StatusFound)
}

func (module *AuthModule) logoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
        Name:     "access_token",
        Value:    "",
        Path:     "/",
        MaxAge:   -1,
        HttpOnly: true,
        SameSite: http.SameSiteLaxMode,
    })

	if err := module.authService.Logout(w, r); err != nil {
		log.Printf("Logout error: %v", err)
		http.Error(w, "Failed to logout", http.StatusInternalServerError)
		return
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
		if err == ErrEmailAlreadyInUse {
			http.Error(w, "Email já está em uso", http.StatusConflict)
			return
		}
		http.Error(w, "Falha ao registrar usuário", http.StatusInternalServerError)
		return
	}

	// Auto-login: gera token após registro bem-sucedido
	tokenResponse, loginErr := module.authService.LoginLocal(r.Context(), LoginInput{
		Email:    request.Email,
		Password: request.Password,
	})
	if loginErr != nil {
		log.Printf("Auto-login after register error: %v", loginErr)
		// Registro OK, mas auto-login falhou - retorna sucesso sem token
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		return
	}

	// Seta cookie com access_token
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    tokenResponse.AccessToken,
		Path:     "/",
		MaxAge:   int(tokenResponse.ExpiresIn),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(tokenResponse); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
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
}
