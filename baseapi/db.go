package baseapi

import (
	"encoding/json"
	"fmt"
	"time"

	bw "github.com/HackRVA/master-base-2019/badgewrangler"
	gm "github.com/HackRVA/master-base-2019/game"
	"github.com/cnf/structhash"
	scribble "github.com/nanobox-io/golang-scribble"
)

// ScheduleGame -- save gamespec to database
func ScheduleGame(game gm.Game) {
	hash, err := structhash.Hash(game, 1)
	if err != nil {
		logger.Error().Msgf("error scheduling game: %s", err)
	}

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

// GetNext -- return the next game
func GetNext() gm.Game {
	t := time.Now()
	var g gm.Game
	g.AbsStart = 0

	games := GetGames()
	for _, game := range games {
		logger.Debug().Msgf(
			"now: %d \nnextGame: %d \nin the future: %t",
			int64(t.Unix()),
			game.AbsStart,
			int64(t.Unix()) < game.AbsStart+int64(game.Duration))

		// return the first game that is greater than now
		if int64(t.Unix()) < game.AbsStart+int64(game.Duration) {
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
		fmt.Println(gameData.GameID)
		//gameDataOut <- gameData
		nextGame := GetNext()
		gameOut <- &nextGame
		SaveGameData(gameData)
	}
}
