package auth_test

import (
	"net/http"
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
	hash, _ := auth.HashPassword(password)

	tests := []struct {
		name           string
		passwordString string
		hashString     string
		wantMatch      bool
		wantErr        bool
	}{
		{
			name:           "Valid password",
			passwordString: password,
			hashString:     hash,
			wantMatch:      true,
			wantErr:        false,
		},
		{
			name:           "Invalid Password",
			passwordString: "Not the same password!",
			hashString:     hash,
			wantMatch:      false,
			wantErr:        false,
		},
		{
			name:           "Invalid hash",
			passwordString: password,
			hashString:     "",
			wantMatch:      false,
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMatch, err := auth.CheckPasswordHash(tt.passwordString, tt.hashString)
			if (err != nil) != tt.wantErr {
				t.Fatalf("CheckPasswordHash() got: %v, want: %v", err, tt.wantErr)
			}
			if gotMatch != tt.wantMatch {
				t.Fatalf("CheckPasswordHash() got: %v, want: %v", gotMatch, tt.wantMatch)
			}
		})
	}
}

// Test for token generation from UserID and secret string

func TestMakeJWT(t *testing.T) {
	userID := uuid.New()
	validToken, _ := auth.MakeJWT(userID, "secret", time.Hour)

	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantUserID  uuid.UUID
		wantErr     bool
	}{
		{
			name:        "Valid Token",
			tokenString: validToken,
			tokenSecret: "secret",
			wantUserID:  userID,
			wantErr:     false,
		},
		{
			name:        "Invalid Token",
			tokenString: "Invalid.Token.String",
			tokenSecret: "secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "Invalid Secret",
			tokenString: validToken,
			tokenSecret: "Not secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := auth.ValidateJWT(tt.tokenString, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateJWT() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if gotUserID != tt.wantUserID {
				t.Fatalf("ValidateJWT() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}

func TestBearerToken(t *testing.T) {
	header := make(http.Header)
	badHeader := make(http.Header)
	token, _ := auth.MakeJWT(uuid.New(), "secret", time.Hour)
	header.Set("Authorization", "Bearer "+token)
	badHeader.Add("Authorization", token)

	tests := []struct {
		name            string
		testHeader      http.Header
		wantTokenString string
		wantErr         bool
	}{
		{
			name:            "ValidHeader",
			testHeader:      header,
			wantTokenString: token,
			wantErr:         false,
		},
		{
			name:            "InvalidHeader",
			testHeader:      make(http.Header),
			wantTokenString: "",
			wantErr:         true,
		},
		{
			name:            "InvalidAuthFormat",
			testHeader:      badHeader,
			wantTokenString: "",
			wantErr:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToken, err := auth.GetBearerToken(tt.testHeader)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetBearerToken error got: %v, want: %v\n", err, tt.wantErr)
			}
			if gotToken != tt.wantTokenString {
				t.Fatalf("GetBearerToken gotToken: %v, want: %v\n", gotToken, tt.wantTokenString)
			}
		})
	}
}
