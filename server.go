package main

import (
	"fmt"
	"net/http"
	"time"

	api "github.com/HackRVA/master-base-2019/baseapi"
	"github.com/gorilla/mux"
)

func sendToLeaderboard() {
	ticker := time.NewTicker(180 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Println("sending data to leaderboard")
				api.PostGameData(api.GetGameData())
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/newgame", api.NewGame).Methods("POST")
	r.HandleFunc("/api/nextgame", api.NextGame).Methods("GET")
	r.HandleFunc("/api/games", api.AllGames).Methods("GET")
	http.Handle("/", r)
	fmt.Println("running web server on port 8000")
	sendToLeaderboard()
	http.ListenAndServe(":8000", nil)
}
