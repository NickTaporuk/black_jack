package main

import (
	"testing"
)

// -------------------------
// Benchmark tests
// -------------------------

func BenchmarkGenerateRandomTokens(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = GenerateRandomTokens(1000)
	}
}

func BenchmarkGetRandomIndex(b *testing.B) {
	tokens := GenerateRandomTokens(1000)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = GetRandomIndex(tokens)
	}
}

func BenchmarkDropBonus(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DropBonus(100)
	}
}
