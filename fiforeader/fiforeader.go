package main

import (
	"fmt"

	"github.com/HackRVA/master-base-2019/fifo"
	irp "github.com/HackRVA/master-base-2019/irpacket"
)

func main() {
	packetsIn := make(chan *irp.Packet)
	go fifo.ReadFifo(packetsIn)
	for {
		pkt := <-packetsIn
		fmt.Println()
		irp.PrintPacket(pkt)
	}
}
