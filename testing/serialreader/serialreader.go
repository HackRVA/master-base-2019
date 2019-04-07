package main

import (
	"fmt"

	irp "github.com/HackRVA/master-base-2019/irpacket"
	serial "github.com/HackRVA/master-base-2019/serial"
)

func main() {
	serial.SetDebug(true)
	packetsIn := make(chan *irp.Packet)
	serial.OpenPort("/dev/ttyUSB1", 9600)
	go serial.ReadSerial(packetsIn)
	for {
		packet := <-packetsIn
		fmt.Println("\nPacket received from packetsIn channel")
		packet.Print()
		packet.PrintPayload()
		fmt.Println()
	}
}
