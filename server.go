package main

import (
	"net/http"

	api "github.com/HackRVA/master-base-2019/baseapi"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/newgame", api.Newgame).Methods("POST")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./client/")))
	http.Handle("/", r)
	http.ListenAndServe(":3000", nil)
}
