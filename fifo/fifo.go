package fifo

import (
	"bufio"
	"fmt"
	"log"
	"os"

	irp "github.com/HackRVA/master-base-2019/irpacket"
)

// FifoInFile - The path of the named pipe from the badge
var FifoInFile = "/tmp/fifo-from_badge"

// FifoOutFile - The path of the named pipe to the badge
var FifoOutFile = "/tmp/fifo-to-badge"

// ReadFifo - Reads a badge packet off of the named pipe (fifo)
func ReadFifo(fifoInFile string, packetsIn chan *irp.Packet) {
	fmt.Printf("Opening named pipe %s\n", fifoInFile)
	fifoFd, err := os.OpenFile(fifoInFile, os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		log.Fatal("Open Named pipe file error:", err)
	}

	buf := make([]byte, 4)

	reader := bufio.NewReader(fifoFd)

	for {
		byteCount, err := reader.Read(buf)
		if err != nil {
			log.Println("Error reading packet:", err)
		}
		if byteCount != 4 {
			log.Println("Packet read is not 4 bytes")
		}
		fmt.Println("bytes in:", buf)
		packetsIn <- irp.ReadPacket(irp.BytesToRawPacket(buf))
	}
}

// WriteFifo - Writes a badge packet to the named pipe (fifo)
func WriteFifo(fifoOutFile string, packetsOut chan *irp.Packet) {
	fmt.Printf("Opening named pipe %s\n", fifoOutFile)
	fifoFd, err := os.OpenFile(fifoOutFile, os.O_WRONLY, os.ModeNamedPipe)
	if err != nil {
		log.Fatal("Open Named pip error:", err)
	}

	writer := bufio.NewWriter(fifoFd)

	for {
		packet := <-packetsOut
		byteCount, err := writer.Write(irp.RawPacketToBytes(irp.WritePacket(packet)))
		if err != nil {
			log.Println("Error writing packet:", err)
		}
		if byteCount != 4 {
			log.Println("Packet written was not 4 bytes")
		}
	}
}
