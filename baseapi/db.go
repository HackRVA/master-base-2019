package baseapi

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/cnf/structhash"
	scribble "github.com/nanobox-io/golang-scribble"
)

// SaveGame -- save gamespec to database
func SaveGame(game Game) {
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
func GetNext() Game {
	t := time.Now()
	var g Game
	// TODO: send next game or *Current Game
	games := GetGames()
	for _, game := range games {
		fmt.Println(
			"now: ",
			int64(t.Unix()),
			"\nnextGame: ",
			game.AbsStart,
			"\nin the future:",
			int64(t.Unix()) < game.AbsStart)

		// return the first game that is greater than now
		if int64(t.Unix()) < game.AbsStart+int64(game.Duration) {
			return game
		}
	}

	return g
}

// GetGames -- retrieves games from DB
func GetGames() []Game {
	// create a new scribble database, providing a destination for the database to live
	db, _ := scribble.New("./data", nil)
	// Read more games from the database
	moregames, _ := db.ReadAll("games")
	// iterate over moregames creating a new game for each record
	games := []Game{}
	for _, game := range moregames {
		g := Game{}
		json.Unmarshal([]byte(game), &g)
		games = append(games, g)
	}

	return games
}

// DataInGameOut - stores the game data and gets the current/next game
func DataInGameOut(gameDataIn chan *GameData, gameOut chan *Game) {
	for {
		gameData := <-gameDataIn

		gameOut <- &ba.GetNext()
	}
}
