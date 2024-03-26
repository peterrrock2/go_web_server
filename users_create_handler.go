package main

import (
	"encoding/json"
	"net/http"
)

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cfg *apiConfig) usersCreateHandler(w http.ResponseWriter, r *http.Request) {
	defer somethingWentWrong(w)

	var newUser User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid JSON"})
		return
	}

	createdUser, err := cfg.DB.CreateUsers(newUser.Email, newUser.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Error creating user"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createdUser.ToUserResponse("", "")); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Error encoding response"})
		return
	}
}
