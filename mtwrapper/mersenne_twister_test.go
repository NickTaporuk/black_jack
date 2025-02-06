package mtwrapper

import (
	"testing"
)

// TestGenerateRandomTokens verifies that tokens are generated correctly
func TestGenerateRandomTokens(t *testing.T) {
	rangeGen := uint64(100)
	tokens, resetSlice := GenerateRandomTokens(rangeGen)
	defer resetSlice()
	if len(tokens) != int(rangeGen) {
		t.Errorf("Expected %d tokens, got %d", rangeGen, len(tokens))
	}

	for _, token := range tokens {
		if token == 0 {
			t.Error("Token value should not be zero")
		}
	}
}

// TestGetRandomIndex ensures the index is within valid range
func TestGetRandomIndex(t *testing.T) {
	tokens, resetSlice := GenerateRandomTokens(50)
	defer resetSlice()
	index := GetRandomIndex(tokens)

	if index < 0 || index >= len(tokens) {
		t.Errorf("Random index out of range: %d", index)
	}
}

// TestGetRandomNumber checks if a random number is retrieved correctly
func TestGetRandomNumber(t *testing.T) {
	tokens, resetSlice := GenerateRandomTokens(50)
	defer resetSlice()
	number := GetRandomNumber(tokens)
	found := false

	for _, token := range tokens {
		if token == number {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Random number %d not found in tokens", number)
	}
}

// TestDropBonus runs multiple trials to check the drop rate
func TestDropBonus(t *testing.T) {
	rangeGen := uint64(100)
	trials := 10000
	drops := 0

	for i := 0; i < trials; i++ {
		if DropBonus(rangeGen) {
			drops++
		}
	}

	expectedRate := 1.0 / float64(rangeGen)
	actualRate := float64(drops) / float64(trials)
	tolerance := expectedRate * 0.5 // Allow ±50% tolerance due to randomness

	if actualRate < expectedRate-tolerance || actualRate > expectedRate+tolerance {
		t.Errorf("Drop rate %.4f is outside expected range (expected ~%.4f ±%.4f)", actualRate, expectedRate, tolerance)
	}
}

// BenchmarkGenerateRandomTokens measures performance of token generation
func BenchmarkGenerateRandomTokens(b *testing.B) {
	rangeGen := uint64(1000)
	for i := 0; i < b.N; i++ {
		GenerateRandomTokens(rangeGen)
	}
}

// BenchmarkGetRandomIndex measures performance of random index retrieval
func BenchmarkGetRandomIndex(b *testing.B) {
	tokens, resetSlice := GenerateRandomTokens(1000)
	defer resetSlice()
	for i := 0; i < b.N; i++ {
		GetRandomIndex(tokens)
	}
}

// BenchmarkGetRandomNumber measures performance of random number retrieval
func BenchmarkGetRandomNumber(b *testing.B) {
	tokens, resetSlice := GenerateRandomTokens(1000)
	defer resetSlice()
	for i := 0; i < b.N; i++ {
		GetRandomNumber(tokens)
	}
}

// BenchmarkDropBonus measures performance of bonus drop simulation
func BenchmarkDropBonus(b *testing.B) {
	rangeGen := uint64(1000)
	for i := 0; i < b.N; i++ {
		DropBonus(rangeGen)
	}
}
