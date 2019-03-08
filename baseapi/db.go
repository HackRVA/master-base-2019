package baseapi

import (
	"encoding/json"
	"fmt"
	"time"

	bw "github.com/HackRVA/master-base-2019/badgewrangler"
	"github.com/cnf/structhash"
	scribble "github.com/nanobox-io/golang-scribble"
)

// SaveGame -- save gamespec to database
func SaveGame(game bw.Game) {
	hash, err := structhash.Hash(game, 1)
	if err != nil {
		panic(err)
	}

	// create a new scribble database, providing a destination for the database to live
	db, _ := scribble.New("./data", nil)
	fmt.Println(game)
	db.Write("games", hash, game)
}

// GetNext -- return the next game
func GetNext() bw.Game {
	t := time.Now()
	var g bw.Game
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
			return game
		}
	}

	return g
}

// GetGames -- retrieves games from DB
func GetGames() []bw.Game {
	// create a new scribble database, providing a destination for the database to live
	db, _ := scribble.New("./data", nil)
	// Read more games from the database
	moregames, _ := db.ReadAll("games")
	// iterate over moregames creating a new game for each record
	games := []bw.Game{}
	for _, game := range moregames {
		g := bw.Game{}
		json.Unmarshal([]byte(game), &g)
		games = append(games, g)
	}

	return games
}

// DataInGameOut - stores the game data and gets the current/next game
func DataInGameOut(gameDataIn chan *bw.GameData, gameOut chan *bw.Game) {
	for {
		gameData := <-gameDataIn
		fmt.Println(gameData.GameID)
		nextGame := GetNext()
		gameOut <- &nextGame
	}
}
