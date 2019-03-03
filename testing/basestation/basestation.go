package main

import (
	"fmt"

	fifo "github.com/HackRVA/master-base-2019/fifo"
	log "github.com/HackRVA/master-base-2019/filelogging"
	irp "github.com/HackRVA/master-base-2019/irpacket"
	msg "github.com/HackRVA/master-base-2019/messages"
	term "github.com/nsf/termbox-go"
)

var logger = log.Ger

func reset() {
	term.Sync() // cosmetic purpose?
}

func main() {
	// Set up input and output channels
	packetsIn := make(chan *irp.Packet)
	packetsOut := make(chan *irp.Packet)
	gameData := make(chan *msg.GameData)
	beaconHold := make(chan bool)
	gameSpec := make(chan *msg.GameSpec)

	go fifo.ReadFifo(fifo.BadgeOutFile, packetsIn)
	go fifo.WriteFifo(fifo.BadgeInFile, packetsOut)
	fifo.SetDebug(true)
	msg.SetDebug(true)

	err := term.Init()
	if err != nil {
		panic(err)
	}

	defer term.Close()

	go msg.ReceivePackets(packetsIn, gameData, beaconHold)

	go msg.TransmitBeacon(packetsOut, beaconHold)

	go msg.TransmitNewGamePackets(packetsOut, gameSpec, beaconHold)

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
