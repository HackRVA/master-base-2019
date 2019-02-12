package messages

// #include <badge-ir-game-protocol.h>
import "C"

import (
	"fmt"

	irp "github.com/HackRVA/master-base-2019/irpacket"
)

func ProcessPackets(packetsIn chan *irp.Packet) {
	fmt.Println("Start processing packets")
	var opcode uint8

	for {
		packet := <-packetsIn
		fmt.Println()
		irp.PrintPacket(packet)
		opcode = uint8(packet.Payload >> 12)
		fmt.Println("  Opcode:", opcode)
		fmt.Println()
		switch opcode {
		case C.OPCODE_GAME_ID:
			gameID := uint16(packet.Payload & 0x0fff)
			fmt.Println("** Game ID Received:", gameID)
		default:
			fmt.Println("** Opcode", opcode, "not handled yet")
		}
	}
}
