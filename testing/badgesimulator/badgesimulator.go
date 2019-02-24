package main

import (
	"fmt"

	fifo "github.com/HackRVA/master-base-2019/fifo"
	irp "github.com/HackRVA/master-base-2019/irpacket"
	term "github.com/nsf/termbox-go"
)

const (
	Listening = iota
	Ignoring  = iota
)

func reset() {
	term.Sync() // cosmetic purpose?
}

func main() {
	// Set up input and output channels
	packetsIn := make(chan *irp.Packet)
	packetsOut := make(chan *irp.Packet)
	go fifo.ReadFifo(fifo.BadgeInFile, packetsIn)
	go fifo.WriteFifo(fifo.BadgeOutFile, packetsOut)

	err := term.Init()
	if err != nil {
		panic(err)
	}

	defer term.Close()

	status := Listening

	reset()
keyPressListenerLoop:
	for {
		switch status {
		case Listening:
			fmt.Println("Listening to base station")
		case Ignoring:
			fmt.Println("Ignoring base station")
		}
		fmt.Println("F5: Listen, F6: Ignore, Esc: Quit")

		switch ev := term.PollEvent(); ev.Type {
		case term.EventKey:
			switch ev.Key {
			case term.KeyEsc:
				break keyPressListenerLoop
			case term.KeyF5:
				status = Listening
			case term.KeyF6:
				status = Ignoring
			}

		case term.EventError:
			panic(ev.Err)
		}
		reset()
	}
}
