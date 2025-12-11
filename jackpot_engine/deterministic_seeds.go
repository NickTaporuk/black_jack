package jackpot_engine

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
)

// DeterministicSeed generates a deterministic uint64 seed
func DeterministicSeed(masterSeed uint64, labels ...string) uint64 {
	h := hmac.New(sha256.New, u64ToBytes(masterSeed))
	for _, s := range labels {
		_, _ = h.Write([]byte(s))
	}
	sum := h.Sum(nil)

	return binary.LittleEndian.Uint64(sum[:8])
}

func u64ToBytes(v uint64) []byte {
	var b [8]byte
	binary.LittleEndian.PutUint64(b[:], v)
	return b[:]
}
