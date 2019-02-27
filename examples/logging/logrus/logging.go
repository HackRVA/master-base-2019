package main

import (
	"os"

	foo "github.com/HackRVA/master-base-2019/examples/logging/logrus/testpkglog"
	log "github.com/sirupsen/logrus"
)

var logger *log.Entry

func init() {
	log.SetOutput(os.Stdout)

	file, err := os.OpenFile("test.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Debug("Failed to log to file, using default stderr")
	}
	logger = log.WithFields(log.Fields{"pkg": "main"})
}

func main() {
	log.SetLevel(log.DebugLevel)
	logger.Info("first message in main")
	foo.Foo()
	logger.Info("This is a\nmultiline log.")
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
