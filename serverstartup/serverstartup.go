package serverstartup

import (
	bw "github.com/HackRVA/master-base-2019/badgewrangler"
	ba "github.com/HackRVA/master-base-2019/baseapi"
	"github.com/HackRVA/master-base-2019/fifo"
	"github.com/HackRVA/master-base-2019/game"
	irp "github.com/HackRVA/master-base-2019/irpacket"
)

// StartBadgeWrangler - Start up the badge wrangler
func StartBadgeWrangler() {
	// Set up input a)nd output channels
	packetsIn := make(chan *irp.Packet)
	packetsOut := make(chan *irp.Packet)
	gameData := make(chan *bw.GameData)
	beaconHold := make(chan bool)
	game := make(chan *game.Game)

	//fifo.SetDebug(true)
	bw.SetDebug(true)

	go fifo.ReadFifo(fifo.BadgeOutFile, packetsIn)
	go fifo.WriteFifo(fifo.BadgeInFile, packetsOut)

	go bw.ReceivePackets(packetsIn, gameData, beaconHold)
	go bw.TransmitBeacon(packetsOut, beaconHold)
	go bw.TransmitNewGamePackets(packetsOut, game, beaconHold)
	go ba.DataInGameOut(gameData, game)
}
