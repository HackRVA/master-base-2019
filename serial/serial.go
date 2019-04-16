package serial

import (
	log "github.com/HackRVA/master-base-2019/filelogging"
	irp "github.com/HackRVA/master-base-2019/irpacket"
	"github.com/hackebrot/go-repr/repr"
	"github.com/tarm/serial"
)

var logger = log.Ger
var debug = false
var connected = true
var serialConn *serial.Port

// SetDebug - set debug on/off
func SetDebug(isDebug bool) {
	debug = isDebug
}

// SetConnected - If true, passes packets to the channels;
//                if false, the packets dispappear into the ether
//                in a simulation of IR communication
func SetConnected(isConnected bool) {
	connected = isConnected
}

// OpenPort - Open a serial port
func OpenPort(portName string, baud int) {
	var err error

	config := &serial.Config{Name: portName, Baud: baud}
	serialConn, err = serial.OpenPort(config)
	if err != nil {
		logger.Fatal().Err(err).Msg("Error opening port")
	}
}

//ReadSerial - Reads a badge packet from the serial port
func ReadSerial(packetsIn chan *irp.Packet) {

	buf := make([]byte, 1)
	packetBuffer := make([]byte, 0, 4)

	for {

		for len(packetBuffer) < 4 {
			buf[0] = 0
			byteCount, err := serialConn.Read(buf)
			if err != nil {
				logger.Debug().Msgf("Error reading packet: %s", err)
			}
			if byteCount != 1 {
				logger.Debug().Msgf("Packet read is not 4 bytes, it is %d bytes", byteCount)
			}

			packetBuffer = append(packetBuffer, buf[0])
		}

		if debug {
			logger.Debug().Msgf("bytes in: %s", repr.Repr(packetBuffer))
		}

		packet := irp.PacketBytes(packetBuffer).Packet()

		if debug {
			packetLogger := packet.Logger(logger)
			packetLogger.Debug().Msgf("Packet read from serial and routed to channel")
		}

		if connected {
			packetsIn <- packet
		}
		packetBuffer = packetBuffer[:0]
	}

}

// InitIR - writes an IR initialization sequence to the serial port
func InitIR() {
	if debug {
		logger.Debug().Msg("Initializing IR")
	}
	byteCount, err := serialConn.Write([]byte("ZsYnC#"))
	if err != nil {
		logger.Fatal().Err(err).Msg("Error initializing IR")
	}
	if byteCount != 6 {
		logger.Fatal().Msg("IR init did not write 6 bytes")
	}

	err = serialConn.Flush()
	if err != nil {
		logger.Fatal().Err(err).Msg("Error flushing the buffer")
	}
}

// WriteSerial - writes a packet to the serial port
func WriteSerial(packetsOut chan *irp.Packet) {

	for {
		packet := <-packetsOut

		if connected {

			if debug {
				packetLogger := packet.Logger(logger)
				packetLogger.Debug().Msgf("Packet to write received from channel")
			}

			bytes := packet.Bytes()

			if debug {
				logger.Debug().Msgf("bytes out: %s", repr.Repr(bytes))
			}

			byteCount, err := serialConn.Write(bytes)
			if err != nil {
				logger.Error().Msgf("Error writing packet: %s", err)
			}
			if byteCount != 4 {
				logger.Error().Msg("Packet written was not 4 bytes")
			}

			err = serialConn.Flush()
			if err != nil {
				logger.Error().Msgf("Error flushing the buffer: %s", err)
			}
		}
	}
}
