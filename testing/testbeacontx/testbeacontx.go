package main

import (
	bw "github.com/HackRVA/master-base-2019/badgewrangler"
	fifo "github.com/HackRVA/master-base-2019/fifo"
	irp "github.com/HackRVA/master-base-2019/irpacket"
)

func main() {
	packetsOut := make(chan *irp.Packet)
	beaconHold := make(chan bool)
	go fifo.WriteFifo(fifo.BadgeInFile, packetsOut)
	go bw.TransmitBeacon(packetsOut, beaconHold)
	for {
	}
}
