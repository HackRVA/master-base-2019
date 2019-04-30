package game

import (
       zl "github.com/rs/zerolog"
       gm "github.com/HackRVA/master-base-2019/game"
       scribble "github.com/nanobox-io/golang-scribble"
       "time"
)

type BadgeInfo struct {
     timeSeen Time // Time a badge was last retreived
     badge    gm.Game
}

type GameInfo struct {
     details gm.Game
     badges map[BadgeInfo][]BadgeInfo
}

func NewGameInfo *GameInfo {
     var gameDetails gm.Game
     g.badges = make(map[uint16][]Game)
     return &g
}

func updateBadge(gameInfo GameInfo, game Game) {

     // Determine if a badge already exists
     // If it does update the badge info 
     badgeInfo, exists := gameInfo.badges[game.BadgeID]

     if (exists) {
     	currentTime := Time.now().UTC()

	badgeInfo.timeSeen = currentTime
	badgeInfo.badge = badge

	// This method always updates with
	// the original badgeID used
	// to retreive badge content
	gameInfo.badges[badgeInfo.BadgeID]
     }

     return gameInfo
}

// Overwrites the 'info' entry with the gameInfo 
func writeGameInfo(db scribble.Driver, gameInfo GameInfo) {

     allInfo, err = db.Write("info", "info", gameInfo)

     if err != nil {
     	logger.Error().Msgf("error writing to the database: %s", err)
     }
}
