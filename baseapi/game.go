package baseapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	bw "github.com/HackRVA/master-base-2019/badgewrangler"
	"github.com/joho/godotenv"
)

// PostGameData -- sends gameData to the leaderboard
func PostGameData(gameData []string) {
	err := godotenv.Load()
	if err != nil {
		logger.Error().Msg("Error loading .env file")
	}
	uri := os.Getenv("LEADERBOARD_API")
	json := `{"data":[` + strings.Join(gameData, ",") + `]}`
	var jsonStr = []byte(json)
	logger.Info().Msg(json)

	req, err := http.NewRequest("POST", uri+"consume", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error().Msgf("error connecting to Leaderboard: %s", err)
	}

	logger.Info().Msg("sending data to leaderboard")
	defer resp.Body.Close()
}

// SendGameData - Send GameData to Leaderboard
func SendGameData(gameDataIn chan *bw.GameData) {
	for {
		gameData := <-gameDataIn

		gameDataJSON, _ := json.Marshal(gameData)
		logger.Debug().Msg("Send Game Data: " + string(gameDataJSON))
	}

}
