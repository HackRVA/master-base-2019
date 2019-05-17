package baseapi

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	bw "github.com/HackRVA/master-base-2019/badgewrangler"
	gm "github.com/HackRVA/master-base-2019/game"
	gi "github.com/HackRVA/master-base-2019/gameinfo"
	"github.com/cnf/structhash"
	scribble "github.com/nanobox-io/golang-scribble"
)

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
			"now: %d \nnextGame: %d \nin the future: %t",
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
		fmt.Println(gameData.GameID)
		//gameDataOut <- gameData
		nextGame := GetNext()
		gameOut <- &nextGame
		SaveGameData(gameData)
	}
}

// Overwrites the 'info' entry with the gameInfo
func WriteGameInfo(gameInfo gi.GameInfo) {
	db, _ := scribble.New("./data", nil)
	if err := db.Write("info", strconv.FormatInt(int64(gameInfo.ID), 10), gameInfo); err != nil {
		logger.Error().Msgf("error writing to the database: %s", err)
	}
}

// GetAllInfo -- retrieves all info entries from the DB
func GetAllInfo() []gi.GameInfo {
	// create a new scribble database, providing a destination for the database to live
	db, _ := scribble.New("./data", nil)
	// Read the info table from the database
	resultSet, _ := db.ReadAll("info")
	// iterate over the info result-set
	allInfo := []gi.GameInfo{}
	for _, result := range resultSet {
		info := gi.GameInfo{}
		json.Unmarshal([]byte(result), &info)
		allInfo = append(allInfo, info)
	}

	return allInfo
}

// GetInfo -- retreive game info from database
func GetInfo(gameID uint16) gi.GameInfo {
	var gameInfo gi.GameInfo

	db, _ := scribble.New("./data", nil)
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

	db, _ := scribble.New("./data", nil)
	err := db.Read("info", key, &gameInfo)

	if err != nil {
		logger.Error().Msgf("Error reading: %s", err)
	}

	return gameInfo
}

func AddNewGameEntryToGameInfo(game gm.Game) {
	var gameInfo gi.GameInfo

	// Establish driver connection
	db, err := scribble.New("./data", nil)
	logger.Info().Msg("Adding new entry")

	if err != nil {
		logger.Error().Msgf("Driver failure: %s", err)
		return
	}

	err = db.Read("info", strconv.FormatUint((uint64)(game.GameID), 10), gameInfo)

	if err != nil {
		logger.Error().Msgf("Read error: %s", err)
	}

	// If no prior gameInfo information exists for this game
	// Initialize gameInfo with game
	if gameInfo.ID == 0 {
		gameInfo.ID = game.GameID
		gameInfo.Details = game
	} else {
		gameInfo.Details = game
	}

	err = db.Write("info", strconv.FormatInt((int64)(gameInfo.ID), 10), &gameInfo)
	if err != nil {
		logger.Error().Msgf("error writing to the database: %s", err)
	}

}

// UpdateGameInfo -- updates the game info present in the database
// this function writes a new record if there isn't already
// a record present
func UpdateGameInfo(gameInfo gi.GameInfo) {
	var storedGameInfo gi.GameInfo
	var oldGameInfo gi.GameInfo

	// Establish driver connection
	db, err := scribble.New("./data", nil)

	if err != nil {
		logger.Error().Msgf("Driver failure: %s", err)
		return
	}

	// retrieve a single entry from info
	err = db.Read("info", strconv.FormatUint(uint64(gameInfo.ID), 10), &storedGameInfo)

	if err != nil {
		logger.Warn().Msgf("game info update failed: %s", err)
		WriteGameInfo(gameInfo)
		return
	}

	// Backup entry
	err = db.Read("info", strconv.FormatUint((uint64)(oldGameInfo.ID), 10)+"old", &oldGameInfo)

	if err != nil {
		logger.Warn().Msgf("Reading game info backup failed: %s", err)
	}

	err = db.Delete("info", strconv.FormatUint((uint64)(oldGameInfo.ID), 10)+"old")

	if err != nil {
		logger.Warn().Msgf("Deleting backup failed: %s", err)
	}

	err = db.Write("info", strconv.FormatUint((uint64)(storedGameInfo.ID), 10)+"old", storedGameInfo)
	if err != nil {
		logger.Error().Msgf("Writing Backup failed: %s", err)
		return
	}

	err = db.Write("info", strconv.FormatUint((uint64)(gameInfo.ID), 10), gameInfo)
	if err != nil {
		logger.Error().Msgf("Writing game info entry failed: %s", err)
	}
}
