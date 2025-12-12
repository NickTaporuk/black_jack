package main

import (
	"time"

	jackpot_engine_v2 "gitlabce.champion.tm/champion-prj/mersenne-twister-shared-object/jackpot/engine/v2"
)

func main() {

	for i := 0; i < 1000000; i++ {
		seed := time.Now().UnixNano()

		dropped, dropPoint, err := jackpot_engine_v2.CheckUniform(
			uint64(seed),
			10_000,
			12_000,
			11111,
		)

		println("Seed:", seed)
		println("Dropped:", dropped)
		println("dropPoint:", dropPoint)
		println("err:", err)

	}
}
