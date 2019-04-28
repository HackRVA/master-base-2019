package main

import (
	"fmt"

	fifo "github.com/HackRVA/master-base-2019/fifo"
	irp "github.com/HackRVA/master-base-2019/irpacket"
	term "github.com/nsf/termbox-go"
)

func reset() {
	term.Sync() // cosmetic purpose?
}

func main() {
	packetsIn := make(chan *irp.Packet)
	go fifo.ReadFifo(fifo.BadgeOutFile, packetsIn)
	//go msg.ProcessPackets(packetsIn)

	err := term.Init()
	if err != nil {
		panic(err)
	}

	defer term.Close()

	reset()

	termEvent := make(chan term.Event)

	go func(termEvent chan term.Event) {
		for {
			ev := term.PollEvent()
			termEvent <- ev
		}
	}(termEvent)

keyPressListenerLoop:
	for {
		select {
		case packet := <-packetsIn:
			fmt.Println("\nPacket received from", fifo.BadgeOutFile, "channel")
			fmt.Println("Esc to quit")
			packet.Print()
			packet.PrintPayload()
			fmt.Println()
		case ev := <-termEvent:
			switch ev.Type {
			case term.EventKey:
				if ev.Key == term.KeyEsc {
					fmt.Println("Esc pressed")
					break keyPressListenerLoop
				}
			case term.EventError:
				panic(ev.Err)
			}
		default:
		}
	}
}
