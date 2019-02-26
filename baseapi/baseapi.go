package baseapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	msg "github.com/HackRVA/master-base-2019/messages"
)

// NewGame - function to schedule newgame
func NewGame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var e Game
	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &e)

	form := "2006-01-02 15:04:05"
	t2, _ := time.Parse(form, e.AbsStart)

	newSpec := msg.GameSpec{
		StartTime: Until(t2),
		Duration:  e.Duration,
		Variant:   e.Variant,
		GameID:    1,
	}
	fmt.Println(newSpec)
	SaveGame(e)

	j, _ := json.Marshal(newSpec)
	w.Write(j)
}

// NextGame -- returns the game that is sheduled next
func NextGame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	GetNext()
}

// AllGames - returns all scheduled games
func AllGames(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	j, _ := json.Marshal(GetGames())
	w.Write(j)
}
