package baseapi

import (
	"encoding/json"
	"fmt"

	msg "github.com/HackRVA/master-base-2019/messages"

	"github.com/cnf/structhash"
	"github.com/nanobox-io/golang-scribble"
)

// Game -- with absolute start time and ID
type Game struct {
	msg.GameSpec
	AbsStart string
}

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
func GetNext() {
	// TODO: sort by date to find the next date
	games := GetGames()
	for _, game := range games {
		fmt.Println(game.AbsStart)
	}
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
