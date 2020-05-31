package main

import (
	"fmt"
)

func main() {
	gamemode := []string{
		"HUMAN",
		"DEVIL HUNTER",
		"SON OF SPARDA",
		"HEAVEN OR HELL",
		"LEGENDARY DARK KNIGHT",
		"DANTE MUST DIE",
		"HELL AND HELL",
		"BLOODY PALACE",
	}

	activity := NewActivity(10, gamemode)
	result := activity.ChooseCommit("choose prev commit")
	fmt.Println(gamemode[result])
}
