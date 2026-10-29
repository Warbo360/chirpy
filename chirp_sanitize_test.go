package main

import (
	"testing"
)

func TestChirpSanitizeBadWord(t *testing.T) {
	chirp := "This is a kerfuffle opinion I need to share with the world"
	want := "This is a **** opinion I need to share with the world"
	got := helperChirpSanitize(chirp)
	if want != got {
		t.Errorf(`chirpSanitize("This is a kerfuffle opinion I need to share with the world") = %q, want match for %q`, got, want)
	}
}

func TestChirpSanitizeNoBadWord(t *testing.T) {
	chirp := "There is not a bad word to sanitize in this chirp so that is good!"
	want := "There is not a bad word to sanitize in this chirp so that is good!"
	got := helperChirpSanitize(chirp)
	if want != got {
		t.Errorf(`chirpSanitize("There is not a bad word to sanitize in this chirp so that is good!") = %q, want match for %q`, got, want)
	}
}
