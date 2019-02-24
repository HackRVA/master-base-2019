package main

import (
	fifo "github.com/HackRVA/master-base-2019/fifo"
	irp "github.com/HackRVA/master-base-2019/irpacket"
	msg "github.com/HackRVA/master-base-2019/messages"
)

func main() {
	packetsOut := make(chan *irp.Packet)
	beaconHold := make(chan bool)
	go fifo.WriteFifo(fifo.FifoOutFile, packetsOut)
	go msg.TransmitBeacon(packetsOut, beaconHold)
	for {
	}
}
