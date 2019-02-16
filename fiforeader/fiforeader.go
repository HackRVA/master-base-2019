package main

import (
	"fmt"

	fifo "github.com/HackRVA/master-base-2019/fifo"
	irp "github.com/HackRVA/master-base-2019/irpacket"
)

func main() {
	packetsIn := make(chan *irp.Packet)
	go fifo.ReadFifo(fifo.FifoOutFile, packetsIn)
	//go msg.ProcessPackets(packetsIn)
	for {
		packet := <-packetsIn
		fmt.Println()
		irp.PrintPacket(packet)
		fmt.Println()
	}
}
