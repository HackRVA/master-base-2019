package fifo

import (
	"bufio"
	"fmt"
	"log"
	"os"

	irp "github.com/HackRVA/master-base-2019/irpacket"
)

var fifoInFile = "/tmp/badge-ir-fifo"

// ReadFifo - Reads a badge packet off of the named pipe (fifo)
func ReadFifo(c chan *irp.Packet) {
	fmt.Printf("Opening named pipe %s\n", fifoInFile)
	fifoFd, err := os.OpenFile(fifoInFile, os.O_RDONLY, os.ModeNamedPipe)
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
			log.Println("Packet not 4 bytes")
		}
		fmt.Println("bytes in:", buf)
		c <- irp.ReadPacket(irp.BytesToRawPacket(buf))
	}
}
