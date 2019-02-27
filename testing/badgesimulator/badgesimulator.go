package main

import (
	"fmt"
	"os"

	fifo "github.com/HackRVA/master-base-2019/fifo"
	irp "github.com/HackRVA/master-base-2019/irpacket"
	msg "github.com/HackRVA/master-base-2019/messages"
	term "github.com/nsf/termbox-go"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(os.Stdout)

	file, err := os.OpenFile("test.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Info("Failed to log to file, using default stderr")
	}
	log.SetLevel(log.DebugLevel)

}

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
	msg.SetDebug(true)

	err := term.Init()
	if err != nil {
		panic(err)
	}

	defer term.Close()

	status := ignoring

	gameData := &msg.GameData{
		BadgeID: uint16(2322),
		GameID:  uint16(1234),
		Hits: []*msg.Hit{
			{BadgeID: uint16(101), Timestamp: uint16(33), Team: uint8(2)},
			{BadgeID: uint16(101), Timestamp: uint16(103), Team: uint8(2)},
			{BadgeID: uint16(101), Timestamp: uint16(203), Team: uint8(2)},
			{BadgeID: uint16(101), Timestamp: uint16(303), Team: uint8(2)},
			{BadgeID: uint16(101), Timestamp: uint16(403), Team: uint8(2)},
			{BadgeID: uint16(101), Timestamp: uint16(503), Team: uint8(2)},
			{BadgeID: uint16(101), Timestamp: uint16(603), Team: uint8(2)},
			{BadgeID: uint16(101), Timestamp: uint16(703), Team: uint8(2)}}}

	go msg.BadgeHandlePackets(packetsIn, packetsOut, gameData)

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
