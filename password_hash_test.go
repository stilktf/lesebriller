package main

import (
	"testing"
)

func TestHashComparison(t *testing.T) {
	InitializeHashing()
	password := "rand0mt3st2025!"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf(`HashPassword("") = %q. returned an error %v`, password, err.Error())
	}
	match, err := ComparePasswordWithHash(password, hash)
	if err != nil {
		t.Fatalf(`ComparePasswordWithHash("") = %q & %q, returned an error %v`, password, hash, err.Error())
	}
	if !match {
		t.Fatalf(`ComparePasswordWithHash("") = %q & %q, apparently don't match`, password, hash)
	}
}
