package auth_test

import (
	"testing"

	"github.com/Warbo360/chirpy/internal/auth"
)

// Tests proper hash generation and passing password verification.

func TestHashPassword(t *testing.T) {
	password := "Password!"

	// Generate hash of the password, assert error should be nil and hash is non-empty

	hash, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error during hash creation, got %v", err)
	}
	if hash == "" {
		t.Fatal("expected a non-empty hash string")
	}

	// Check that password is verified correctly

	match, err := auth.CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("expected no error during password verification, got %v", err)
	}
	if !match {
		t.Fatal("expected matching password to be true, got false")
	}
}

// Tests proper error generation for a mismatched password

func TestMismatchHashPassword(t *testing.T) {
	password := "Password!"

	// Generate hash of the password, assert error should be nil and hash is non-empty

	hash, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error during hash creation, got %v", err)
	}
	if hash == "" {
		t.Fatal("expected a non-empty hash string")
	}

	// Check that password is not verified when mismatched password is passed

	match, err := auth.CheckPasswordHash("NOTTHERIGHTPASSWORD!", hash)
	if err != nil {
		t.Fatalf("expected no error during password verification, got %v", err)
	}
	if match {
		t.Fatal("expected matching password to be false, got true")
	}
}
