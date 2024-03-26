package database

import (
	"time"
)

type RevokedToken struct {
	RefreshToken string    `json:"refresh_token"`
	Time         time.Time `json:"time"`
}

func (db *DB) GetRevokedToken(token string) (RevokedToken, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return RevokedToken{}, err
	}

	revokedToken, ok := dbStructure.RevokedTokens[token]
	if !ok {
		return RevokedToken{}, ErrNotExist
	}

	return RevokedToken{RefreshToken: token, Time: revokedToken}, nil
}

func (db *DB) AddRevokedToken(refresh_token string) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	dbStructure.RevokedTokens[refresh_token] = time.Now()

	err = db.WriteDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}
