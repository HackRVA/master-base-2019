package messages

// #include <../badge-ir-game-protocol.h>
import "C"

import (
	"fmt"
	"time"

	irp "github.com/HackRVA/master-base-2019/irpacket"
)

const (
	beaconInterval = 2 * time.Second
	beaconDelay    = 15 * time.Second
)

// Values for expecting
const (
	GameID      = C.OPCODE_GAME_ID
	RecordCount = C.OPCODE_BADGE_RECORD_COUNT
	BadgeID     = C.OPCODE_BADGE_UPLOAD_HIT_RECORD_BADGE_ID
	Timestamp   = C.OPCODE_BADGE_UPLOAD_HIT_RECORD_TIMESTAMP
	Team        = C.OPCODE_SET_BADGE_TEAM
)

var debug = false

// SetDebug - sets the debugging on and off
func SetDebug(isDebug bool) {
	debug = isDebug
}

// Hit - The data comprising a Hit
type Hit struct {
	BadgeID   uint16
	Timestamp uint16
	Team      uint8
}

func (h *Hit) BadgeIDPacket(badgeID uint16) *irp.Packet {
	return BuildBadgeUploadHitRecordBadgeID(badgeID, h.BadgeID)
}

func (h *Hit) TimestampPacket(badgeID uint16) *irp.Packet {
	return BuildBadgeUploadHitRecordTimestamp(badgeID, h.Timestamp)
}

func (h *Hit) TeamPacket(badgeID uint16) *irp.Packet {
	return BuildBadgeUploadHitRecordTeam(badgeID, h.Team)
}

// GameData - The game data dump from a badge
type GameData struct {
	BadgeID uint16
	GameID  uint16
	Hits    []*Hit
}

func (gd *GameData) BadgeIDPacket() *irp.Packet {
	return BuildBadgeUploadHitRecordGameID(gd.BadgeID, gd.GameID)
}

func (gd *GameData) HitCountPacket(hitCount uint16) *irp.Packet {
	return BuildBadgeUploadRecordCount(gd.BadgeID, hitCount)
}

func (gd *GameData) Packets() []*irp.Packet {
	packets := make([]*irp.Packet, len(gd.Hits)+2)
	packets = append(packets, gd.BadgeIDPacket())
	packets = append(packets, gd.HitCountPacket(uint16(len(gd.Hits))))
	for _, hit := range gd.Hits {
		packets = append(packets, hit.BadgeIDPacket(gd.BadgeID))
		packets = append(packets, hit.TimestampPacket(gd.BadgeID))
		packets = append(packets, hit.TeamPacket(gd.BadgeID))
	}
	return packets
}

func (gd *GameData) TransmitBadgeDump(packetsOut chan *irp.Packet) {
	for _, packet := range gd.Packets() {
		packetsOut <- packet
	}
}

// GameSpec - The game specification sent to the badge
type GameSpec struct {
	BadgeID   uint16
	StartTime int16
	Duration  uint16 // 0x0fff
	Variant   uint8  // 0x0f
	Team      uint8  // 0x0f
	GameID    uint16 // 0x0fff
}

// PrintUnexpectedPacketError - print expected vs. unexpected character error
func PrintUnexpectedPacketError(expected uint8, got uint8) {
	fmt.Printf("Expected \"%s\" packet but got \"%s\" packet instead\n",
		irp.GetPayloadSpecs(expected).Description,
		irp.GetPayloadSpecs(got).Description)
}

// ReceivePackets - Receives incoming Packets, supresses beacon, and sends out GameData
func ReceivePackets(packetsIn chan *irp.Packet, gameDataOut chan *GameData, beaconHoldOut chan bool) {
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
				beaconHoldOut <- false
			}
		}

		packet := <-packetsIn
		fmt.Println()
		irp.PrintPacket(packet)
		opcode = packet.Opcode()
		fmt.Println("  Opcode:", opcode)
		fmt.Println()
		switch opcode {
		case C.OPCODE_GAME_ID:
			if expecting == GameID {
				beaconHoldOut <- true
				startTime = time.Now()
				gameData = &GameData{
					BadgeID: packet.BadgeID,
					GameID:  uint16(packet.Payload & 0x0fff)}
				expecting = RecordCount
				if debug {
					fmt.Println("** Game ID Received:", gameData.GameID)
				}
			} else {
				PrintUnexpectedPacketError(expecting, opcode)
			}
		case C.OPCODE_BADGE_RECORD_COUNT:
			if expecting == RecordCount {
				hitCount = uint16(packet.Payload & 0x0fff)
				hitsRecorded = 0
				if debug {
					fmt.Println("** Badge Record Count Received:", hitCount)
				}
			} else {
				PrintUnexpectedPacketError(expecting, opcode)
			}
		case C.OPCODE_BADGE_UPLOAD_HIT_RECORD_BADGE_ID:
			if expecting == BadgeID && hitsRecorded < hitCount {
				hit := &Hit{
					BadgeID: uint16(packet.Payload & 0x01ff)}
				gameData.Hits[hitsRecorded] = hit
				expecting = Timestamp
				if debug {
					fmt.Println("** Badge Upload Hit Record Badge ID Received:", gameData.Hits[hitsRecorded].BadgeID)
				}
			} else {
				PrintUnexpectedPacketError(expecting, opcode)
			}
		case C.OPCODE_BADGE_UPLOAD_HIT_RECORD_TIMESTAMP:
			if expecting == Timestamp && hitsRecorded < hitCount {
				gameData.Hits[hitsRecorded].Timestamp = uint16(packet.Payload & 0x0fff)
				expecting = Team
				if debug {
					fmt.Println("** Badge Upload Hit Record Timestamp Received:", gameData.Hits[hitsRecorded].Timestamp)
				}
			} else {
				PrintUnexpectedPacketError(expecting, opcode)
			}
		case C.OPCODE_SET_BADGE_TEAM:
			if expecting == Team && hitsRecorded < hitCount {
				gameData.Hits[hitsRecorded].Team = uint8(packet.Payload & 0x0fff)
				if debug {
					fmt.Println("** Badge Upload Hit Record Team Received:", gameData.Hits[hitsRecorded].Team)
				}
				if hitsRecorded++; hitsRecorded == hitCount {
					if debug {
						fmt.Println("GameData Complete!")
					}
					gameDataOut <- gameData
					hitsRecorded = 0
					hitCount = 0
					gameData = nil
					expecting = GameID
				} else {
					expecting = BadgeID
				}
			} else {
				PrintUnexpectedPacketError(expecting, opcode)

			}
		default:
			fmt.Println("** Opcode", opcode, "not handled yet")
		}
	}
}

// BuildGameStartTime - Build a game start time packet
func BuildGameStartTime(gameSpec *GameSpec) *irp.Packet {
	return irp.BuildPacket(gameSpec.BadgeID, C.OPCODE_SET_GAME_START_TIME<<12|uint16(gameSpec.StartTime&0x0fff))
}

// BuildGameDuration - Build a game duration packet
func BuildGameDuration(gameSpec *GameSpec) *irp.Packet {
	return irp.BuildPacket(gameSpec.BadgeID, C.OPCODE_SET_GAME_DURATION<<12|gameSpec.Duration&0x0fff)
}

// BuildGameVariant - Build a game variant packet
func BuildGameVariant(gameSpec *GameSpec) *irp.Packet {
	return irp.BuildPacket(gameSpec.BadgeID, C.OPCODE_SET_GAME_VARIANT<<12|uint16(gameSpec.Variant))
}

// BuildGameTeam - Build a game team packet
func BuildGameTeam(gameSpec *GameSpec) *irp.Packet {
	return irp.BuildPacket(gameSpec.BadgeID, C.OPCODE_SET_BADGE_TEAM<<12|uint16(gameSpec.Team))
}

// BuildGameID - Build a game ID packet
func BuildGameID(gameSpec *GameSpec) *irp.Packet {
	return irp.BuildPacket(gameSpec.BadgeID, C.OPCODE_GAME_ID<<12|uint16(gameSpec.GameID&0x0fff))
}

// BuildBeacon - Build the "beacon" packet
func BuildBeacon() *irp.Packet {
	return irp.BuildPacket(uint16(C.BASE_STATION_BADGE_ID), C.OPCODE_REQUEST_BADGE_DUMP<<12)
}

// BuildBadgeUploadHitRecordGameID - Build the game ID packet for the hit record
func BuildBadgeUploadHitRecordGameID(badgeID uint16, gameID uint16) *irp.Packet {
	return irp.BuildPacket(badgeID, C.OPCODE_GAME_ID|gameID&0x0fff)
}

// BuildBadgeUploadRecordCount - Build the badge record count packet
func BuildBadgeUploadRecordCount(badgeID uint16, recordCount uint16) *irp.Packet {
	return irp.BuildPacket(badgeID, C.OPCODE_BADGE_RECORD_COUNT<<12|recordCount&0x0fff)
}

// BuildBadgeUploadHitRecordBadgeID - Build the badge ID packet for a hit record
func BuildBadgeUploadHitRecordBadgeID(badgeID uint16, hitBadgeID uint16) *irp.Packet {
	return irp.BuildPacket(badgeID, C.OPCODE_BADGE_UPLOAD_HIT_RECORD_BADGE_ID<<12|hitBadgeID&0x01ff)
}

// BuildBadgeUploadHitRecordTeam - Build the team packet for the hit record
func BuildBadgeUploadHitRecordTeam(badgeID uint16, team uint8) *irp.Packet {
	return irp.BuildPacket(badgeID, C.OPCODE_SET_BADGE_TEAM|uint16(team&0x0f))
}

// BuildBadgeUploadHitRecordTimestamp - Build the timestamp packet for the hit record
func BuildBadgeUploadHitRecordTimestamp(badgeID uint16, timestamp uint16) *irp.Packet {
	return irp.BuildPacket(badgeID, C.OPCODE_BADGE_UPLOAD_HIT_RECORD_TIMESTAMP|timestamp&0x0fff)
}

// TransmitNewGamePackets - Receives GameData, Transmits packets to the badge, and re-enables beacon
func TransmitNewGamePackets(packetsOut chan *irp.Packet, gameSpecIn chan *GameSpec, beaconHold chan bool) {

	for {
		gameSpec := <-gameSpecIn

		packetsOut <- BuildGameStartTime(gameSpec)
		packetsOut <- BuildGameDuration(gameSpec)
		packetsOut <- BuildGameVariant(gameSpec)
		packetsOut <- BuildGameTeam(gameSpec)
		packetsOut <- BuildGameID(gameSpec)

		time.Sleep(beaconDelay)

		beaconHold <- false
	}
}

// TransmitBeacon - Transmits "beacon" packets to the badge to trigger gameData upload
//                  Switchable based on input from beaconHoldIn channel
func TransmitBeacon(packetsOut chan *irp.Packet, beaconHoldIn chan bool) {

	beaconHold := false
	for {
		select {
		case beaconHold = <-beaconHoldIn:
		default:
		}
		if !beaconHold {
			packetsOut <- BuildBeacon()
			time.Sleep(beaconInterval)
		}
	}
}

// BadgeHandlePackets - packet handler for the badge simulator
func BadgeHandlePackets(packetsIn chan *irp.Packet, packetsOut chan *irp.Packet, beaconIgnoreIn chan bool, gameData *GameData) {
	fmt.Println("Start handling packets")
	// beaconIgnoredChan := make(chan bool)
	var opcode uint8
	// isBeaconIgnored := false

	// timer := time.NewTimer(beaconDelay)
	// timer.Stop()

	/*
		go func(beaconIgnored chan bool) {
			<-timer.C
			beaconIgnored <- true
			if debug {
				fmt.Println("beacon ignore put enchanneled")
			}
		}(beaconIgnoredChan)
	*/

	for {
		packet := <-packetsIn
		opcode = packet.Opcode()

		/*
			select {
			case isBeaconIgnored = <-beaconIgnoredChan:
			default:
			}
		*/

		switch opcode {
		case C.OPCODE_REQUEST_BADGE_DUMP:
			//		if !isBeaconIgnored {
			gameData.TransmitBadgeDump(packetsOut)
			//			isBeaconIgnored = true
			//			timer.Reset(beaconDelay)
			//		}
		default:
			fmt.Printf("\"%s\" packet not handled yet.\n", irp.GetPayloadSpecs(opcode).Description)
		}

	}
}
