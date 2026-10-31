package auth_test

import (
	"testing"
	"time"

	"github.com/Warbo360/chirpy/internal/auth"
	"github.com/google/uuid"
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

// Test for token generation from UserID and secret string

func TestMakeJWT(t *testing.T) {
	expiryDelta := time.Minute * 30
	userID := uuid.New()

	// Test no error during token string generation

	tokenString, err := auth.MakeJWT(userID, "tokenSecret", expiryDelta)
	if err != nil {
		t.Fatalf("expected no error during tokenString generation, got %v", err)
	}

	// Test no error during jwt.Token gen from ValidateJWT

	regenUserID, err := auth.ValidateJWT(tokenString, "tokenSecret")
	if err != nil {
		t.Fatalf("expected no error during token validation, got %v", err)
	}
	if userID != regenUserID {
		t.Fatalf("expected validated UserID to match, wanted: %v, got: %v", userID, regenUserID)
	}
}

// Test for invalid token
