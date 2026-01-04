package portfolio

import (
	"encoding/json"
	"log"
	"net/http"
	"portfolio/internal/jwt"

	"github.com/gorilla/mux"
)

type PortfolioModule struct {
	service      *PortfolioService
	tokenService *jwt.TokenService
}

func NewPortfolioModule(service *PortfolioService, tokenService *jwt.TokenService) *PortfolioModule {
	return &PortfolioModule{
		service:      service,
		tokenService: tokenService,
	}
}

func (module *PortfolioModule) RegisterRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/me", jwt.AuthMiddleware(module.tokenService, module.getMyProfile)).Methods("GET")
	router.HandleFunc("/", jwt.AuthMiddleware(module.tokenService, module.createProfile)).Methods("POST")
	router.HandleFunc("/", jwt.AuthMiddleware(module.tokenService, module.updateProfile)).Methods("PUT")
	router.HandleFunc("/", jwt.AuthMiddleware(module.tokenService, module.patchProfile)).Methods("PATCH")
	router.HandleFunc("/", jwt.AuthMiddleware(module.tokenService, module.deleteProfile)).Methods("DELETE")

	return router
}

func (module *PortfolioModule) getMyProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(jwt.UserIDKey).(string) // Assumindo que o AuthMiddleware injeta "user_id"

	profile, err := module.service.GetMyProfile(r.Context(), userID)
	if err != nil {
		if err == ErrProfileNotFound {
			http.Error(w, "Profile not found", http.StatusNotFound)
			return
		}
		log.Printf("GetMyProfile error: %v", err)
		http.Error(w, "Failed to get profile", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}

func (module *PortfolioModule) createProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(jwt.UserIDKey).(string)

	var input SaveProfileInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Printf("CreateProfile decode error: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	profile, err := module.service.CreateProfile(r.Context(), userID, input)
	if err != nil {
		if err == ErrProfileAlreadyExists {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		log.Printf("CreateProfile error: %v", err)
		http.Error(w, "Failed to create profile", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(profile)
}

func (module *PortfolioModule) updateProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(jwt.UserIDKey).(string)

	var input SaveProfileInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	profile, err := module.service.UpdateProfile(r.Context(), userID, input)
	if err != nil {
		if err == ErrProfileNotFound {
			http.Error(w, "Profile not found", http.StatusNotFound)
			return
		}
		log.Printf("UpdateProfile error: %v", err)
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(profile)
}

func (module *PortfolioModule) patchProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(jwt.UserIDKey).(string)

	var input PatchProfileDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	profile, err := module.service.PatchProfile(r.Context(), userID, input)
	if err != nil {
		if err == ErrProfileNotFound {
			http.Error(w, "Profile not found", http.StatusNotFound)
			return
		}
		log.Printf("UpdateProfile error: %v", err)
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(profile)
}

func (module *PortfolioModule) deleteProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(jwt.UserIDKey).(string)

	if err := module.service.DeleteProfile(r.Context(), userID); err != nil {
		if err == ErrProfileNotFound {
			http.Error(w, "Profile not found", http.StatusNotFound)
			return
		}
		log.Printf("DeleteProfile error: %v", err)
		http.Error(w, "Failed to delete profile", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
