package main

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"gorm.io/gorm"
)

func SyncsProgress(w http.ResponseWriter, r *http.Request) {
	progressStruct := Progress{}

	err := json.NewDecoder(r.Body).Decode(&progressStruct)
	if err != nil {
		ReturnError(w, 500, 2000, "Unknown server error.")
	}

	// auth user
	authed, err := AuthenticateUser(r.Header.Get("x-auth-user"), r.Header.Get("x-auth-key"))
	if err != nil {
		ReturnError(w, 500, 2000, "Unknown server error.")
		return
	}
	if !authed {
		ReturnError(w, 401, 2001, "Unauthorized")
		return
	}

	var user User
	var document Document
	DbConnection.First(&User{
		Username: r.Header.Get("x-auth-user"),
	}).Scan(&user)

	currentTimestamp := time.Now().Unix()
	freshDocument := Document{
		DocumentID: progressStruct.Document,
		Progress:   progressStruct.Progress,
		Percentage: progressStruct.Percentage,
		DeviceId:   progressStruct.DeviceID,
		Device:     progressStruct.Device,
		Timestamp:  currentTimestamp,
	}

	// See if the document already exists
	if err := DbConnection.Where("document_id = ? AND user_id = ?", progressStruct.Document, &user.ID).First(&document).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			DbConnection.Model(&user).Association("Documents").Append(&freshDocument)
		}
	} else {
		document.Progress = progressStruct.Progress
		document.Percentage = progressStruct.Percentage
		document.Timestamp = currentTimestamp

		if err := DbConnection.Save(&document).Error; err != nil {
			slog.Error("Could not save document!")
		}
	}

	returnStruct := Progress{
		Timestamp: int32(currentTimestamp),
		Document:  progressStruct.Document,
	}

	returnJson, err := json.Marshal(returnStruct)
	if err != nil {
		ReturnError(w, 502, 2000, "Unknown server error.")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(returnJson)
}

// This is the function for the route that gets progress data
func SyncsProgressPull(w http.ResponseWriter, r *http.Request) {

	// auth user
	authed, err := AuthenticateUser(r.Header.Get("x-auth-user"), r.Header.Get("x-auth-key"))
	if err != nil {
		ReturnError(w, 502, 2000, "Unknown server error.")
	}

	if !authed {
		ReturnError(w, 401, 2001, "Unauthorized")
	}

	requestedDocument := r.PathValue("document")

	var user User
	var document Document

	DbConnection.First(&User{
		Username: r.Header.Get("x-auth-user"),
	}).Scan(&user)

	if err := DbConnection.Where("document_id = ? AND user_id = ?", requestedDocument, user.ID).First(&document).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ReturnError(w, 403, 2004, "Field 'document' not provided.")
			return
		}
	} else {
		progressStruct := Progress{
			DeviceID:   document.DeviceId,
			Percentage: document.Percentage,
			Device:     document.Device,
			Timestamp:  int32(document.Timestamp),
			Progress:   document.Progress,
			Document:   document.DocumentID,
		}
		progressJson, err := json.Marshal(progressStruct)
		if err != nil {
			ReturnError(w, 502, 2000, "Unknown server error.")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(progressJson)
	}
}
