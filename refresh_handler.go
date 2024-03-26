package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type RefreshResponse struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) refreshHandler(w http.ResponseWriter, r *http.Request) {
	defer somethingWentWrong(w)

	authorization := r.Header.Get("Authorization")
	if authorization == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Unauthorized"})
		return
	}

	authorization = strings.TrimPrefix(authorization, "Bearer ")

	_, err := cfg.DB.GetRevokedToken(authorization)
	if err == nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Unauthorized"})
		return
	}

	token, err := cfg.parseToken(authorization)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Unauthorized"})
		return
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Unauthorized"})
		return
	}

	if claims.Issuer != "chirpy-refresh" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Unauthorized"})
		return
	}

	userId, err := strconv.Atoi(claims.Subject)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Error determining user"})
		return
	}

	access_token, err := cfg.generateAccessToken(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Error creating token"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(RefreshResponse{Token: access_token}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Error creating token"})
		return
	}
}
