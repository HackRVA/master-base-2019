package main

import (
	"fmt"

	fifo "github.com/HackRVA/master-base-2019/fifo"
	irp "github.com/HackRVA/master-base-2019/irpacket"
)

func main() {
	packetsIn := make(chan *irp.Packet)
	go fifo.ReadFifo(fifo.BadgeOutFile, packetsIn)
	for {
		packet := <-packetsIn
		fmt.Println("\nPacket received from packetsIn channel")
		packet.Print()
		packet.PrintPayload()
		fmt.Println()
	}
}
