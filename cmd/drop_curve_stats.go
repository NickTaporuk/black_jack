package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gitlabce.champion.tm/champion-prj/mersenne-twister-shared-object/jackpot_engine"
)

type Result struct {
	DropPoint   uint64  `json:"drop_point"`
	AvgSteps    float64 `json:"avg_steps_to_drop"`
	StdDev      float64 `json:"std_dev"` // опционально, можно убрать
	Volatility  string  `json:"volatility"`
	Runs        int     `json:"runs"`
	GeneratedAt string  `json:"generated_at"`
}

func main() {
	// === НАСТРОЙКИ ===
	const runsPerPoint = 1000                    // точность ±0.3%
	const step = 1                               // как в игре
	fromDropPoint := uint64(6600)
	toDropPoint := uint64(9800)
	vol := jackpot_engine.VolatilityHigh        // меняй на Low/Medium для сравнения

	// === ПОДГОТОВКА ===
	results := make([]Result, 0, toDropPoint-fromDropPoint+1)

	fmt.Printf("Запуск симуляции dropPoint [%d..%d], %d прогонов на точку, volatility=%s\n",
		fromDropPoint, toDropPoint, runsPerPoint, jackpot_engine.VolatilityIntToName[vol])

	for dp := fromDropPoint; dp <= toDropPoint; dp++ {
		var sumSteps uint64 = 0
		var sumSquares uint64 = 0

		for r := uint64(0); r < runsPerPoint; r++ {
			seed := jackpot_engine.DeterministicSeed(dp*100000 + r)
			eng, err := jackpot_engine.New(seed)
			if err != nil {
				panic(err)
			}

			cfg := jackpot_engine.Config{
				MinPoint:   0,
				MaxPoint:   10_000_000_000, // не влияет
				Volatility: vol,
			}

			id, err := eng.AddJackpot(cfg)
			if err != nil {
				panic(err)
			}

			// ВРУЧНУЮ ЗАДАЁМ dropPoint!
			if err := eng.SetDropPoint(id, dp); err != nil {
				panic(err)
			}

			steps := uint64(0)
			for {
				steps += step
				dropped, err := eng.ProcessOne(id, step)
				if err != nil {
					panic(err)
				}
				if dropped {
					sumSteps += steps
					sumSquares += steps * steps
					break
				}
			}
			eng.Close()
		}

		avg := float64(sumSteps) / float64(runsPerPoint)
		variance := float64(sumSquares)/float64(runsPerPoint) - avg*avg
		stdDev := 0.0
		if variance > 0 {
			stdDev = variance
		}

		results = append(results, Result{
			DropPoint:   dp,
			AvgSteps:    avg,
			StdDev:      stdDev,
			Volatility:  jackpot_engine.VolatilityIntToName[vol],
			Runs:        runsPerPoint,
			GeneratedAt: time.Now().Format(time.RFC3339),
		})

		if dp%50 == 0 || dp == toDropPoint {
			fmt.Printf("dp=%d → avg=%.2f (expected ≈%.0f)\n", dp, avg, float64(dp))
		}
	}

	// === СОХРАНЕНИЕ ===
	filename := fmt.Sprintf("drop_curve_%s_%d-%d.json",
		jackpot_engine.VolatilityIntToName[vol], fromDropPoint, toDropPoint)

	dir := "./jackpot_engine/stats"
	os.MkdirAll(dir, 0755)
	path := filepath.Join(dir, filename)

	file, _ := os.Create(path)
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	enc.Encode(results)

	fmt.Printf("\nГОТОВО! Файл сохранён:\n   %s\n", path)
	fmt.Printf("Теперь запусти Python-скрипт → получишь график уровня Aristocrat!\n")
}