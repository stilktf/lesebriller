package main

import (
	"encoding/json"
	"net/http"
)

func UsersAuth(w http.ResponseWriter, r *http.Request) {
	// Get headers
	usernameHeader := r.Header.Get("x-auth-user")
	if usernameHeader == "" {
		ReturnError(w, 403, 2003, "Invalid request.")
		return
	}
	authKeyHeader := r.Header.Get("x-auth-key")
	if authKeyHeader == "" {
		ReturnError(w, 403, 2003, "Invalid request.")
		return
	}

	authed, err := AuthenticateUser(usernameHeader, authKeyHeader)
	if err != nil {
		ReturnError(w, 502, 2000, "Unknown server error.")
		return
	}

	if authed {
		returnStruct := Authenticated{
			Authenticated: "OK",
		}

		returnJson, err := json.Marshal(returnStruct)
		if err != nil {
			ReturnError(w, 502, 2000, "Unknown server error.")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(returnJson)
		return
	} else {
		ReturnError(w, 401, 2001, "Unauthorized")
	}
}

func AuthenticateUser(username string, authkey string) (authed bool, err error) {
	// Find user
	userFetch := User{
		Username: username,
	}

	fetchedUser := User{}

	result := DbConnection.Model(&User{}).Find(&userFetch).Scan(&fetchedUser)
	if result.Error != nil {
		return false, result.Error
	}

	// check if password matches hash
	match, err := ComparePasswordWithHash(authkey, fetchedUser.AuthKey)
	if err != nil {
		return false, err
	}
	return match, err
}
