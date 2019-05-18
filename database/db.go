package database

import (
	"encoding/json"
	"fmt"
	"time"

	bw "github.com/HackRVA/master-base-2019/badgewrangler"
	log "github.com/HackRVA/master-base-2019/filelogging"
	gm "github.com/HackRVA/master-base-2019/game"
	"github.com/cnf/structhash"
	scribble "github.com/nanobox-io/golang-scribble"
)

var logger = log.Ger.With().Str("pkg", "baseapi").Logger()

var gamesSent = 0

type patientZero struct {
	gameID uint16
}

type zombieGames struct {
	patientZero []patientZero
}

// ScheduleGame -- save gamespec to database
func ScheduleGame(game gm.Game) {
	hash, err := structhash.Hash(game, 1)
	if err != nil {
		logger.Error().Msgf("error scheduling game: %s", err)
	}

	currentTime := time.Now().UTC().UnixNano()

	var bitMask int64

	// Init mask
	for i := 0; i < 16; i++ {
		bitMask++
		bitMask <<= 1
	}

	game.GameID = (uint16)(currentTime * bitMask)

	// create a new scribble database, providing a destination for the database to live
	db, _ := scribble.New("./data", nil)
	logger.Info().Msg("scheduling game")
	db.Write("games", hash, game)
}

// SaveGameData -- save game data to db
func SaveGameData(data *bw.GameData) {
	d := &GameDataWithSent{}
	d.GameData = *data
	d.Sent = false
	hash, err := structhash.Hash(d, 1)
	if err != nil {
		logger.Error().Msgf("error saving game data: %s", err)
	}

	db, _ := scribble.New("./data", nil)
	logger.Info().Msg("saving game data")
	db.Write("game_data", hash, d)
}

func killGameData() {
	db, _ := scribble.New("./data", nil)

	// Delete all fish from the database
	if err := db.Delete("game_data", ""); err != nil {
		fmt.Println("Error", err)
	}
}

// ZeroGameData -- Sets all game data as sent
func ZeroGameData() {
	gameData := GetGameData()

	db, _ := scribble.New("./data", nil)

	for _, g := range gameData {
		g.Sent = true
		logger.Debug().Msgf("zeroing game_data for badgeID: %d", g.BadgeID)
		hash, err := structhash.Hash(g, 1)
		if err != nil {
			logger.Error().Msgf("error saving zeroed game data: %s", err)
		}
		db.Write("game_data", hash, g)
		logger.Debug().Msg("zeroing game data")
	}

}

func notSent(gd []GameDataWithSent, f func(GameDataWithSent) bool) []GameDataWithSent {
	notSent := make([]GameDataWithSent, 0)
	for _, g := range gd {
		if f(g) {
			notSent = append(notSent, g)
		}
	}
	return notSent
}

// StrGameData -- GameData Returned as []String
func StrGameData() []string {
	var s []string

	for _, c := range GetGameData() {
		s = append(s, c.ToString())
	}
	return s
}

// GetGameData -- retrieves gamedata from the db
func GetGameData() []GameDataWithSent {
	// create a new scribble database, providing a destination for the database to live
	db, _ := scribble.New("./data", nil)

	// Read more games from the database
	records, _ := db.ReadAll("game_data")
	games := []GameDataWithSent{}

	for _, g := range records {
		gameFound := GameDataWithSent{}
		if err := json.Unmarshal([]byte(g), &gameFound); err != nil {
			fmt.Println("Error", err)
		}
		games = append(games, gameFound)
	}

	pendingData := notSent(games, func(g GameDataWithSent) bool {
		return g.Sent == false
	})

	return pendingData
}

func determineTeam(variant uint8, gameID uint16) uint8 {
	gamesSent++
	switch variant {
	case 0:
		// "FREE FOR ALL",
		return 1
	case 1:
		// "TEAM BATTLE",
		return uint8(gamesSent%2 + 1)
	case 2:
		// "ZOMBIES!",
		// TEAM 1 is zombie
		// TEAM 2 is non-zombie
		var z zombieGames
		for _, c := range z.patientZero {
			if c.gameID == gameID {
				return 2
			}
		}
		pZero := &patientZero{
			gameID: gameID,
		}

		z.patientZero = append(z.patientZero, *pZero)

		return 1
	case 3:
		// "CAPTURE BADGE",
		return 1
	}

	return 1 // default return -- we should not use a zero value for team
}

// GetNext -- return the next game
func GetNext() gm.Game {
	t := time.Now()
	var g gm.Game
	g.AbsStart = 0

	games := GetGames()
	for _, game := range games {
		logger.Debug().Msgf(
			"now: %d nextGame: %d in the future: %t",
			int64(t.Unix()),
			game.AbsStart,
			int64(t.Unix()) < game.AbsStart+int64(game.Duration))

		// return the first game that is greater than now
		if int64(t.Unix()) < game.AbsStart+int64(game.Duration) {
			game.Team = determineTeam(game.Variant, game.GameID)
			game.StartTime = int16(game.AbsStart - t.Unix())
			return game
		}
	}

	return g
}

// GetGames -- retrieves games from DB
func GetGames() []gm.Game {
	// create a new scribble database, providing a destination for the database to live
	db, _ := scribble.New("./data", nil)
	// Read more games from the database
	moregames, _ := db.ReadAll("games")
	// iterate over moregames creating a new game for each record
	games := []gm.Game{}
	for _, game := range moregames {
		g := gm.Game{}
		json.Unmarshal([]byte(game), &g)
		games = append(games, g)
	}

	return games
}

// DataInGameOut - stores the game data and gets the current/next game
func DataInGameOut(gameDataIn chan *bw.GameData, gameDataOut chan *bw.GameData, gameOut chan *gm.Game) {
	for {
		gameData := <-gameDataIn
		logger.Debug().Msg("DataInGameOut received gameData from GameDataIn channel")
		fmt.Println(gameData.GameID)
		nextGame := GetNext()
		gameOut <- &nextGame
		logger.Debug().Msg("DataInGameOut sent nextGame to gameOut channel")
		SaveGameData(gameData)
	}
}
