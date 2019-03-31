package main

import (
	"fmt"

	msg "github.com/HackRVA/master-base-2019/badgewrangler"
	fifo "github.com/HackRVA/master-base-2019/fifo"
	irp "github.com/HackRVA/master-base-2019/irpacket"
	term "github.com/nsf/termbox-go"
)

func reset() {
	term.Sync() // cosmetic purpose?
}

func main() {
	packetsOut := make(chan *irp.Packet)
	go fifo.WriteFifo(fifo.BadgeInFile, packetsOut)
	packet := msg.BuildBeacon()

	err := term.Init()
	if err != nil {
		panic(err)
	}

	defer term.Close()

	reset()

keyPressListenerLoop:
	for {
		switch ev := term.PollEvent(); ev.Type {
		case term.EventKey:
			if ev.Key == term.KeyEsc {
				fmt.Println("Esc pressed")
				break keyPressListenerLoop
			} else if ev.Ch == 'p' {
				fmt.Println("\nPacket built:")
				packet.Print()
				fmt.Println()

				packetsOut <- packet
			}
		case term.EventError:
			panic(ev.Err)
		}
	}
}
