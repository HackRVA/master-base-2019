package database

import (
	"encoding/json"

	bw "github.com/HackRVA/master-base-2019/badgewrangler"
)

// GameDataWithSent -- extending GameData locally to easily unmarshal
type GameDataWithSent struct {
	bw.GameData
	Sent bool
}

// MarshalBinary -
func (g *GameDataWithSent) MarshalBinary() ([]byte, error) {
	return json.Marshal(g)
}

// UnmarshalBinary -
func (g GameDataWithSent) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &g); err != nil {
		return err
	}

	return nil
}

// ToString -- takes in GameData returns string
func (g *GameDataWithSent) ToString() string {
	b, _ := g.MarshalBinary()

	return string(b)
}
