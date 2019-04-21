package serverstartup

import (
	"fmt"

	bw "github.com/HackRVA/master-base-2019/badgewrangler"
	ba "github.com/HackRVA/master-base-2019/baseapi"
	"github.com/HackRVA/master-base-2019/game"
	irp "github.com/HackRVA/master-base-2019/irpacket"
	"github.com/HackRVA/master-base-2019/serial"
	"github.com/spf13/viper"
)

// StartBadgeWrangler - Start up the badge wrangler
func StartBadgeWrangler() {

	// Config init
	fmt.Println("Configuration Settings...")
	viper.SetDefault("serialPort", "/dev/ttyACM0")
	viper.SetDefault("baud", 9600)
	viper.SetDefault("ir", true)
	viper.SetDefault("serialDebug", false)
	viper.SetDefault("bwDebug", false)

	viper.SetConfigName("baseconfig")
	viper.AddConfigPath("/etc/basestation")
	viper.AddConfigPath("$HOME/etc/basestation")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("No config file: %s\nUsing Config Defaults\n", err)
	}

	fmt.Println("        serialPort:", viper.GetString("serialPort"))
	fmt.Println("              baud:", viper.GetInt("baud"))
	fmt.Println("                ir:", viper.GetBool("ir"))
	fmt.Println("       serialDebug:", viper.GetBool("serialDebug"))
	fmt.Println("badgeWranglerDebug:", viper.GetBool("bwDebug"))

	// Set up input a)nd output channels
	packetsIn := make(chan *irp.Packet)
	packetsOut := make(chan *irp.Packet)
	gameDataIn := make(chan *bw.GameData)
	gameDataOut := make(chan *bw.GameData)
	beaconHold := make(chan bool)
	gameOut := make(chan *game.Game)

	serial.SetDebug(viper.GetBool("serialDebug"))
	bw.SetDebug(viper.GetBool("bwDebug"))

	serial.OpenPort(viper.GetString("serialPort"), viper.GetInt("baud"))
	if viper.GetBool("ir") {
		serial.InitIR()
	}

	go serial.ReadSerial(packetsIn)
	go serial.WriteSerial(packetsOut)

	go bw.ReceivePackets(packetsIn, gameDataIn, beaconHold)
	go bw.TransmitBeacon(packetsOut, beaconHold)
	go bw.TransmitNewGamePackets(packetsOut, gameOut, beaconHold)
	go ba.DataInGameOut(gameDataIn, gameDataOut, gameOut)
	go ba.SendGameData(gameDataOut)
}
