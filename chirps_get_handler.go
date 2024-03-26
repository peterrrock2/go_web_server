package main

import (
	"chirpy/internal/database"
	"encoding/json"
	"net/http"
	"sort"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (cfg *apiConfig) chirpsGetAllHandler(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Query().Get("author_id")
	// s is a string that contains the value of the author_id query parameter
	// if it exists, or an empty string if it doesn't

	var author_id int
	var dbChirps []database.Chirp
	var err error
	if s != "" {
		author_id, err = strconv.Atoi(s)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid author ID"})
			return
		}
	} else {
		author_id = -1
	}

	dbChirps, err = cfg.DB.GetChirps(author_id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Could not retrieve any chirps"})
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{AuthorId: dbChirp.AuthorId, Id: dbChirp.Id, Body: dbChirp.Body})
	}

	var sortasc bool
	if r.URL.Query().Get("sort") == "desc" {
		sortasc = false
	} else {
		sortasc = true
	}

	sort.Slice(chirps, func(i, j int) bool {
		if sortasc {
			return chirps[i].Id < chirps[j].Id
		}
		return chirps[i].Id > chirps[j].Id
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
