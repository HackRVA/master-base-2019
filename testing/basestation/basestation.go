package main

import (
	"fmt"

	log "github.com/HackRVA/master-base-2019/filelogging"
	ss "github.com/HackRVA/master-base-2019/serverstartup"
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

	ss.StartBadgeWrangler("/dev/ttyACM0", 9600)

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
