package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) usersUpdateHandler(w http.ResponseWriter, r *http.Request) {
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

	var updatedUser User
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid JSON"})
		return
	}

	id, err := strconv.Atoi(claims.Subject)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Error updating user"})
		return
	}

	user, err := cfg.DB.UpdateUser(id, updatedUser.Email, updatedUser.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Error updating user"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user.ToUserResponse("", ""))
}

func (cfg *apiConfig) parseToken(tokenString string) (*jwt.Token, error) {

	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil
}
