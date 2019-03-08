package serverstartup

import (
	bw "github.com/HackRVA/master-base-2019/badgewrangler"
	ba "github.com/HackRVA/master-base-2019/baseapi"
	"github.com/HackRVA/master-base-2019/fifo"
	irp "github.com/HackRVA/master-base-2019/irpacket"
)

// StartBadgeWrangler - Start up the badge wrangler
func StartBadgeWrangler() {
	// Set up input a)nd output channels
	packetsIn := make(chan *irp.Packet)
	packetsOut := make(chan *irp.Packet)
	gameData := make(chan *bw.GameData)
	beaconHold := make(chan bool)
	game := make(chan *bw.Game)

	go fifo.ReadFifo(fifo.BadgeOutFile, packetsIn)
	go fifo.WriteFifo(fifo.BadgeInFile, packetsOut)
	fifo.SetDebug(true)

	go bw.ReceivePackets(packetsIn, gameData, beaconHold)
	go bw.TransmitBeacon(packetsOut, beaconHold)
	go bw.TransmitNewGamePackets(packetsOut, game, beaconHold)
	go ba.DataInGameOut(gameData, game)
}
