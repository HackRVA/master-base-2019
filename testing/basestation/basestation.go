package main

import (
	"fmt"

	bw "github.com/HackRVA/master-base-2019/badgewrangler"
	fifo "github.com/HackRVA/master-base-2019/fifo"
	log "github.com/HackRVA/master-base-2019/filelogging"
	irp "github.com/HackRVA/master-base-2019/irpacket"
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
	gameData := make(chan *bw.GameData)
	beaconHold := make(chan bool)
	gameSpec := make(chan *bw.GameSpec)

	go fifo.ReadFifo(fifo.BadgeOutFile, packetsIn)
	go fifo.WriteFifo(fifo.BadgeInFile, packetsOut)
	fifo.SetDebug(true)
	bw.SetDebug(true)

	err := term.Init()
	if err != nil {
		panic(err)
	}

	defer term.Close()

	go bw.ReceivePackets(packetsIn, gameData, beaconHold)

	go bw.TransmitBeacon(packetsOut, beaconHold)

	go bw.TransmitNewGamePackets(packetsOut, gameSpec, beaconHold)

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
