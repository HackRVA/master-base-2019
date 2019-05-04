package gameinfo

import (
       bw "github.com/HackRVA/master-base-2019/badgewrangler"
       gm "github.com/HackRVA/master-base-2019/game"
       "time"
)

type BadgeInfo struct {
     TimeSeen time.Time // Time a badge was last retreived
     ID       uint16
     Data     bw.GameData
}

type GameInfo struct {
     ID      uint16
     Details gm.Game
     Badges  map[uint16]BadgeInfo
}

func NewGameInfo(game gm.Game) *GameInfo {
     
     var g GameInfo = GameInfo{}
     g.Badges = make(map[uint16]BadgeInfo)

     // If game is not empty add game details
     if (gm.Game{}) != game {
     	g.Details = game
     	g.ID = game.GameID
     }

     return &g
}

func UpdateBadgeData(gameInfo GameInfo, gameData bw.GameData) GameInfo {

     // If gameData is empty return
     if len(gameData.Hits) == 0 && gameData.GameID == 0 {
     	return gameInfo
     }
     
     currentTime := time.Now().UTC()
     
     // Determine if a badge already exists
     // If it does update the badge info 
     badgeInfo, exists := gameInfo.Badges[gameData.BadgeID]

     if (exists) {

	badgeInfo.TimeSeen = currentTime
	badgeInfo.Data = gameData

	// This method always updates with
	// the original badgeID used
	// to retreive badge content
	gameInfo.Badges[badgeInfo.ID] = badgeInfo
     } else {

       var badge BadgeInfo = BadgeInfo{}
       
       badge.TimeSeen = currentTime
       badge.ID       = gameData.BadgeID
       badge.Data     = gameData

       gameInfo.Badges[gameData.BadgeID] = badge
       
     }

     return gameInfo
}

