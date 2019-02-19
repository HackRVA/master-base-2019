package main

import (
	"fmt"
	"io/ioutil"

	"github.com/HackRVA/master-base-2019/handler"
)

func printschedule(games []byte) {
	for _, game := range handler.DecodeJSON(games) {
		fmt.Println("StartTime: ", game.StartTime)
		fmt.Println("Duration: ", game.Duration)
		fmt.Println("Variant: ", game.Variant)
		fmt.Println("Team: ", game.Team)
		fmt.Println("GameID: ", game.GameID)
		fmt.Println("-------------------")
	}
}

func main() {
	games, _ := ioutil.ReadFile("games.json")
	printschedule(games)
}
