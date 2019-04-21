package main

import (
	"fmt"

	irp "github.com/HackRVA/master-base-2019/irpacket"
	serial "github.com/HackRVA/master-base-2019/serial"
)

func main() {
	serial.SetDebug(true)
	packetsIn := make(chan *irp.Packet)
	serial.OpenPort("/dev/ttyACM0", 9600)
	serial.InitIR()
	go serial.ReadSerial(packetsIn)
	for {
		packet := <-packetsIn
		fmt.Println("\nPacket received from packetsIn channel")
		packet.Print()
		packet.PrintPayload()
		fmt.Println()
	}
}
