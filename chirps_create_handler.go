package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type Chirp struct {
	Body string `json:"body"`
	Id   int    `json:"id"`
}

func (cfg *apiConfig) chirpsCreateHandler(w http.ResponseWriter, r *http.Request) {
	defer somethingWentWrong(w, r)

	var newChirp Chirp
	if err := json.NewDecoder(r.Body).Decode(&newChirp); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid JSON"})
		return
	}

	cleaned, err := validateChirp(newChirp.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
		return
	}

	createdChirp, err := cfg.DB.CreateChirp(cleaned)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Error creating chirp"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(&createdChirp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Error encoding response"})
		return
	}
}

func validateChirp(body string) (string, error) {
	if len(body) > 140 {
		return "", errors.New("Chirp is too long")
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	cleaned := cleanChirp(body, badWords)
	return cleaned, nil
}

func cleanChirp(body string, badWords map[string]struct{}) string {
	chirpWords := strings.Split(body, " ")
	for i, word := range chirpWords {
		if _, ok := badWords[word]; ok {
			chirpWords[i] = "****"
		}
	}
	cleaned := strings.Join(chirpWords, " ")
	return cleaned
}
