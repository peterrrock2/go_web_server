package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (cfg *apiConfig) chirpsGetAllHandler(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.DB.GetChirps()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Could not retrieve any chirps"})
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{Id: dbChirp.Id, Body: dbChirp.Body})
	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].Id < chirps[j].Id
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&chirps); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Error encoding response"})
	}
}

func (cfg *apiConfig) chirpsGetHandler(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&dbChirp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Error encoding response"})
	}
}
