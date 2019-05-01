package leaderboard

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	api "github.com/HackRVA/master-base-2019/baseapi"
	log "github.com/HackRVA/master-base-2019/filelogging"
	"github.com/spf13/viper"
)

var logger = log.Ger.With().Str("pkg", "leaderboard").Logger()

// an in browser editor exists on the leaderboard webserver
// users scripts are fetched and queued up for transmitting to the badge

type script struct {
	Content string `json:"content"`
	Name    string `json:"name"`
}

// UserScripts -- stores scripts of user and hash
type UserScripts struct {
	Scripts []string
}

// GameDataResponse -- response from leaderboard when we send GameData
type GameDataResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// PostGameData -- sends gameData to the leaderboard
func postGameData(gameData []string) {
	uri := viper.GetString("leaderBoard_API") + "consume"

	payload := strings.NewReader(`{"data":[` + strings.Join(gameData, ",") + `]}`)

	req, _ := http.NewRequest("POST", uri, payload)

	req.Header.Add("Content-Type", "application/json")

	res, sendErr := http.DefaultClient.Do(req)
	if sendErr != nil {
		logger.Error().Msg("error sending to leaderboard")
	}

	defer res.Body.Close()
	var g GameDataResponse
	body, _ := ioutil.ReadAll(res.Body)

	err := json.Unmarshal(body, &g)
	if err != nil {
		logger.Error().Msg("error unmarshalling Json response from leaderboard")
	}

	if g.Status == "ok" {
		logger.Info().Msg("sent data to leaderboard")
		api.ZeroGameData()
	}
}

// FetchScripts -- fetch user's scripts from leaderboard api
func FetchScripts(BadgeID uint16) {

	uri := viper.GetString("leaderBoard_API")

	b := strconv.Itoa(int(BadgeID))

	resp, err := http.Get(uri + "users/" + b + "/scripts")
	// req.Header.Set("Content-Type", "application/json")

	if err != nil {
		logger.Error().Msgf("error fetching user %d scripts", BadgeID)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(body)
}

func sendToLeaderboard(interval *time.Ticker, quit chan struct{}) {
	for {
		select {
		case <-interval.C:
			logger.Debug().Msg("attempt to send data to leaderboard")
			postGameData(api.StrGameData())
		case <-quit:
			logger.Debug().Msg("stopping routine that sends data to leaderboard.")
			interval.Stop()
			return
		}
	}
}

// StartLeaderboardLoop -- loop to start go routine that sends data to leaderboard
func StartLeaderboardLoop() {
	interval := time.NewTicker(30 * time.Second)
	quit := make(chan struct{})
	go sendToLeaderboard(interval, quit)
}
