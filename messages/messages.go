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
		case C.OPCODE_BADGE_RECORD_COUNT:
			hitCount := uint16(packet.Payload & 0x0fff)
			fmt.Println("** Badge Record Count Received:", hitCount)
		case C.OPCODE_BADGE_UPLOAD_HIT_RECORD_BADGE_ID:
			hitBadgeID := uint16(packet.Payload & 0x01ff)
			fmt.Println("** Badge Upload Hit Record Badge ID Received:", hitBadgeID)
		case C.OPCODE_BADGE_UPLOAD_HIT_RECORD_TIMESTAMP:
			hitTimestamp := uint16(packet.Payload & 0x0fff)
			fmt.Println("** Badge Upload Hit Record Timestamp Received:", hitTimestamp)
		case C.OPCODE_SET_BADGE_TEAM:
			hitTeam := uint16(packet.Payload & 0x0fff)
			fmt.Println("** Badge Upload Hit Record Team Received:", hitTeam)
		default:
			fmt.Println("** Opcode", opcode, "not handled yet")
		}
	}
}
