package main

import (
	"fmt"

	bw "github.com/HackRVA/master-base-2019/badgewrangler"
	log "github.com/HackRVA/master-base-2019/filelogging"
	term "github.com/nsf/termbox-go"
)

var logger = log.Ger

func reset() {
	term.Sync() // cosmetic purpose?
}

func main() {

	err := term.Init()
	if err != nil {
		panic(err)
	}

	defer term.Close()

	bw.StartWrangler()

	reset()

keyPressListenerLoop:
	for {
		switch ev := term.PollEvent(); ev.Type {
		case term.EventKey:
			if ev.Key == term.KeyEsc {
				fmt.Println("Esc pressed")
				break keyPressListenerLoop
			}
		case term.EventError:
			panic(ev.Err)
		}
	}
}
