package fifo

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	irp "github.com/HackRVA/master-base-2019/irpacket"
)

// BadgeOutFile - The path of the named pipe from the badge
var BadgeOutFile = "/tmp/fifo-from-badge"

// BadgeInFile - The path of the named pipe to the badge
var BadgeInFile = "/tmp/fifo-to-badge"

var debug = false

// SetDebug - sets the debugging on or off
func SetDebug(isDebug bool) {
	debug = isDebug
}

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
		buf[0], buf[1], buf[2], buf[3] = 0, 0, 0, 0
		byteCount, err := reader.Read(buf)
		if err != io.EOF {
			if err != nil {
				log.Println("Error reading packet:", err)
			}
			if byteCount != 4 {
				log.Println("Packet read is not 4 bytes")
			}

			if debug {
				fmt.Println("bytes in:", buf)
			}

			//packet := irp.ReadPacket(irp.BytesToRawPacket(buf))
			packet := irp.PacketBytes(buf).Packet()

			if debug {
				fmt.Println("\nPacket read and going to channel:")
				packet.Print()
				fmt.Println()
			}

			packetsIn <- packet
		}
	}
}

// WriteFifo - Writes a badge packet to the named pipe (fifo)
func WriteFifo(fifoOutFile string, packetsOut chan *irp.Packet) {
	fmt.Printf("Opening named pipe %s\n", fifoOutFile)
	fifoFd, err := os.OpenFile(fifoOutFile, os.O_WRONLY, os.ModeNamedPipe)
	if err != nil {
		log.Fatal("Open Named pipe error:", err)
	}

	writer := bufio.NewWriter(fifoFd)

	for {
		packet := <-packetsOut

		if debug {
			fmt.Println("\nPacket to write recieved from channel:")
			packet.Print()
			fmt.Println()
		}

		bytes := packet.Bytes()
		//irp.RawPacketToBytes(irp.WritePacket(packet))

		if debug {
			fmt.Println("bytes out:", bytes)
			fmt.Println()
		}

		byteCount, err := writer.Write(bytes)
		if err != nil {
			log.Println("Error writing packet:", err)
		}
		if byteCount != 4 {
			log.Println("Packet written was not 4 bytes")
		}
		writer.Flush()
		if err != nil {
			log.Println("Error flushing buffer:", err)
		}
	}
}
