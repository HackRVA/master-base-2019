package main

import (
	"fmt"

	msg "github.com/HackRVA/master-base-2019/badgewrangler"
	"github.com/HackRVA/master-base-2019/fifo"
	irp "github.com/HackRVA/master-base-2019/irpacket"
	serial "github.com/HackRVA/master-base-2019/serial"
	term "github.com/nsf/termbox-go"
)

func reset() {
	term.Sync() // cosmetic purpose?
}

func main() {
	serial.SetDebug(true)
	packetsIn := make(chan *irp.Packet)
	packetsOut := make(chan *irp.Packet)
	serial.OpenPort("/dev/ttyACM0", 19200)
	serial.InitIR()
	go serial.ReadSerial(packetsIn)
	go serial.WriteSerial(packetsOut)
	beaconPacket := msg.BuildBeacon()

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
				} else if ev.Ch == 'b' {
					fmt.Println("\nBeacon packet built:")
					beaconPacket.Print()
					fmt.Println()

					packetsOut <- beaconPacket
				}
			case term.EventError:
				panic(ev.Err)
			}
		default:
		}
	}
}
