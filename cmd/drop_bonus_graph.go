package main

import (
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/wcharczuk/go-chart/v2"

	"gitlabce.champion.tm/champion-prj/mersenne-twister-shared-object/mtwrapper"
)

func main() {
	rangeGens := []uint64{
		1, 2, 3, 5, 10, 20, 50, 100,
		150, 200, 300, 500, 700, 1000,
		5000, 10000, 50000, 100000,
	}

	const iterations = 200_000

	type Result struct {
		rg       float64
		actual   float64
		expected float64
	}

	var results []Result

	// ---- CSV ----
	csvFilePath := "./stats/drop_bonus_stats.csv"
	absCsvPath, _ := filepath.Abs(csvFilePath)

	f, _ := os.Create(absCsvPath)
	w := csv.NewWriter(f)
	w.Write([]string{"rangeGen", "actualDropRate", "expectedDropRate"})

	for _, rg := range rangeGens {
		drops := 0

		for i := 0; i < iterations; i++ {
			if mtwrapper.DropBonus(rg) {
				drops++
			}
		}

		actual := float64(drops) / float64(iterations)
		expected := 1.0 / float64(rg)

		results = append(results, Result{
			rg:       float64(rg),
			actual:   actual,
			expected: expected,
		})

		w.Write([]string{
			strconv.FormatUint(rg, 10),
			strconv.FormatFloat(actual, 'f', 8, 64),
			strconv.FormatFloat(expected, 'f', 8, 64),
		})

		log.Printf("rangeGen=%d actual=%f expected=%f", rg, actual, expected)
	}

	w.Flush()
	f.Close()
	log.Println("CSV saved:", absCsvPath)

	// ---- GRAPH ----

	x := []float64{}
	yActual := []float64{}
	yExpected := []float64{}

	var minY, maxY float64 = 1.0, 0.0

	for _, r := range results {
		x = append(x, r.rg)
		yActual = append(yActual, r.actual)
		yExpected = append(yExpected, r.expected)

		if r.actual < minY {
			minY = r.actual
		}
		if r.expected < minY {
			minY = r.expected
		}
		if r.actual > maxY {
			maxY = r.actual
		}
		if r.expected > maxY {
			maxY = r.expected
		}
	}

	// Добавляем padding чтобы график не "прилипал"
	padding := (maxY - minY) * 0.15
	if padding == 0 {
		padding = 0.001
	}

	graph := chart.Chart{
		Title:  "DropBonus Probability (Actual vs Expected)",
		Width:  2200,
		Height: 900,
		XAxis: chart.XAxis{
			Name:  "rangeGen",
			Range: &chart.LogarithmicRange{},
		},
		YAxis: chart.YAxis{
			Name: "Probability",
			Range: &chart.ContinuousRange{
				Min: minY - padding,
				Max: maxY + padding,
			},
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				Name:    "Actual Drop Rate",
				XValues: x,
				YValues: yActual,
				Style: chart.Style{
					StrokeColor: chart.ColorBlue,
					StrokeWidth: 3,
				},
			},
			chart.ContinuousSeries{
				Name:    "Expected 1/rangeGen",
				XValues: x,
				YValues: yExpected,
				Style: chart.Style{
					StrokeColor: chart.ColorRed,
					StrokeWidth: 3,
				},
			},
		},
	}

	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}

	pngPath, _ := filepath.Abs("./stats/drop_bonus_chart.png")
	out, _ := os.Create(pngPath)
	defer out.Close()

	graph.Render(chart.PNG, out)
	log.Println("Graph saved:", pngPath)
}
