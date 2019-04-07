package main

import (
	"fmt"

	msg "github.com/HackRVA/master-base-2019/badgewrangler"
	irp "github.com/HackRVA/master-base-2019/irpacket"
	"github.com/HackRVA/master-base-2019/serial"
	term "github.com/nsf/termbox-go"
)

func reset() {
	term.Sync() // cosmetic purpose?
}

func main() {
	serial.SetDebug(true)
	packetsOut := make(chan *irp.Packet)
	serial.OpenPort("/dev/ttyUSB0", 9600)
	go serial.WriteSerial(packetsOut)
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
