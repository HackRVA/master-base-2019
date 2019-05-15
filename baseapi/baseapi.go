package baseapi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	log "github.com/HackRVA/master-base-2019/filelogging"
	gm "github.com/HackRVA/master-base-2019/game"
	gi "github.com/HackRVA/master-base-2019/gameinfo"
	mux "github.com/gorilla/mux"
)

var logger = log.Ger.With().Str("pkg", "baseapi").Logger()

// NewGame - function to schedule newgame
func NewGame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var e gm.Game
	currentTime := time.Now().UTC().UnixNano()

	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &e)

	var bitMask int64

	// Init mask
	for i := 0; i < 16; i++ {
		bitMask++
		bitMask <<= 1
	}

	e.GameID = (uint16)(currentTime * bitMask)

	ScheduleGame(e)
	AddNewGameEntryToGameInfo(e)

	j, _ := json.Marshal(e)
	w.Write(j)
}

// NextGame -- returns the game that is sheduled next
func NextGame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	next := GetNext()
	var gameInfo gi.GameInfo
	AddNewGameEntryToGameInfo(next)

	gameInfo = GetInfo(next.GameID)

	logger.Info().Msgf("gameInfo ID: %d", next.GameID)
	if gameInfo.ID == 0 {
		AddNewGameEntryToGameInfo(next)
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
	j, _ := json.Marshal(GetGames())
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

	gameInfo = GetInfo(uint16(gameIDFromRequest))
	if gameInfo.ID == 0 {
		AddNewGameEntryToGameInfo(game)
		gameInfo = GetInfo(game.GameID)
	}
	j, _ := json.Marshal(gameInfo)
	w.Write(j)
}

// AllInfo -- handler that returns all GameInfo
func AllInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	j, _ := json.Marshal(GetAllInfo())
	w.Write(j)
}
