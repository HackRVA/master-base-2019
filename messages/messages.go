package messages

// #include <badge-ir-game-protocol.h>
import "C"

import (
	"fmt"
	"time"

	irp "github.com/HackRVA/master-base-2019/irpacket"
)

// Values for expecting
const (
	GameID      = iota
	RecordCount = iota
	BadgeID     = iota
	Timestamp   = iota
	Team        = iota
)

// Hit - The data comprising a Hit
type Hit struct {
	BadgeID   uint16
	Timestamp uint16
	Team      uint8
}

// GameData - The game data dump from a badge
type GameData struct {
	BadgeID uint16
	GameID  uint16
	Hits    []*Hit
}

// GameSpec - The game specification sent to the badge
type GameSpec struct {
	BadgeID   uint16
	StartTime uint16
	Duration  uint16
	Variant   uint8 // 0x0f
	Team      uint8 // 0x0f
	GameID    uint8 // 0x0f
}

// ProcessPackets - Processes incoming Packets, supresses beacon, and sends out GameData
func ProcessPackets(packetsIn chan *irp.Packet, gameDataOut chan *GameData, beaconHold chan bool) {
	fmt.Println("Start processing packets")
	var opcode uint8
	var expecting uint8 = GameID
	var gameData *GameData
	var hitCount uint16
	var hitsRecorded uint16
	var startTime time.Time

	for {
		if expecting > GameID {
			elapsedTime := time.Now()
			timeoutInterval, _ := time.ParseDuration("2s")
			if elapsedTime.Sub(startTime) > timeoutInterval {
				expecting = GameID
			}
		}

		packet := <-packetsIn
		fmt.Println()
		irp.PrintPacket(packet)
		opcode = uint8(packet.Payload >> 12)
		fmt.Println("  Opcode:", opcode)
		fmt.Println()
		switch opcode {
		case C.OPCODE_GAME_ID:
			if expecting == GameID {
				beaconHold <- true
				startTime = time.Now()
				gameData = &GameData{
					BadgeID: packet.BadgeID,
					GameID:  uint16(packet.Payload & 0x0fff)}
				expecting = RecordCount
				fmt.Println("** Game ID Received:", gameData.GameID)
			}
		case C.OPCODE_BADGE_RECORD_COUNT:
			if expecting == RecordCount {
				hitCount = uint16(packet.Payload & 0x0fff)
				hitsRecorded = 0
				fmt.Println("** Badge Record Count Received:", hitCount)
			}
		case C.OPCODE_BADGE_UPLOAD_HIT_RECORD_BADGE_ID:
			if expecting == BadgeID && hitsRecorded < hitCount {
				hit := &Hit{
					BadgeID: uint16(packet.Payload & 0x01ff)}
				gameData.Hits[hitsRecorded] = hit
				expecting = Timestamp
				fmt.Println("** Badge Upload Hit Record Badge ID Received:", gameData.Hits[hitsRecorded].BadgeID)
			}
		case C.OPCODE_BADGE_UPLOAD_HIT_RECORD_TIMESTAMP:
			if expecting == Timestamp && hitsRecorded < hitCount {
				gameData.Hits[hitsRecorded].Timestamp = uint16(packet.Payload & 0x0fff)
				expecting = Team
				fmt.Println("** Badge Upload Hit Record Timestamp Received:", gameData.Hits[hitsRecorded].Timestamp)
			}
		case C.OPCODE_SET_BADGE_TEAM:
			if expecting == Team && hitsRecorded < hitCount {
				gameData.Hits[hitsRecorded].Team = uint8(packet.Payload & 0x0fff)
				fmt.Println("** Badge Upload Hit Record Team Received:", gameData.Hits[hitsRecorded].Team)
				if hitsRecorded++; hitsRecorded == hitCount {
					fmt.Println("GameData Complete!")
					gameDataOut <- gameData
					hitsRecorded = 0
					hitCount = 0
					gameData = nil
					expecting = GameID
					beaconHold <- false
				} else {
					expecting = BadgeID
				}
			}
		default:
			fmt.Println("** Opcode", opcode, "not handled yet")
		}
	}
}

func buildGameStartTime(gameSpec *GameSpec) *irp.Packet {
	return irp.BuildPacket(1, 1, C.BADGE_IR_GAME_ADDRESS, gameSpec.BadgeID, C.OPCODE_SET_GAME_START_TIME<<12|gameSpec.StartTime)
}
