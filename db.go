package main

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        int
	Username  string
	AuthKey   string // This will be hashed using argon2. Uses an md5 hashing format for some reason internally, though. Can't do much about it, this is how KOReader works
	Documents []Document
}

type Document struct {
	gorm.Model
	ID         int
	UserID     uint    // foreign key to User
	DocumentID string  // MD5 hash of the document being synced
	Progress   string  // for example "/html/body/div[1]/img.0" bla bla bla
	Percentage float32 // for example 0.18 or similar
	Device     string  // the name of device being used (like Kobo BW)
	DeviceId   string  // not sure how this is generated, but seems to be an MD5
	Timestamp  int64   // UNIX timestamp (when the document was last synced?)
}

var DbConnection *gorm.DB

func Db() {
	slog.Info("creating db folder if it does not exist")
	newpath := filepath.Join(".", "db")
	err := os.MkdirAll(newpath, os.ModePerm)
	if err != nil {
		slog.Error("something went wrong creating folder", "err", err.Error())
	}

	slog.Info("connecting to database")
	DbConnection, err = gorm.Open(sqlite.Open("db/lesebriller.db"), &gorm.Config{})
	if err != nil {
		slog.Error("can't connect to database :(", "err", err.Error())
		return
	}
	slog.Info("migrating structs")
	DbConnection.AutoMigrate(&User{})
	DbConnection.AutoMigrate(&Document{})
}
