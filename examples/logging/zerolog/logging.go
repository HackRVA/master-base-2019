package main

import (
	log "github.com/HackRVA/master-base-2019/logging"
	"github.com/rs/zerolog"
)

func main() {
	log.SetGlobalLevel(zerolog.DebugLevel)
	log.Ger.Info().Msg("first message in main")
	//foo.Foo()
	log.Ger.Info().Msg("This is a\nmultiline log.")
}

/*
fund doit() {
	gameDataC := make(chan *GameData)
	gameSpecC = make(chan *gameSpec)
	beaconHoldC = make(chan *beaconHoldC)

	go msg.ReceivePackets(badgeOutC, gameDataC, beaconHoldC)
	go msg.TransmitPackets(badgeInC, gameSpecC, beaconHoldC)
	go dst.HandOutGameSpecs(gameDataC, gameSpecC)
	for {}
}
*/
