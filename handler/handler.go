package handler

import (
	"fmt"

	msg "github.com/HackRVA/master-base-2019/messages"
)

// go routine -- takes in game data -
// unwraps them and stores them and passes back --
// return game spec packet onto another channel

// Receive - takes in game data,
// stores it and decides whether or not to respond
func Receive(data *msg.GameData, spec *msg.GameSpec) {
	store(data)
	Send(data, spec)
}

// channel in - game data
func store(data *msg.GameData) {
	fmt.Println(
		"storing game data :::",
		"BadgeID: ", data.BadgeID,
		"GameID: ", data.GameID)
}

// Send -- responds to badge with gamespec
func Send(data *msg.GameData, spec *msg.GameSpec) {
	fmt.Println("responding to: ", data.BadgeID)
	fmt.Println(
		"BadgeID: ", spec.BadgeID,
		"StartTime: ", spec.StartTime,
		"Duration: ", spec.Duration,
		"Variant: ", spec.Variant,
		"Team: ", spec.Team,
		"GameID: ", spec.GameID)
}
