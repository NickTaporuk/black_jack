package main

import (
	"fmt"

	"gitlabce.champion.tm/champion-prj/mersenne-twister-shared-object/jackpot_engine"
)

func main() {
	const masterSeed uint64 = 0xC0FFEE1234567890
	seed := jackpot_engine.DeterministicSeed(masterSeed)
	jpe, err := jackpot_engine.New(seed)
	if err != nil {
		panic(err)
	}

	cfg := jackpot_engine.Config{
		MinPoint:   11000,
		MaxPoint:   12000,
		Volatility: jackpot_engine.VolatilityHigh,
	}
	fmt.Printf("Jackpot Engine initialized: %+v\n", jpe)
	fmt.Printf("Jackpot Config: %+v\n", cfg)
	jackpotID, err := jpe.AddJackpot(cfg)
	if err != nil {
		panic(err)
	}

	var droppedCount uint64
	var undroppedCount uint64
	for i := 0; i < 5000000; i++ {
		dropped, err := jpe.ProcessOne(jackpotID, 100)
		if err != nil {
			panic(err)
		}
		if dropped {
			droppedCount++
		} else {
			undroppedCount++
		}
	}
	fmt.Printf("After 5,000,000 steps:\n")
	fmt.Printf("Dropped count: %d\n", (droppedCount*100)/5000000)
	fmt.Printf("Undropped count: %d\n", (undroppedCount*100)/5000000)

	defer func() {
		err := jpe.Close()
		if err != nil {
			panic(err)
		}
	}()
}
