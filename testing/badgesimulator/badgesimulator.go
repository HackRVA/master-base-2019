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

const (
	listening = iota
	ignoring  = iota
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
	fifo.SetConnected(false)
	fifo.SetDebug(true)
	bw.SetDebug(true)

	err := term.Init()
	if err != nil {
		panic(err)
	}

	defer term.Close()

	status := ignoring

	gameData := &bw.GameData{
		BadgeID: uint16(333),
		GameID:  uint16(1234),
		Hits: []*bw.Hit{
			{BadgeID: uint16(101), Timestamp: uint16(33), Team: uint8(2)},
			{BadgeID: uint16(101), Timestamp: uint16(103), Team: uint8(2)},
			{BadgeID: uint16(101), Timestamp: uint16(203), Team: uint8(2)},
			{BadgeID: uint16(101), Timestamp: uint16(303), Team: uint8(2)},
			{BadgeID: uint16(101), Timestamp: uint16(403), Team: uint8(2)},
			{BadgeID: uint16(101), Timestamp: uint16(503), Team: uint8(2)},
			{BadgeID: uint16(101), Timestamp: uint16(603), Team: uint8(2)},
			{BadgeID: uint16(101), Timestamp: uint16(703), Team: uint8(2)}}}

	go bw.BadgeHandlePackets(packetsIn, packetsOut, gameData)

	reset()
keyPressListenerLoop:
	for {
		switch status {
		case listening:
			fmt.Println("Listening to base station")
			fifo.SetConnected(true)
		case ignoring:
			fmt.Println("Ignoring base station")
			fifo.SetConnected(false)
		}
		fmt.Println("F5: Listen, F6: Ignore, Esc: Quit")

		switch ev := term.PollEvent(); ev.Type {
		case term.EventKey:
			switch ev.Key {
			case term.KeyEsc:
				break keyPressListenerLoop
			case term.KeyF5:
				status = listening
			case term.KeyF6:
				status = ignoring
			}

		case term.EventError:
			panic(ev.Err)
		}
		reset()
	}
}
