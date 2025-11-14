package mtwrapper_test

import (
	"testing"

	"gitlabce.champion.tm/champion-prj/mersenne-twister-shared-object/mtwrapper"
)

// Benchmark DropBonus
func BenchmarkDropBonus(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = mtwrapper.DropBonus(100)
	}
}
