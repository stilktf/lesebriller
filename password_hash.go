package main

import (
	"github.com/matthewhartstonge/argon2"
)

var Argon argon2.Config

// Initialize parameters for password hashing functions.
func InitializeHashing() {
	Argon = argon2.DefaultConfig()
}

func HashPassword(password string) (encodedHash string, err error) {
	encoded, err := Argon.HashEncoded([]byte(password))
	if err != nil {
		return "", err
	}

	return string(encoded), nil
}

func ComparePasswordWithHash(password, encodedHash string) (match bool, err error) {
	ok, err := argon2.VerifyEncoded([]byte(password), []byte(encodedHash))
	if err != nil {
		return false, err
	}
	return ok, nil
}
