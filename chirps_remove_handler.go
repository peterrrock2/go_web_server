package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) chirpsRemoveHandler(w http.ResponseWriter, r *http.Request) {
	defer somethingWentWrong(w)

	authorization := r.Header.Get("Authorization")
	if authorization == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Unauthorized"})
		return
	}

	authorization = strings.TrimPrefix(authorization, "Bearer ")

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

	if claims.Issuer != "chirpy-access" {
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

	chirpIDstr := chi.URLParam(r, "id")
	chirpID, err := strconv.Atoi(chirpIDstr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid chirp ID"})
		return
	}

	dbChirp, err := cfg.DB.GetChirp(chirpID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Chirp not found"})
		return
	}

	if dbChirp.AuthorId != userId {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Forbidden"})
		return
	}

	err = cfg.DB.RemoveChirp(chirpID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Error removing chirp"})
		return
	}
}
