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
	hash, err := structhash.Hash(data, 1)
	if err != nil {
		logger.Error().Msgf("error saving game data: %s", err)
	}

	db, _ := scribble.New("./data", nil)
	logger.Info().Msg("saving game data")
	db.Write("game_data", hash, data)
}

// GetGameData -- retrieves gamedata from the db
func GetGameData() []string {
	// create a new scribble database, providing a destination for the database to live
	db, _ := scribble.New("./data", nil)
	// Read more games from the database
	gameData, _ := db.ReadAll("game_data")
	return gameData
}

// GetNext -- return the next game
func GetNext() gm.Game {
	t := time.Now()
	var g gm.Game
	g.AbsStart = 0

	games := GetGames()
	for _, game := range games {
		fmt.Println(
			"now: ",
			int64(t.Unix()),
			"\nnextGame: ",
			game.AbsStart,
			"\nin the future:",
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

func GetInfo() gm.GameInfo {
     db, _ := scribble.New("./info", nil)
     gameInfo, _ := db.Read("info", "info", nil)
     
     return gameInfo
}

func updateGameInfo(gameInfo GameInfo) {

     // Establish driver connection
     db, err = scribble.New("./info", nill)

     if err != nil {
     	logger.Error().Msgf("Unable to: %s", err)
     	writeGameInfo(db, gameInfo)
	return 
     }

     // retreive a single entry from info
     info, err = db.Read("info", "info", nil)

     if err != nil {
     	logger.Warn().Msgf("game info update failed: %s", err)
     	writeGameInfo(db, gameInfo)
	return 
     }
     
     // Backup entry
     infoOld, err = db.Read("info","old",nil)

     if err != nil {
     	logger.Warn().Msgf("Reading game info backup failed: %s", err)
     }
     
     err = db.Delete("info","old",nil)

     if err != nil {
     	logger.Warn().Msgf("Deleting backup failed: %s", err)
     }

     err = db.Write("info", "old", info)
     if err != nil {
     	logger.Error().Msgf("Writing Backup failed: %s", err)
	return
     }
     
     err = db.Delete("info","info",nil)

     if err != nil {
     	logger.Warn().Msgf("Deleting game info entry failed: %s", err)
     }

     err = db.Write("info", "info", info)
     if err != nil {
     	logger.Error().Msgf("Writing game info entry failed: %s", err)
     }
}