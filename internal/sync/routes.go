package sync

import (
	"encoding/json"
	"log"
	"net/http"
	"portfolio/internal/auth"
	"portfolio/internal/jwt"

	"github.com/gorilla/mux"
)

type GithubSyncModule struct {
	jwtService *jwt.JWTService
	userRepo   auth.UserRepository
}

func NewGithubSyncModule(jwtService *jwt.JWTService,userRepo   auth.UserRepository) *GithubSyncModule {
	return &GithubSyncModule{
		jwtService: jwtService,
		userRepo:   userRepo,
	}
}

func (module *GithubSyncModule) RegisterRoutes() *mux.Router {
	// No routes for now
	router := mux.NewRouter()
	router.HandleFunc("/github", module.jwtService.RequiredAutenticationMiddleware(module.syncGithubUserData)).Methods("GET")
	return router
}

func (module *GithubSyncModule) syncGithubUserData(w http.ResponseWriter, r *http.Request){
	userCtx := jwt.GetUserCurrentUser(r.Context())
	userId := userCtx.ID
	log.Printf("Syncing GitHub data for user ID: %s", userId)
	user,err := module.userRepo.Find(r.Context(), userId)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	githubAccessToken := user.GetGithubAccessToken()

	if githubAccessToken == nil {
		http.Error(w, "GitHub access token not found", http.StatusBadRequest)
		return
	}

	githubData,err := SyncGithubData(r.Context(), *githubAccessToken)
	if err != nil {
		http.Error(w, "Failed to sync GitHub data", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(githubData)
}