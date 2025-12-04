package main

import (
	"encoding/json"
	"fmt"
	"os"

	"gitlabce.champion.tm/champion-prj/mersenne-twister-shared-object/jackpot_engine"
)

type Sample struct {
	Progress float64 `json:"progress"` // counter / dropPoint
	Hit      int     `json:"hit"`      // 0 or 1
}

type VolatilityReport struct {
	Volatility string   `json:"volatility"`
	Samples    []Sample `json:"samples"`
}

type Report struct {
	MinPoint uint64             `json:"min_point"`
	MaxPoint uint64             `json:"max_point"`
	Reports  []VolatilityReport `json:"reports"`
}

func runSimulation(vol jackpot_engine.Volatility, min, max uint64) VolatilityReport {
	seed := jackpot_engine.DeterministicSeed(999123)
	engine, _ := jackpot_engine.New(seed)

	cfg := jackpot_engine.Config{
		MinPoint:   min,
		MaxPoint:   max,
		Volatility: vol,
	}
	id, _ := engine.AddJackpot(cfg)

	drop, err := engine.GetDropPoint(id)
	if err != nil {
		panic(err)
	}
	if drop == 0 {
		panic("dropPoint=0 => wrong config")
	}

	samples := make([]Sample, 0, drop)

	for i := uint64(0); i < drop; i++ {
		hit, _ := engine.ProcessOne(id, 1)
		progress := float64(i) / float64(drop)
		samples = append(samples, Sample{
			Progress: progress,
			Hit:      boolToInt(hit),
		})
		if hit {
			break
		}
	}

	engine.Close()

	vname := map[jackpot_engine.Volatility]string{
		jackpot_engine.VolatilityLow:    "Low",
		jackpot_engine.VolatilityMedium: "Medium",
		jackpot_engine.VolatilityHigh:   "High",
	}[vol]

	return VolatilityReport{
		Volatility: vname,
		Samples:    samples,
	}
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func main() {
	min := uint64(10000)
	max := uint64(20000)

	report := Report{
		MinPoint: min,
		MaxPoint: max,
		Reports: []VolatilityReport{
			runSimulation(jackpot_engine.VolatilityLow, min, max),
			runSimulation(jackpot_engine.VolatilityMedium, min, max),
			runSimulation(jackpot_engine.VolatilityHigh, min, max),
		},
	}

	f, err := os.Create("volatility_stats.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	enc.Encode(report)

	fmt.Println("Saved → volatility_stats.json")
}
