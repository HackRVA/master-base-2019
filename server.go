package main

import (
	"fmt"
	"net/http"

	log "github.com/HackRVA/master-base-2019/filelogging"
	lb "github.com/HackRVA/master-base-2019/leaderboard"
	ss "github.com/HackRVA/master-base-2019/serverstartup"
	"github.com/HackRVA/master-base-2019/sync"
	api "github.com/HackRVA/master-base-2019/webapi"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

var logger = log.Ger

func main() {
	ss.InitConfiguration()

	r := mux.NewRouter()
	r.HandleFunc("/api/newgame", api.NewGame).Methods("POST")
	r.HandleFunc("/api/nextgame", api.NextGame).Methods("GET")
	r.HandleFunc("/api/games", api.AllGames).Methods("GET")
	r.HandleFunc("/api/info/all", api.AllInfo).Methods("GET")
	r.HandleFunc("/api/info/{id}", api.Info).Methods("GET")
	r.HandleFunc("/api/zombie", api.Zombie).Methods("POST")
	http.Handle("/", r)
	fmt.Println("running web server on port 8000")

	isMaster := viper.GetBool("isMaster")
	if isMaster == false {
		sync.StartSyncLoop()
	}

	lb.StartLeaderboardLoop()
	ss.StartBadgeWrangler()
	http.ListenAndServe(":8000", nil)
}
