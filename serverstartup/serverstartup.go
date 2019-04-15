package serverstartup

import (
	bw "github.com/HackRVA/master-base-2019/badgewrangler"
	ba "github.com/HackRVA/master-base-2019/baseapi"
	"github.com/HackRVA/master-base-2019/game"
	irp "github.com/HackRVA/master-base-2019/irpacket"
	"github.com/HackRVA/master-base-2019/serial"
)

// StartBadgeWrangler - Start up the badge wrangler
func StartBadgeWrangler(port string, baud int) {
	// Set up input a)nd output channels
	packetsIn := make(chan *irp.Packet)
	packetsOut := make(chan *irp.Packet)
	gameDataIn := make(chan *bw.GameData)
	gameDataOut := make(chan *bw.GameData)
	beaconHold := make(chan bool)
	gameOut := make(chan *game.Game)

	serial.SetDebug(true)
	bw.SetDebug(true)

	serial.OpenPort(port, baud)

	go serial.ReadSerial(packetsIn)
	go serial.WriteSerial(packetsOut)

	go bw.ReceivePackets(packetsIn, gameDataIn, beaconHold)
	go bw.TransmitBeacon(packetsOut, beaconHold)
	go bw.TransmitNewGamePackets(packetsOut, gameOut, beaconHold)
	go ba.DataInGameOut(gameDataIn, gameDataOut, gameOut)
	go ba.SendGameData(gameDataOut)
}
