package main

import (
	"chirpy/internal/database"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
	JWTSecret      string
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func main() {
	godotenv.Load()

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatalf("Error creating database: %s\n", err)
	}

	cfg := &apiConfig{
		fileserverHits: 0,
		DB:             db,
		JWTSecret:      os.Getenv("JWT_SECRET"),
	}
	mainRouter := chi.NewRouter()

	mainRouter.Use(middlewareCors)

	fsHandler := cfg.middlewareMetricsInc(
		http.StripPrefix(
			"/app",
			http.FileServer(http.Dir("./app")),
		),
	)

	mainRouter.Handle("/app", fsHandler)
	mainRouter.Handle("/app/*", fsHandler)

	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", readinessHandler)
	apiRouter.Get("/reset", cfg.resetHandler)

	apiRouter.Post("/chirps", cfg.chirpsCreateHandler)
	apiRouter.Get("/chirps", cfg.chirpsGetAllHandler)
	apiRouter.Get("/chirps/{id}", cfg.chirpsGetHandler)
	apiRouter.Delete("/chirps/{id}", cfg.chirpsRemoveHandler)

	apiRouter.Post("/users", cfg.usersCreateHandler)
	apiRouter.Put("/users", cfg.usersUpdateHandler)

	apiRouter.Post("/login", cfg.loginHandler)

	apiRouter.Post("/refresh", cfg.refreshHandler)
	apiRouter.Post("/revoke", cfg.revokeHandler)

	apiRouter.Post("/polka/webhooks", cfg.polkaWebhookHandler)

	mainRouter.Mount("/api", apiRouter)

	adminRouter := chi.NewRouter()
	adminRouter.Get("/metrics", cfg.metricsGetHandler)

	mainRouter.Mount("/admin", adminRouter)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mainRouter,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error starting server: %s\n", err)
	}
}

func somethingWentWrong(w http.ResponseWriter) {
	if err := recover(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Something went wrong"})
	}
}
