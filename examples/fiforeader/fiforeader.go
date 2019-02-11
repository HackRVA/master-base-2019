package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	irp "github.com/wjodon/master-base/irpacket"
)

var fifofile = "/tmp/badge-ir-fifo"

func main() {
	fmt.Printf("Opening named pipe %s\n", fifofile)
	fifoFd, err := os.OpenFile(fifofile, os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		log.Fatal("Open Named pipe file error:", err)
	}

	buf := make([]byte, 4)

	reader := bufio.NewReader(fifoFd)

	for {
		bytes, err := reader.Read(buf)
		if err != nil {
			log.Println("Error reading packet")
		}
		if bytes != 4 {
			log.Println("Packet is wrong length")
		}
		fmt.Println(buf)
		irp.PrintPacket(irp.ReadPacket(irp.BytesToRawPacket(buf)))
	}

}
