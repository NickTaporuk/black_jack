//go:build darwin

package mtwrapper

/*
#cgo LDFLAGS: -L${SRCDIR} -lmtgenerator_macos
#include <stdint.h>
#include <stdbool.h>
#include <stdlib.h>

// External C++ functions
extern uint64_t* generate_random_tokens(uint64_t rangeGen);
extern int get_random_index(uint64_t* tokens, uint64_t rangeGen);
extern uint64_t get_random_number(uint64_t* tokens, uint64_t rangeGen);
extern bool drop_bonus(uint64_t rangeGen);
*/
import "C"
import "unsafe"

// GenerateRandomTokens wraps the C++ function to generate random tokens
// GenerateRandomTokens returns the slice and a free function
func GenerateRandomTokens(rangeGen uint64) ([]uint64, func()) {
	tokensPtr := C.generate_random_tokens(C.uint64_t(rangeGen))
	tokensSlice := (*[1 << 30]uint64)(unsafe.Pointer(tokensPtr))[:rangeGen:rangeGen]

	// Free function to be called by the caller
	freeFunc := func() {
		C.free(unsafe.Pointer(tokensPtr))
	}

	return tokensSlice, freeFunc
}

// GetRandomIndex wraps the C++ function to get a random index
func GetRandomIndex(tokens []uint64) int {
	return int(C.get_random_index((*C.uint64_t)(unsafe.Pointer(&tokens[0])), C.uint64_t(len(tokens))))
}

// GetRandomNumber wraps the C++ function to get a random number
func GetRandomNumber(tokens []uint64) uint64 {
	return uint64(C.get_random_number((*C.uint64_t)(unsafe.Pointer(&tokens[0])), C.uint64_t(len(tokens))))
}

// DropBonus wraps the C++ function to determine if a bonus should drop
func DropBonus(rangeGen uint64) bool {
	return bool(C.drop_bonus(C.uint64_t(rangeGen)))
}
