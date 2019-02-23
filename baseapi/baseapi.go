package baseapi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	msg "github.com/HackRVA/master-base-2019/messages"
)

// Newgame - function to schedule newgame
func Newgame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var m msg.GameSpec
	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &m)

	j, _ := json.Marshal(m)
	w.Write(j)

}
