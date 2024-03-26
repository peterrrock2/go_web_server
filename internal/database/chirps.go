package database

import (
	"encoding/json"
	"os"
	"sort"
)

type Chirp struct {
	AuthorId int    `json:"author_id"`
	Body     string `json:"body"`
	Id       int    `json:"id"`
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps(id int) ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	var chirpsArray []Chirp
	for _, chirp := range dbStructure.Chirps {
		if id == -1 {
			chirpsArray = append(chirpsArray, chirp)
		} else if chirp.AuthorId == id {
			chirpsArray = append(chirpsArray, chirp)
		}
	}

	sort.Slice(chirpsArray, func(i, j int) bool {
		return chirpsArray[i].Id < chirpsArray[j].Id
	})

	return chirpsArray, nil
}

func (db *DB) GetChirp(id int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirp, ok := dbStructure.Chirps[id]
	if !ok {
		return Chirp{}, ErrNotExist
	}

	return chirp, nil
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(author_id int, body string) (Chirp, error) {
	allChirps, err := db.GetChirps(-1)
	if err != nil {
		return Chirp{}, err
	}
	dbStructure, _ := db.loadDB()

	newChirp := Chirp{
		AuthorId: author_id,
		Body:     body,
		Id:       len(allChirps) + 1,
	}

	allChirps = append(allChirps, newChirp)

	chirps := make(map[int]Chirp)
	for _, chirp := range allChirps {
		chirps[chirp.Id] = chirp
	}

	chirps_json, err := json.Marshal(DBStructure{Chirps: chirps, Users: dbStructure.Users, RevokedTokens: dbStructure.RevokedTokens})
	if err != nil {
		return Chirp{}, err
	}

	db.mux.Lock()
	err = os.WriteFile(db.path, chirps_json, 0644)
	db.mux.Unlock()

	if err != nil {
		return Chirp{}, err
	}

	return newChirp, nil
}

func (db *DB) RemoveChirp(id int) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	_, ok := dbStructure.Chirps[id]
	if !ok {
		return ErrNotExist
	}

	delete(dbStructure.Chirps, id)

	chirps_json, err := json.Marshal(DBStructure{Chirps: dbStructure.Chirps, Users: dbStructure.Users, RevokedTokens: dbStructure.RevokedTokens})
	if err != nil {
		return err
	}

	db.mux.Lock()
	err = os.WriteFile(db.path, chirps_json, 0644)
	db.mux.Unlock()

	if err != nil {
		return err
	}

	return nil
}
