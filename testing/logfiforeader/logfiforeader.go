package main

import (
	"flag"

	fifo "github.com/HackRVA/master-base-2019/fifo"
	log "github.com/HackRVA/master-base-2019/filelogging"
	irp "github.com/HackRVA/master-base-2019/irpacket"
)

var namedPipe string
var logger = log.Ger

func main() {
	channelInPtr := flag.Bool("in", false, "read ingoing badge Channel (otherwise outgoing)")

	flag.Parse()

	if *channelInPtr {
		namedPipe = fifo.BadgeInFile
	} else {
		namedPipe = fifo.BadgeOutFile
	}

	packetsIn := make(chan *irp.Packet)
	go fifo.ReadFifo(namedPipe, packetsIn)
	for {
		packet := <-packetsIn
		packetLogger := packet.Logger(logger)
		packetLogger.Info().Msg("Packet received from " + fifo.BadgeOutFile)
	}
}
