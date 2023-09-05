package main

import "C"
import (
	"math/rand"
	"time"

	"github.com/seehuhn/mt19937"
)

// GenerateRandomTokens generates random tokens
// rangeGen - amount of the tokens
// tokensSlice - slice of the generated tokens
func GenerateRandomTokens(rangeGen uint8) (tokensSlice []uint64) {
	rng := mt19937.New()
	// calculate the range from min to max by % from the cents
	// generate x random tokens like a 100% of the range
	tokensSlice = make([]uint64, rangeGen)
	for i := uint8(0); i < rangeGen; i++ {
		tokensSlice[i] = rng.Uint64()
	}

	return tokensSlice
}

// GetRandomIndex retrieves the randomly selected index
func GetRandomIndex(tokensSlice []uint64) (randomIndex int) {
	rng := rand.New(mt19937.New())
	rng.Seed(time.Now().UnixNano())
	randomIndex = rng.Intn(len(tokensSlice))

	return randomIndex
}

// GetRandomNumber retrieves the randomly selected number
func GetRandomNumber(tokensSlice []uint64) (randomNumber uint64) {
	randomIndex := GetRandomIndex(tokensSlice)
	// Retrieve the randomly selected number
	randomNumber = tokensSlice[randomIndex]

	return randomNumber
}

// DropPersonalBonus drops the personal bonus
func DropPersonalBonus(rangeGen uint8) bool {
	tokensSlice := GenerateRandomTokens(rangeGen)

	randomNumber1 := GetRandomNumber(tokensSlice)
	randomNumber2 := GetRandomNumber(tokensSlice)
	if randomNumber2 == randomNumber1 {
		return true
	}
	return false
}

func main() {}
