package main

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func ReturnError(w http.ResponseWriter, responseCode int, errorCode int, errorMessage string) {
	errStruct := Error{
		Code:    errorCode,
		Message: errorMessage,
	}

	errJson, _ := json.Marshal(errStruct)

	// now return stuff
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseCode)
	w.Write(errJson)
}
