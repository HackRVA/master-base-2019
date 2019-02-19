package handler

import (
	"encoding/json"
	"log"
)

// Game - struct for our json games array
type Game []struct {
	StartTime string `json:"StartTime"`
	Duration  string `json:"Duration"`
	Variant   string `json:"Variant"`
	Team      string `json:"Team"`
	GameID    string `json:"GameID"`
}

// DecodeJSON - decoding json data
func DecodeJSON(data []byte) Game {
	games := Game{}
	err := json.Unmarshal(data, &games)
	if err != nil {
		log.Fatal(err)
	}
	return games
}
