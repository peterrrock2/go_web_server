package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Id           int    `json:"id"`
	Email        string `json:"email"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	defer somethingWentWrong(w)

	var loginRequest LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid JSON"})
		return
	}

	user, err := cfg.DB.GetUserByEmail(loginRequest.Email)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid email or password"})
		return
	}

	if !user.CheckPassword(loginRequest.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid email or password"})
		return
	}

	access_token, err := cfg.generateAccessToken(user.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Error creating token"})
		return
	}

	refresh_token, err := cfg.generateRefreshToken(user.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Error creating token"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user.ToUserResponse(access_token, refresh_token)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Error encoding response"})
		return
	}
}

func (cfg *apiConfig) generateAccessToken(userId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24).UTC()),
		Subject:   strconv.Itoa(userId),
	})
	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (cfg *apiConfig) generateRefreshToken(userId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy-refresh",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 60).UTC()),
		Subject:   strconv.Itoa(userId),
	})
	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
