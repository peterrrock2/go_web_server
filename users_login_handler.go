package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type LoginRequest struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

type LoginResponse struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("loginHandler")
	defer somethingWentWrong(w, r)

	var loginRequest LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid JSON"})
		return
	}
	fmt.Println(loginRequest)

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

	if loginRequest.ExpiresInSeconds <= 0 {
		loginRequest.ExpiresInSeconds = 24 * 60 * 60
	}
	loginRequest.ExpiresInSeconds = min(loginRequest.ExpiresInSeconds, 24*60*60)
	fmt.Println(loginRequest.ExpiresInSeconds)
	token, err := cfg.generateToken(user.Id, loginRequest.ExpiresInSeconds)
	fmt.Println(token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Error creating token"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user.ToUserResponse(token)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Error encoding response"})
		return
	}
}

func (cfg *apiConfig) generateToken(userId, expiresInSeconds int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(expiresInSeconds)).UTC()),
		Subject:   strconv.Itoa(userId),
	})
	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
