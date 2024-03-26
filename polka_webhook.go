package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
)

type webhookRequest struct {
	Data struct {
		User_id int `json:"user_id"`
	} `json:"data"`
	Event string `json:"event"`
}

func (cfg *apiConfig) polkaWebhookHandler(w http.ResponseWriter, r *http.Request) {
	defer somethingWentWrong(w)

	authorization := r.Header.Get("Authorization")
	api_token := strings.TrimPrefix(authorization, "ApiKey ")
	if api_token != os.Getenv("POLKA_API_TOKEN") {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Unauthorized"})
		return
	}

	var webhook webhookRequest
	if err := json.NewDecoder(r.Body).Decode(&webhook); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid JSON"})
		return
	}

	if webhook.Event != "user.upgraded" {
		w.WriteHeader(http.StatusOK)
		return
	}

	userId := webhook.Data.User_id

	err := cfg.DB.UpgradeUser(userId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "User not found"})
		return
	}

	w.WriteHeader(http.StatusOK)
}
