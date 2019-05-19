package webapi

import (
	"encoding/json"

	bw "github.com/HackRVA/master-base-2019/badgewrangler"
)

// SendGameData - Send GameData to Leaderboard
func SendGameData(gameDataIn chan *bw.GameData) {
	for {
		gameData := <-gameDataIn

		gameDataJSON, _ := json.Marshal(gameData)
		logger.Debug().Msg("Send Game Data: " + string(gameDataJSON))
	}
}
