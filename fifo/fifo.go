package fifo

import (
	"bufio"
	"io"
	"os"

	log "github.com/HackRVA/master-base-2019/filelogging"
	irp "github.com/HackRVA/master-base-2019/irpacket"
	"github.com/hackebrot/go-repr/repr"
)

var debug = false
var logger = log.Ger

// SetDebug - sets the debugging on or off
func SetDebug(isDebug bool) {
	debug = isDebug
}

var connected = true

// SetConnected - If true, passes packets to the channels;
//                if false, the packets dispappear into the ether
//                in a simulation of IR communication
func SetConnected(isConnected bool) {
	connected = isConnected
}

// BadgeOutFile - The path of the named pipe from the badge
var BadgeOutFile = "/tmp/fifo-from-badge"

// BadgeInFile - The path of the named pipe to the badge
var BadgeInFile = "/tmp/fifo-to-badge"

// ReadFifo - Reads a badge packet off of the named pipe (fifo)
func ReadFifo(fifoInFile string, packetsIn chan *irp.Packet) {
	if debug {
		logger.Debug().Msgf("Opening named pipe %s\n", fifoInFile)
	}
	fifoFd, err := os.OpenFile(fifoInFile, os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		logger.Fatal().Msgf("Open Named pipe file error: %s", err)
	}

	buf := make([]byte, 4)

	reader := bufio.NewReader(fifoFd)

	for {
		buf[0], buf[1], buf[2], buf[3] = 0, 0, 0, 0
		byteCount, err := reader.Read(buf)
		if err != io.EOF {
			if err != nil {
				logger.Debug().Msgf("Error reading packet: %s", err)
			}
			if byteCount != 4 {
				logger.Debug().Msg("Packet read is not 4 bytes")
			}

			if debug {
				logger.Debug().Msgf("bytes in: %s", repr.Repr(buf))
			}

			packet := irp.PacketBytes(buf).Packet()

			if debug {
				packetLogger := packet.Logger(logger)
				packetLogger.Debug().Msgf("Packet read and routed to channel from: %s", fifoInFile)
			}

			if connected {
				packetsIn <- packet
			}
		}
	}
}

// WriteFifo - Writes a badge packet to the named pipe (fifo)
func WriteFifo(fifoOutFile string, packetsOut chan *irp.Packet) {
	if debug {
		logger.Debug().Msgf("Opening named pipe %s\n", fifoOutFile)
	}
	fifoFd, err := os.OpenFile(fifoOutFile, os.O_WRONLY, os.ModeNamedPipe)
	if err != nil {
		logger.Fatal().Msgf("Open Named pipe error: %s", err)
	}

	writer := bufio.NewWriter(fifoFd)

	for {
		packet := <-packetsOut

		if connected {
			if debug {
				packetLogger := packet.Logger(logger)
				packetLogger.Debug().Msgf("Packet to write received from channel: %s", fifoOutFile)
			}

			bytes := packet.Bytes()
			//irp.RawPacketToBytes(irp.WritePacket(packet))

			if debug {
				logger.Debug().Msgf("bytes out: %s", repr.Repr(bytes))
			}

			byteCount, err := writer.Write(bytes)
			if err != nil {
				logger.Error().Msgf("Error writing packet: %s", err)
			}
			if byteCount != 4 {
				logger.Error().Msg("Packet written was not 4 bytes")
			}
			err = writer.Flush()
			if err != nil {
				logger.Error().Msgf("Error flushing buffer: %s", err)
			}
		}
	}
}
