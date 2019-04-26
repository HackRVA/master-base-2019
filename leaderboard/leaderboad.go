package leaderboard

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	api "github.com/HackRVA/master-base-2019/baseapi"
	log "github.com/HackRVA/master-base-2019/filelogging"
	"github.com/joho/godotenv"
)

var logger = log.Ger

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

// PostGameData -- sends gameData to the leaderboard
func postGameData(gameData []string) {
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

// FetchScripts -- fetch user's scripts from leaderboard api
func FetchScripts(BadgeID uint16) {
	err := godotenv.Load()
	if err != nil {
		logger.Error().Msg("Error loading .env file")
	}
	uri := os.Getenv("LEADERBOARD_API")

	b := strconv.Itoa(int(BadgeID))

	resp, err := http.Get(uri + "users/" + b + "/scripts")
	// req.Header.Set("Content-Type", "application/json")

	if err != nil {
		logger.Error().Msgf("error fetching user %d scripts", BadgeID)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(body)
}

func sendToLeaderboard(interval *time.Ticker, quit chan struct{}) {
	for {
		select {
		case <-interval.C:
			logger.Debug().Msg("sending data to leaderboard")
			postGameData(api.GetGameData())
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
