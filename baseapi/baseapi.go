package baseapi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	msg "github.com/HackRVA/master-base-2019/messages"
)

// GameSpec with absolute start time
type extendedGameSpec struct {
	msg.GameSpec
	AbsStart string
}

// Newgame - function to schedule newgame
func Newgame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var e extendedGameSpec
	var s msg.GameSpec

	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &e)

	form := "2006-01-02 15:04:05"
	t2, _ := time.Parse(form, e.AbsStart)

	s.StartTime = Until(t2)
	s.Duration = e.Duration
	s.Variant = e.Variant
	s.GameID = 1

	j, _ := json.Marshal(s)
	w.Write(j)
}
