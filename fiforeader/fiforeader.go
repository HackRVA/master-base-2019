package main

import (
	fifo "github.com/HackRVA/master-base-2019/fifo"
	irp "github.com/HackRVA/master-base-2019/irpacket"
	msg "github.com/HackRVA/master-base-2019/messages"
)

func main() {
	packetsIn := make(chan *irp.Packet)
	go fifo.ReadFifo(packetsIn)
	go msg.ProcessPackets(packetsIn)
	for {
	}
}
