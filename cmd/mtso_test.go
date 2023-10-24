package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRandomTokens(t *testing.T) {
	rangeGen := uint64(10)
	tokens := GenerateRandomTokens(rangeGen)

	if len(tokens) != int(rangeGen) {
		t.Errorf("Expected %d tokens, got %d", rangeGen, len(tokens))
	}

	// Check for uniqueness (this might not always be true, but for our range, it should be)
	tokenMap := make(map[uint64]bool)
	for _, token := range tokens {
		if tokenMap[token] {
			t.Errorf("Token %d appears more than once", token)
		}
		tokenMap[token] = true
	}
}

func TestGetRandomIndex(t *testing.T) {
	tokens := GenerateRandomTokens(10)
	randomIndex := GetRandomIndex(tokens)

	if randomIndex < 0 || randomIndex >= len(tokens) {
		t.Errorf("Random index out of range: %d", randomIndex)
	}
}

func TestGetRandomNumber(t *testing.T) {
	tokens := GenerateRandomTokens(10)
	randomNumber := GetRandomNumber(tokens)

	found := false
	for _, token := range tokens {
		if token == randomNumber {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Random number %d not found in tokens slice", randomNumber)
	}
}

func TestDropBonus(t *testing.T) {
	dropped := DropBonus(1)
	// The DropBonus function is mostly about randomness.
	// There isn't much deterministic behavior to check.
	// But you can check it multiple times and see if it always/never returns true/false if needed.

	assert.True(t, dropped)
}
