package main

import (
	"fmt"
	"os"
	"os/signal"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c)
	s := <-c
	fmt.Println("Got signal:", s)
	if s.String() == "^C" {
		fmt.Println("Okey Dokey")
	}
}
