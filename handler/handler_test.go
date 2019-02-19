package handler_test

import (
	"fmt"
	"io/ioutil"

	handler "github.com/HackRVA/master-base-2019/handler"
)

func main() {
	games, _ := ioutil.ReadFile("games.json")
	fmt.Println(handler.Readgames(games))
}
