package main

import (
	"fmt"

	"github.com/spf13/viper"
)

func main() {

	viper.SetDefault("serialPort", "/dev/ttyACM0")
	viper.SetDefault("ir", true)
	viper.SetDefault("serialDebug", false)
	viper.SetDefault("bwDebug", false)

	fmt.Println("testing configs...")
	viper.SetConfigName("baseconfig")
	viper.AddConfigPath("/etc/basestation")
	viper.AddConfigPath("$HOME/etc/basestation")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("No config file: %s\nUsing Config Defaults\n", err)
	}

	fmt.Println("serialPort:", viper.GetString("serialPort"))
	fmt.Println("ir:", viper.GetBool("ir"))
	fmt.Println("serialDebug:", viper.GetBool("serialDebug"))
	fmt.Println("badgeWranglerDebug:", viper.GetBool("bwDebug"))
}
