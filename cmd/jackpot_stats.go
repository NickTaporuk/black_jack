package main

import (
	"encoding/json"
	"fmt"
	"os"

	"gitlabce.champion.tm/champion-prj/mersenne-twister-shared-object/jackpot_engine"
)

type StatRow struct {
	Volatility int    `json:"volatility"`
	Step       uint64 `json:"step"`
	HitIndex   int    `json:"hit_index"`
	Iterations uint64 `json:"iterations"`
	DropPoint  uint64 `json:"drop_point"`
}

func main() {
	seed := jackpot_engine.DeterministicSeed(123456789)
	jpe, err := jackpot_engine.New(seed)
	if err != nil {
		panic(err)
	}
	defer jpe.Close()

	minPoint := uint64(10000)
	maxPoint := uint64(12000)

	steps := []uint64{1, 5, 10, 50, 100}
	vols := []jackpot_engine.Volatility{
		jackpot_engine.VolatilityLow,
		jackpot_engine.VolatilityMedium,
		jackpot_engine.VolatilityHigh,
	}

	var stats []StatRow
	const hitsToCollect = 200

	for _, vol := range vols {
		for _, step := range steps {

			cfg := jackpot_engine.Config{
				MinPoint:   minPoint,
				MaxPoint:   maxPoint,
				Volatility: vol,
			}

			jpID, err := jpe.AddJackpot(cfg)
			if err != nil {
				panic(err)
			}

			for hitIndex := 1; hitIndex <= hitsToCollect; hitIndex++ {
				var iterations uint64

				for {
					drop, err := jpe.ProcessOne(jpID, step)
					if err != nil {
						panic(err)
					}
					iterations++

					if drop {
						dropPoint, err := jpe.GetDropPoint(jpID)
						if err != nil {
							panic(err)
						}
						row := StatRow{
							Volatility: int(vol),
							Step:       step,
							HitIndex:   hitIndex,
							Iterations: iterations,
							DropPoint:  dropPoint,
						}
						stats = append(stats, row)
						break
					}
				}
			}
		}
	}

	file, _ := os.Create("stats.json")
	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	enc.Encode(stats)
	file.Close()

	fmt.Println("✔ Saved stats.json")
}
