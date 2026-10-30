package auth_test

import (
	"testing"

	"github.com/Warbo360/chirpy/internal/auth"
)

func TestHashPassword(t *testing.T) {
	password := "Password!"

	// Generate hash of the password

	hash, err := auth.HashPassword(password)
}
