//go:build linux

package jackpot_engine

/*
#cgo LDFLAGS: -L${SRCDIR} -ljackpot_engine_v1_linux

#include "jackpot_engine_v1.h"
*/
import "C"

import (
	"fmt"
	"sync"
	"sync/atomic"
	"unsafe"
)

// Volatility mirrors C++ enum
type Volatility int32

const (
	VolatilityLow    Volatility = 1
	VolatilityMedium Volatility = 2
	VolatilityHigh   Volatility = 3
)

var VolatilityName = map[string]Volatility{
	"low":    VolatilityLow,
	"medium": VolatilityMedium,
	"high":   VolatilityHigh,
}

var VolatilityIntToName = map[Volatility]string{
	VolatilityLow:    "low",
	VolatilityMedium: "medium",
	VolatilityHigh:   "high",
}

type Config struct {
	MinPoint   uint64
	MaxPoint   uint64
	Volatility Volatility
}

// Engine is safe for concurrent use from thousands of goroutines
type Engine struct {
	handle unsafe.Pointer // *C.JackpotEngineHandle
	closed uint64         // atomic
	mu     sync.Mutex     // protects destroy sequence
}

// New creates a fully deterministic engine from seed.
// Use the same seed on two different machines → identical jackpot sequences (perfect for audit & replay)
func New(seed uint64) (*Engine, error) {
	h := C.jackpot_engine_v1(C.uint64_t(seed))
	if h == nil {
		return nil, fmt.Errorf("failed to create jackpot engine (OOM or exception)")
	}
	return &Engine{handle: unsafe.Pointer(h)}, nil
}

// AddJackpot creates a new independent jackpot and returns its ID
func (e *Engine) AddJackpot(cfg Config) (int, error) {
	if atomic.LoadUint64(&e.closed) == 1 {
		return -1, fmt.Errorf("engine already destroyed")
	}

	id := C.jackpot_add((*C.JackpotEngineHandle)(e.handle),
		C.uint64_t(cfg.MinPoint),
		C.uint64_t(cfg.MaxPoint),
		C.int(cfg.Volatility))

	if id == -1 {
		return -1, fmt.Errorf("failed to add jackpot")
	}
	return int(id), nil
}

// ProcessOne convenience single-jackpot version used 99% of the time
func (e *Engine) ProcessOne(id int, step uint64) (bool, error) {
	if atomic.LoadUint64(&e.closed) == 1 {
		return false, fmt.Errorf("engine already destroyed")
	}
	ok := C.jackpot_process((*C.JackpotEngineHandle)(e.handle),
		C.int(id), C.uint64_t(step))
	return ok != 0, nil
}

// Close frees native memory – call on shutdown
func (e *Engine) Close() error {
	if !atomic.CompareAndSwapUint64(&e.closed, 0, 1) {
		return nil
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	if e.handle != nil {
		C.jackpot_engine_destroy((*C.JackpotEngineHandle)(e.handle))
		e.handle = nil
	}
	return nil
}

// GetDropPoint returns the drop point for the given jackpot ID
func (e *Engine) GetDropPoint(id int) (uint64, error) {
	if atomic.LoadUint64(&e.closed) == 1 {
		return 0, fmt.Errorf("engine already destroyed")
	}

	dp := C.jackpot_get_drop_point(
		(*C.JackpotEngineHandle)(e.handle),
		C.int(id),
	)

	return uint64(dp), nil
}

// Version returns the engine version string
func (e *Engine) Version() (string, error) {
	if atomic.LoadUint64(&e.closed) == 1 {
		return "", fmt.Errorf("engine closed")
	}
	cstr := C.jackpot_engine_version()
	if cstr == nil {
		return "", fmt.Errorf("version not available")
	}
	return C.GoString(cstr), nil
}
