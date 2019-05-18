package baseapi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	db "github.com/HackRVA/master-base-2019/database"
	log "github.com/HackRVA/master-base-2019/filelogging"
	gm "github.com/HackRVA/master-base-2019/game"
	gi "github.com/HackRVA/master-base-2019/gameinfo"
	info "github.com/HackRVA/master-base-2019/info"
	mux "github.com/gorilla/mux"
)

var logger = log.Ger.With().Str("pkg", "baseapi").Logger()

// NewGame - function to schedule newgame
func NewGame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var e gm.Game

	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &e)

	db.ScheduleGame(e)
	info.AddNewGameEntryToGameInfo(e)

	j, _ := json.Marshal(e)
	w.Write(j)
}

// NextGame -- returns the game that is sheduled next
func NextGame(w http.ResponseWriter, r *http.Request) {
	t := time.Now()
	w.Header().Set("Content-Type", "application/json")
	next := func() gm.Game {
		var g gm.Game
		g.AbsStart = 0
		games := db.GetGames()
		for _, game := range games {
			// return the first game that is greater than now
			if int64(t.Unix()) < game.AbsStart+int64(game.Duration) {
				game.StartTime = int16(game.AbsStart - t.Unix())
				return game
			}
		}
		return g
	}()

	var gameInfo gi.GameInfo
	info.AddNewGameEntryToGameInfo(next)

	gameInfo = info.GetInfo(next.GameID)

	logger.Info().Msgf("gameInfo ID: %d", next.GameID)
	if gameInfo.ID == 0 {
		info.AddNewGameEntryToGameInfo(next)
	}

	if next.AbsStart == 0 {
		j, _ := json.Marshal("There are no games scheduled")
		w.Write(j)
	} else {
		j, _ := json.Marshal(next)
		w.Write(j)
	}
}

// AllGames - returns all scheduled games
func AllGames(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	j, _ := json.Marshal(db.GetGames())
	w.Write(j)
}

// Info -- Handler that serves up gameInfo
func Info(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var game gm.Game
	var gameInfo gi.GameInfo

	params := mux.Vars(r)
	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &game)

	gameIDFromRequest, err := strconv.ParseUint(params["id"], 10, 16)

	if err != nil {
		logger.Error().Msgf("Info: %s", err)
	}

	gameInfo = info.GetInfo(uint16(gameIDFromRequest))
	if gameInfo.ID == 0 {
		info.AddNewGameEntryToGameInfo(game)
		gameInfo = info.GetInfo(game.GameID)
	}
	j, _ := json.Marshal(gameInfo)
	w.Write(j)
}

// AllInfo -- handler that returns all GameInfo
func AllInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	j, _ := json.Marshal(info.GetAllInfo())
	w.Write(j)
}
