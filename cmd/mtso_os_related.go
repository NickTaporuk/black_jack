package main

import (
	"fmt"

	"gitlabce.champion.tm/champion-prj/mersenne-twister-shared-object/mtwrapper"
)

func main() {
	//fmt.Println("Random Token:", mtwrapper.GenerateRandomToken())
	if mtwrapper.DropBonus(100) {
		fmt.Println("🎉 Bonus Dropped!")
	} else {
		fmt.Println("❌ No Bonus This Time.")
	}
}
