//go:build windows

package mtwrapper

/*
#cgo CFLAGS: -I${SRCDIR}
#cgo LDFLAGS: -lmtgenerator_windows
#include "mt_generator.h"
*/
import "C"
import "unsafe"

func GenerateRandomTokens(rangeGen uint64) ([]uint64, func()) {
	if rangeGen == 0 {
		return nil, func() {}
	}

	ptr := C.generate_random_tokens(C.uint64_t(rangeGen))
	if ptr == nil {
		return nil, func() {}
	}

	slice := unsafe.Slice((*uint64)(unsafe.Pointer(ptr)), rangeGen)

	freeFn := func() {
		C.free(unsafe.Pointer(ptr))
	}

	return slice, freeFn
}

func GetRandomIndex(tokens []uint64) int {
	if len(tokens) == 0 {
		return 0
	}
	return int(C.get_random_index(
		(*C.uint64_t)(unsafe.Pointer(&tokens[0])),
		C.uint64_t(len(tokens)),
	))
}

func GetRandomNumber(tokens []uint64) uint64 {
	if len(tokens) == 0 {
		return 0
	}
	return uint64(C.get_random_number(
		(*C.uint64_t)(unsafe.Pointer(&tokens[0])),
		C.uint64_t(len(tokens)),
	))
}

func DropBonus(rangeGen uint64) bool {
	if rangeGen == 0 {
		return false
	}
	return bool(C.drop_bonus(C.uint64_t(rangeGen)))
}
