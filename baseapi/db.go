package baseapi

import (
	"encoding/json"
	"fmt"
	"time"
	"strconv"

	bw "github.com/HackRVA/master-base-2019/badgewrangler"
	gm "github.com/HackRVA/master-base-2019/game"
	gi "github.com/HackRVA/master-base-2019/gameinfo"
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


// Overwrites the 'info' entry with the gameInfo 
func writeGameInfo(db *scribble.Driver, gameInfo gi.GameInfo) {

     if err := db.Write("info", strconv.FormatInt(int64(gameInfo.ID), 10), gameInfo); err != nil {
     	logger.Error().Msgf("error writing to the database: %s", err)
     }
}

// GetInfo -- retreive game info from database
func GetInfo(gameID uint16) gi.GameInfo {
     var gameInfo gi.GameInfo

     db, _ := scribble.New("./info", nil)
     err := db.Read("info", strconv.FormatUint(uint64(gameID), 10), &gameInfo)
     if err != nil {
     	logger.Error().Msgf("Error reading: %s", err)
     }
     
     return gameInfo
}

// GetOldInfo -- retreive old game info from database
func GetOldInfo(gameID uint16) gi.GameInfo {
     var gameInfo gi.GameInfo
     var key string = strconv.FormatUint((uint64)(gameID), 10) + "old"

     db, _ := scribble.New("./info", nil)
     err := db.Read("info", key, &gameInfo)

     if err != nil {
     	logger.Error().Msgf("Error reading: %s", err)
     }
     
     return gameInfo
}

// UpdateGameInfo -- updates the game info present in the database
// this function writes a new record if there isn't already
// a record present
func UpdateGameInfo(gameInfo gi.GameInfo) {
     var storedGameInfo gi.GameInfo
     var oldGameInfo gi.GameInfo
     
     // Establish driver connection
     db, err := scribble.New("./info", nil)

     if err != nil {
     	logger.Error().Msgf("Driver failure: %s", err)
	return 
     }

     // retreive a single entry from info
     err = db.Read("info", strconv.FormatUint(uint64(gameInfo.ID), 10), &storedGameInfo)

     if err != nil {
     	logger.Warn().Msgf("game info update failed: %s", err)
     	writeGameInfo(db, gameInfo)
	return 
     }
     
     // Backup entry
     err = db.Read("info", strconv.FormatUint((uint64)(oldGameInfo.ID), 10) + "old", &oldGameInfo)

     if err != nil {
     	logger.Warn().Msgf("Reading game info backup failed: %s", err)
     }
     
     err = db.Delete("info", strconv.FormatUint((uint64)(oldGameInfo.ID), 10) + "old")

     if err != nil {
     	logger.Warn().Msgf("Deleting backup failed: %s", err)
     }

     err = db.Write("info", strconv.FormatUint((uint64)(storedGameInfo.ID), 10) + "old", storedGameInfo)
     if err != nil {
     	logger.Error().Msgf("Writing Backup failed: %s", err)
	return
     }
     
     err = db.Write("info", strconv.FormatUint((uint64)(gameInfo.ID), 10) , gameInfo)
     if err != nil {
     	logger.Error().Msgf("Writing game info entry failed: %s", err)
     }
}
