package serial

import (
	log "github.com/HackRVA/master-base-2019/filelogging"
	irp "github.com/HackRVA/master-base-2019/irpacket"
	"github.com/hackebrot/go-repr/repr"
	"github.com/tarm/serial"
)

var debug = false
var logger = log.Ger

// SetDebug - set debug on/off
func SetDebug(isDebug bool) {
	debug = isDebug
}

var portName = "comX"
var baud = 9600

var serialConn *serial.Port

// OpenPort - Open a serial port
func OpenPort(portname string, baud int) {
	var err error

	config := &serial.Config{Name: portName, Baud: baud}
	serialConn, err = serial.OpenPort(config)
	if err != nil {
		logger.Fatal().Err(err)
	}
}

//ReadSerial - Reads a badge packet from a serial port
func ReadSerial(packetsIn chan *irp.Packet) {

	buf := make([]byte, 128)

	for {
		byteCount, err := serialConn.Read(buf)
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
			packetLogger.Debug().Msgf("Packet read from serial and routed to channel")
		}

		packetsIn <- packet
	}

}

// WriteSerial - writes a packet to a serial port
func WriteSerial(packetsOut chan *irp.Packet) {

	for {
		packet := <-packetsOut

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

	}

}
