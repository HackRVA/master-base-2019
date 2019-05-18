package sync

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	db "github.com/HackRVA/master-base-2019/database"
	log "github.com/HackRVA/master-base-2019/filelogging"
	gm "github.com/HackRVA/master-base-2019/game"
	"github.com/spf13/viper"
)

var logger = log.Ger.With().Str("pkg", "leaderboard").Logger()

func closeResponse(res *http.Response) {
	if res != nil {
		res.Body.Close()
	}
}

// Fetches scheduled games from the MASTER master base station
func fetchScheduledGames() []gm.Game {
	uri := viper.GetString("master_URL") + "/api/games"

	resp, _ := http.Get(uri)

	defer closeResponse(resp)

	body, _ := ioutil.ReadAll(resp.Body)
	var gms []gm.Game

	jsonErr := json.Unmarshal(body, &gms)
	if jsonErr != nil {
		logger.Error().Msg("could not get schedule game from MASTER master base station")
	}

	return gms
}

// isInLocalDB -- returns true if this game already exists in the localDB
func isInLocalDB(game gm.Game) bool {
	dbGames := db.GetGames()

	for _, g := range dbGames {
		if game.GameID == g.GameID {
			return true
		}
	}

	return false
}

// pushToLocalDB -- grabs future games and schedules them in local database
func pushToLocalDB(games []gm.Game) {
	t := time.Now()

	for _, game := range games {
		if int64(t.Unix()) < game.AbsStart+int64(game.Duration) {
			if !isInLocalDB(game) {
				db.ScheduleGame(game)
			}
		}
	}
}

func fetchGames(interval *time.Ticker, quit chan struct{}) {
	for {
		select {
		case <-interval.C:
			logger.Debug().Msg("attempt to send data to leaderboard")
			pushToLocalDB(fetchScheduledGames())
		case <-quit:
			logger.Debug().Msg("stopping routine that sends data to leaderboard.")
			interval.Stop()
			return
		}
	}
}

// StartSyncLoop -- starts a go routine to fetch games from MASTER master base station
// this runs on an interval
func StartSyncLoop() {
	interval := time.NewTicker(30 * time.Second)
	quit := make(chan struct{})
	go fetchGames(interval, quit)
}
