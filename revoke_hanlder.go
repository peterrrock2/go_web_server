package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (cfg *apiConfig) revokeHandler(w http.ResponseWriter, r *http.Request) {
	defer somethingWentWrong(w)

	authorization := r.Header.Get("Authorization")
	if authorization == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Unauthorized"})
		return
	}

	authorization = strings.TrimPrefix(authorization, "Bearer ")

	err := cfg.DB.AddRevokedToken(authorization)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Error revoking token"})
		return
	}

	w.WriteHeader(http.StatusOK)
}
