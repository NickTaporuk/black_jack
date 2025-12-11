package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gitlabce.champion.tm/champion-prj/mersenne-twister-shared-object/jackpot_engine"
)

type Interval struct {
	Min uint64 `json:"min"`
	Max uint64 `json:"max"`
}

type ResultEntry struct {
	DropPoint uint64  `json:"drop_point"`
	AvgSteps  float64 `json:"avg_steps"`
}

type IntervalReport struct {
	Interval Interval      `json:"interval"`
	Results  []ResultEntry `json:"results"`
}

type FullReport struct {
	Volatility string           `json:"volatility"`
	Step       uint64           `json:"step"`
	Runs       int              `json:"runs"`
	Data       []IntervalReport `json:"data"`
}

func main() {
	intervals := []Interval{
		// 		{10000, 12000},
		//{6600, 9800},
		//{36700, 45000},
		{265000, 310000},
	}

	const runsPerPoint = 500
	const step = 50
	vol := jackpot_engine.VolatilityHigh
	// 	vol := jackpot_engine.VolatilityMedium
	//vol := jackpot_engine.VolatilityLow

	report := FullReport{
		Volatility: jackpot_engine.VolatilityIntToName[vol],
		Step:       step,
		Runs:       runsPerPoint,
	}

	for _, iv := range intervals {

		intReport := IntervalReport{
			Interval: iv,
		}

		for dp := iv.Min; dp <= iv.Max; dp++ {
			var total uint64
			seed := jackpot_engine.DeterministicSeed(uint64(1111))
			eng, err := jackpot_engine.New(seed)
			if err != nil {
				panic(err)
			}

			cfg := jackpot_engine.Config{
				MinPoint:   dp,
				MaxPoint:   iv.Max,
				Volatility: vol,
			}

			id, err := eng.AddJackpot(cfg)
			if err != nil {
				panic(err)
			}
			for i := 0; i < runsPerPoint; i++ {
				var steps uint64
				drop, err := eng.ProcessOne(id, step)
				if err != nil {
					panic(err)
				}
				steps++

				if drop {
					total++
				}
			}
			eng.Close()
			avg := float64(total) / float64(runsPerPoint)
			fmt.Printf("Interval %d-%d | dp=%d avg=%f\n",
				iv.Min, iv.Max, dp, avg)

			intReport.Results = append(intReport.Results, ResultEntry{
				DropPoint: dp,
				AvgSteps:  avg,
			})
			iv.Min++
		}

		report.Data = append(report.Data, intReport)
	}

	filename := fmt.Sprintf("jackpot_curve_stats_%d.json", time.Now().Unix())
	directory := "./jackpot_engine/stats/"
	absPath, err := filepath.Abs(directory)
	if err != nil {
		panic(err)
	}
	filePath := filepath.Join(absPath, filename)
	f, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "    ")
	enc.Encode(report)

	fmt.Printf("Saved %s\n", filePath)
}
