package jackpot_engine

/*
#cgo CFLAGS: -Iv1/src
#cgo LDFLAGS: -L${SRCDIR}/v1/lib -ljackpot_engine_v1_linux
#cgo LDFLAGS: -L${SRCDIR}/v1/lib -ljackpot_engine_v1_macos

#include "jackpot_engine.h"
*/
import "C"

import (
    "context"
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

type Config struct {
	MinPoint    uint64
	MaxPoint    uint64
	Volatility  Volatility
}

// Engine is safe for concurrent use from thousands of goroutines
type Engine struct {
	handle unsafe.Pointer // *C.JackpotEngineHandle
	closed uint64          // atomic
	mu     sync.Mutex      // protects destroy sequence
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

// Process advances one or more jackpots in a single atomic step.
// Returns slice of jackpot IDs that dropped this step.
func (e *Engine) Process(ctx context.Context, increments map[int]uint64) ([]int, error) {
	if atomic.LoadUint64(&e.closed) == 1 {
		return nil, fmt.Errorf("engine already destroyed")
	}

	// Fast path: single increment (most common)
	if len(increments) == 1 {
		for id, amt := range increments {
			ok := C.jackpot_process((*C.JackpotEngineHandle)(e.handle),
				C.int(id), C.uint64_t(amt))
			if ok != 0 {
				return []int{id}, nil
			}
			return nil, nil
		}
	}

	// General path – lock once, process all
	e.mu.Lock()
	defer e.mu.Unlock()

	var dropped []int
	for id, amt := range increments {
		select {
		case <-ctx.Done():
			return dropped, ctx.Err()
		default:
		}

		ok := C.jackpot_process((*C.JackpotEngineHandle)(e.handle),
			C.int(id), C.uint64_t(amt))
		if ok != 0 {
			dropped = append(dropped, id)
		}
	}
	return dropped, nil
}

// ProcessOne convenience single-jackpot version used 99% of the time
func (e *Engine) ProcessOne(id int, amount uint64) (bool, error) {
	if atomic.LoadUint64(&e.closed) == 1 {
		return false, fmt.Errorf("engine already destroyed")
	}
	ok := C.jackpot_process((*C.JackpotEngineHandle)(e.handle),
		C.int(id), C.uint64_t(amount))
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
