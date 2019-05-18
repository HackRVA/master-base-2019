package info

import (
	"encoding/json"
	"strconv"

	log "github.com/HackRVA/master-base-2019/filelogging"
	gm "github.com/HackRVA/master-base-2019/game"
	gi "github.com/HackRVA/master-base-2019/gameinfo"
	scribble "github.com/nanobox-io/golang-scribble"
)

var logger = log.Ger.With().Str("pkg", "info").Logger()

// WriteGameInfo - Overwrites the 'info' entry with the gameInfo
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
	var key = strconv.FormatUint((uint64)(gameID), 10) + "old"

	db, _ := scribble.New("./data", nil)
	err := db.Read("info", key, &gameInfo)

	if err != nil {
		logger.Error().Msgf("Error reading: %s", err)
	}

	return gameInfo
}

// AddNewGameEntryToGameInfo -
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
