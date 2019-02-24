package main

import (
	"fmt"

	fifo "github.com/HackRVA/master-base-2019/fifo"
	irp "github.com/HackRVA/master-base-2019/irpacket"
	msg "github.com/HackRVA/master-base-2019/messages"
)

func main() {
	packetsOut := make(chan *irp.Packet)
	go fifo.WriteFifo(fifo.BadgeOutFile, packetsOut)
	packet := msg.BuildBeacon()

	fmt.Println("\nPacket built:")
	packet.Print()
	fmt.Println()

	packetsOut <- packet
	for {
	}
}
