package main

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"gorm.io/gorm"
)

func UserCreate(w http.ResponseWriter, r *http.Request) {
	// Unmarshal request json
	us := UserCreateStruct{}

	err := json.NewDecoder(r.Body).Decode(&us)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if user already exists
	result := DbConnection.First(User{
		Username: us.Username,
	})
	// If it already exists, bail
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		ReturnError(w, 402, 2002, "Username is already registered.")
		return
	}

	// Turn password into hash
	hash, err := HashPassword(us.Password)
	if err != nil {
		ReturnError(w, 502, 2000, "Unknown server error.")
		slog.Error("something went wrong while hashing a password", "err", err.Error())
		return
	}

	// Create record in db
	user := User{Username: us.Username, AuthKey: hash}
	result = DbConnection.Create(&user)
	if result.Error != nil {
		ReturnError(w, 502, 2000, "Unknown server error.")
		slog.Error("something went wrong while creating a user", "err", result.Error.Error())
		return
	}

	// Return json confirming that the account has been created
	returnUs := UserCreateStruct{
		Username: us.Username,
	}

	returnJson, err := json.Marshal(returnUs)
	if err != nil {
		ReturnError(w, 502, 2000, "Unknown server error.")
		slog.Error("something went wrong while creating a user", "err", err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(returnJson)
}
