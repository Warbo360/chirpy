// Package auth that provides authentication functionality and helper functions related to that for users in the database
package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType string

const (
	TokenTypeAccess TokenType = "chirpy-access"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hashedPassword, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}
	return match, nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    string(TokenTypeAccess),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	})
	tokenString, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	// Claims struct to get payload as into the token

	claims := jwt.RegisteredClaims{}

	// Parses token string with pointer to claims type, callback checks for proper signing method

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	// Access issuer field from the token and see if it matches access issuer

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}
	if issuer != string(TokenTypeAccess) {
		return uuid.Nil, errors.New("invalid issuer")
	}

	// Extracts UserID from subject field in the token

	userID, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	// Parses UserID string into a UUID type and returns it

	parsedUUID, err := uuid.Parse(userID)
	if err != nil {
		return uuid.Nil, err
	}
	return parsedUUID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	// Extract Authorization header from http.Header if it exists

	auth := headers.Get("Authorization")
	if auth == "" {
		return "", errors.New("no Authorization header")
	}

	// Split at whitespace once, if in the <BEARER TOKEN_STRING> format will give a length of 2, verification of both
	// fact as well

	splitHeader := strings.Fields(auth)
	if len(splitHeader) != 2 || !strings.EqualFold(splitHeader[0], "Bearer") {
		return "", errors.New("unrecognized authorization header format")
	}

	return splitHeader[1], nil
}

func MakeRefreshToken() string {
	key := make([]byte, 32)
	rand.Read(key)
	return hex.EncodeToString(key)
}

// Functionally the same as GetBearerToken but looking for "ApiKey" instead

func GetAPIKey(headers http.Header) (string, error) {
	auth := headers.Get("Authorization")
	if auth == "" {
		return "", errors.New("no Authorization header")
	}

	splitAuthHeader := strings.Fields(auth)
	if len(splitAuthHeader) != 2 || !strings.EqualFold(splitAuthHeader[0], "ApiKey") {
		return "", errors.New("unrecognized Authorization header format for ApiKey")
	}

	return splitAuthHeader[1], nil
}
