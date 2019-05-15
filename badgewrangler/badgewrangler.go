package badgewrangler

// #include <../lasertag-protocol.h>
import "C"

import (
	"fmt"
	"time"

	log "github.com/HackRVA/master-base-2019/filelogging"
	gm "github.com/HackRVA/master-base-2019/game"
	irp "github.com/HackRVA/master-base-2019/irpacket"
	"github.com/hackebrot/go-repr/repr"
)

const (
	beaconInterval        = 2 * time.Second
	beaconDelay           = 5 * time.Second
	receivedPacketTimeout = 500 * time.Millisecond
)

// Values for expecting
const (
	SenderBadgeID = C.OPCODE_BADGE_IDENTITY
	GameID        = C.OPCODE_GAME_ID
	RecordCount   = C.OPCODE_BADGE_RECORD_COUNT
	BadgeID       = C.OPCODE_BADGE_UPLOAD_HIT_RECORD_BADGE_ID
	Timestamp     = C.OPCODE_BADGE_UPLOAD_HIT_RECORD_TIMESTAMP
	Team          = C.OPCODE_SET_BADGE_TEAM
	UserName      = C.OPCODE_USERNAME_DATA
)

var debug = false
var irErr = false
var logger = log.Ger.With().Str("pkg", "badgewrangler").Logger()

// SetDebug - sets the debugging on and off
func SetDebug(isDebug bool) {
	debug = isDebug
}

// SetIrErr - sets IR error simulation on and off
func SetIrErr(isIrErr bool) {
	irErr = isIrErr
}

// Hit - The data comprising a Hit
type Hit struct {
	BadgeID   uint16
	Timestamp uint16
	Team      uint8
}

// BadgeIDPacket - return a hit's badgeID packet
func (h *Hit) BadgeIDPacket() *irp.Packet {
	return BuildBadgeUploadHitRecordBadgeID(h.BadgeID)
}

// TimestampPacket - return a hit's timestamp packet
func (h *Hit) TimestampPacket() *irp.Packet {
	return BuildBadgeUploadHitRecordTimestamp(h.Timestamp)
}

// TeamPacket - return a hit's team packet
func (h *Hit) TeamPacket() *irp.Packet {
	return BuildBadgeUploadHitRecordTeam(h.Team)
}

// GameData - The game data dump from a badge
type GameData struct {
	BadgeID  uint16
	GameID   uint16
	Hits     []*Hit
	UserName string
}

// BadgeIDPacket - return gameData's BadgeID packet
func (gd *GameData) BadgeIDPacket() *irp.Packet {
	return BuildBadgeIdentity(gd.BadgeID)
}

// GameIDPacket - return a hit's gameID packet
func (gd *GameData) GameIDPacket() *irp.Packet {
	return BuildBadgeUploadHitRecordGameID(gd.GameID)
}

// HitCountPacket - return gameData's hit count packet
func (gd *GameData) HitCountPacket(hitCount uint16) *irp.Packet {
	return BuildBadgeUploadRecordCount(hitCount)
}

// Packets - return a slice containing all the gameData packets
func (gd *GameData) Packets(irErr bool) []*irp.Packet {
	packetIndex := 0
	packets := make([]*irp.Packet, len(gd.Hits)*3+3)
	packets[packetIndex] = gd.BadgeIDPacket()
	packetIndex++
	packets[packetIndex] = gd.GameIDPacket()
	packetIndex++
	packets[packetIndex] = gd.HitCountPacket(uint16(len(gd.Hits)))
	packetIndex++
	for i, hit := range gd.Hits {
		if i == 1 && irErr {
			packets[packetIndex] = irp.BuildPacket(uint16(0), C.OPCODE_SET_GAME_START_TIME<<12|uint16(19&0x0fff))
		} else {
			packets[packetIndex] = hit.BadgeIDPacket()
		}
		packetIndex++
		packets[packetIndex] = hit.TimestampPacket()
		packetIndex++
		packets[packetIndex] = hit.TeamPacket()
		packetIndex++
	}
	return packets
}

// TransmitBadgeDump - place the gameData element's packets on an outbound *Packet channel
func (gd *GameData) TransmitBadgeDump(packetsOut chan *irp.Packet, irErr bool) {
	for _, packet := range gd.Packets(irErr) {
		packetsOut <- packet
	}
}

// PrintUnexpectedPacketError - print expected vs. unexpected character error
func PrintUnexpectedPacketError(expected uint8, got uint8) {
	logger.Error().Msgf("Expected \"%s\" packet but got \"%s\" packet instead\n",
		irp.GetPayloadSpecs(expected).Description,
		irp.GetPayloadSpecs(got).Description)
}

// ReceivePackets - Receives incoming Packets, supresses beacon, and sends out GameData
func ReceivePackets(packetsIn chan *irp.Packet, gameDataOut chan *GameData, beaconHoldOut chan bool) {
	if debug {
		logger.Debug().Msg("Start processing packets")
	}
	var opcode uint8
	var expecting uint8 = SenderBadgeID
	var gameData *GameData
	var hitCount uint16
	var hitsRecorded uint16
	var startTime time.Time
	var packet *irp.Packet
	var letters uint8

	for {
		if expecting != SenderBadgeID {
			if time.Now().Sub(startTime) > receivedPacketTimeout {
				if debug {
					logger.Debug().Msg("Game dump timeout")
				}
				expecting = SenderBadgeID
				beaconHoldOut <- false
			}
		}

		select {
		case packet = <-packetsIn:
		case <-time.After(500 * time.Millisecond):
			continue
		}

		fmt.Println()
		irp.PrintPacket(packet)
		opcode = packet.Opcode()
		fmt.Println("  Opcode:", opcode)
		fmt.Println()
		switch opcode {
		case C.OPCODE_BADGE_IDENTITY:
			if expecting == SenderBadgeID {
				beaconHoldOut <- true
				startTime = time.Now()
				gameData = &GameData{
					BadgeID: uint16(packet.Payload & 0x01ff)}
				expecting = GameID
				if debug {
					logger.Debug().Msgf("** Sender Badge ID Received: %s", repr.Repr(gameData.BadgeID))
				}
			} else {
				PrintUnexpectedPacketError(expecting, opcode)
			}
		case C.OPCODE_GAME_ID:
			if expecting == GameID {
				startTime = time.Now()
				gameData.GameID = uint16(packet.Payload & 0x0fff)
				expecting = RecordCount
				if debug {
					logger.Debug().Msgf("** Game ID Received: %s", repr.Repr(gameData.GameID))
				}
			} else {
				PrintUnexpectedPacketError(expecting, opcode)
			}
		case C.OPCODE_BADGE_RECORD_COUNT:
			if expecting == RecordCount {
				startTime = time.Now()
				hitCount = uint16(packet.Payload & 0x0fff)
				hitsRecorded = 0

				gameData.Hits = make([]*Hit, hitCount)
				if hitCount > 0 {
					expecting = BadgeID
					if debug {
						logger.Debug().Msgf("** Badge Record Count Received: %s", repr.Repr(hitCount))
					}
				} else {
					expecting = UserName
					letters = 0
				}
			} else {
				PrintUnexpectedPacketError(expecting, opcode)
			}
		case C.OPCODE_BADGE_UPLOAD_HIT_RECORD_BADGE_ID:
			if expecting == BadgeID && hitsRecorded < hitCount {
				startTime = time.Now()
				hit := &Hit{
					BadgeID: uint16(packet.Payload & 0x01ff)}
				gameData.Hits[hitsRecorded] = hit
				expecting = Timestamp
				if debug {
					logger.Debug().Msgf("** Badge Upload Hit Record Badge ID Received: %s", repr.Repr(gameData.Hits[hitsRecorded].BadgeID))
				}
			} else {
				PrintUnexpectedPacketError(expecting, opcode)
			}
		case C.OPCODE_BADGE_UPLOAD_HIT_RECORD_TIMESTAMP:
			if expecting == Timestamp && hitsRecorded < hitCount {
				startTime = time.Now()
				gameData.Hits[hitsRecorded].Timestamp = uint16(packet.Payload & 0x0fff)
				expecting = Team
				if debug {
					logger.Debug().Msgf("** Badge Upload Hit Record Timestamp Received: %s", repr.Repr(gameData.Hits[hitsRecorded].Timestamp))
				}
			} else {
				PrintUnexpectedPacketError(expecting, opcode)
			}
		case C.OPCODE_SET_BADGE_TEAM:
			if expecting == Team && hitsRecorded < hitCount {
				startTime = time.Now()
				gameData.Hits[hitsRecorded].Team = uint8(packet.Payload & 0x07)
				if debug {
					logger.Debug().Msgf("** Badge Upload Hit Record Team Received: %s", repr.Repr(gameData.Hits[hitsRecorded].Team))
				}
				if hitsRecorded++; hitsRecorded == hitCount {
					expecting = UserName
					letters = 0
				} else {
					expecting = BadgeID
				}
			} else {
				PrintUnexpectedPacketError(expecting, opcode)
			}
		case C.OPCODE_USERNAME_DATA:
			if expecting == UserName && letters < 10 {
				startTime = time.Now()
				letters += 2
				if debug {
					logger.Debug().Msgf("** Badge User Name Letters Received: %d", letters)
				}
				if letters == 10 {
					if debug {
						logger.Debug().Msg("GameData Complete!")
					}
					gameDataOut <- gameData
					hitsRecorded = 0
					hitCount = 0
					gameData = nil
					expecting = SenderBadgeID
				}
			}

		default:
			{
			}
			if debug {
				logger.Error().Msgf("** Opcode %s never expected", repr.Repr(opcode))
			}
		}
	}
}

// BuildGameStartTime - Build a game start time packet
func BuildGameStartTime(game *gm.Game) *irp.Packet {
	return irp.BuildPacket(game.BadgeID, C.OPCODE_SET_GAME_START_TIME<<12|uint16(game.StartTime&0x0fff))
}

// BuildGameDuration - Build a game duration packet
func BuildGameDuration(game *gm.Game) *irp.Packet {
	return irp.BuildPacket(game.BadgeID, C.OPCODE_SET_GAME_DURATION<<12|game.Duration&0x0fff)
}

// BuildGameVariant - Build a game variant packet
func BuildGameVariant(game *gm.Game) *irp.Packet {
	return irp.BuildPacket(game.BadgeID, C.OPCODE_SET_GAME_VARIANT<<12|uint16(game.Variant))
}

// BuildGameTeam - Build a game team packet
func BuildGameTeam(game *gm.Game) *irp.Packet {
	return irp.BuildPacket(game.BadgeID, C.OPCODE_SET_BADGE_TEAM<<12|uint16(game.Team))
}

// BuildGameID - Build a game ID packet)
func BuildGameID(game *gm.Game) *irp.Packet {
	return irp.BuildPacket(game.BadgeID, C.OPCODE_GAME_ID<<12|uint16(game.GameID&0x0fff))
}

// BuildBeacon - Build the "beacon" packet
func BuildBeacon() *irp.Packet {
	return irp.BuildPacket(uint16(C.BADGE_IR_BROADCAST_ID), C.OPCODE_REQUEST_BADGE_DUMP<<12)
}

// BuildBadgeUploadHitRecordGameID - Build the game ID packet for the hit record
func BuildBadgeUploadHitRecordGameID(gameID uint16) *irp.Packet {
	return irp.BuildPacket(uint16(C.BADGE_IR_BROADCAST_ID), C.OPCODE_GAME_ID<<12|gameID&0x0fff)
}

// BuildBadgeUploadRecordCount - Build the badge record count packet
func BuildBadgeUploadRecordCount(recordCount uint16) *irp.Packet {
	return irp.BuildPacket(uint16(C.BADGE_IR_BROADCAST_ID), C.OPCODE_BADGE_RECORD_COUNT<<12|recordCount&0x0fff)
}

// BuildBadgeUploadHitRecordBadgeID - Build the badge ID packet for a hit record
func BuildBadgeUploadHitRecordBadgeID(hitBadgeID uint16) *irp.Packet {
	return irp.BuildPacket(uint16(C.BADGE_IR_BROADCAST_ID), C.OPCODE_BADGE_UPLOAD_HIT_RECORD_BADGE_ID<<12|hitBadgeID&0x01ff)
}

// BuildBadgeUploadHitRecordTeam - Build the team packet for the hit record
func BuildBadgeUploadHitRecordTeam(team uint8) *irp.Packet {
	return irp.BuildPacket(uint16(C.BADGE_IR_BROADCAST_ID), C.OPCODE_SET_BADGE_TEAM<<12|uint16(team&0x07))
}

// BuildBadgeUploadHitRecordTimestamp - Build the timestamp packet for the hit record
func BuildBadgeUploadHitRecordTimestamp(timestamp uint16) *irp.Packet {
	return irp.BuildPacket(uint16(C.BADGE_IR_BROADCAST_ID), C.OPCODE_BADGE_UPLOAD_HIT_RECORD_TIMESTAMP<<12|timestamp&0x0fff)
}

// BuildBadgeIdentity - Build the badge identity packet
func BuildBadgeIdentity(senderBadgeID uint16) *irp.Packet {
	return irp.BuildPacket(uint16(C.BADGE_IR_BROADCAST_ID), C.OPCODE_BADGE_IDENTITY<<12|senderBadgeID&0x01ff)
}

// TransmitNewGamePackets - Receives GameData, Transmits packets to the badge, and re-enables beacon
func TransmitNewGamePackets(packetsOut chan *irp.Packet, gameIn chan *gm.Game, beaconHold chan bool) {

	for {
		game := <-gameIn

		if debug {
			packetLogger := game.Logger(logger)
			packetLogger.Debug().Msg("Send game to badge")
		}

		NewGamePackets(packetsOut, game)
		NewGamePackets(packetsOut, game)

		time.Sleep(beaconDelay)

		beaconHold <- false
	}
}

// NewGamePackets - Put the new game packets on the channel
func NewGamePackets(packetsOut chan *irp.Packet, game *gm.Game) {
	packetsOut <- BuildGameStartTime(game)
	packetsOut <- BuildGameDuration(game)
	packetsOut <- BuildGameVariant(game)
	packetsOut <- BuildGameTeam(game)
	packetsOut <- BuildGameID(game)
}

// TransmitBeacon - Transmits "beacon" packets to the badge to trigger gameData upload
//                  Switchable based on input from beaconHoldIn channel
func TransmitBeacon(packetsOut chan *irp.Packet, beaconHoldIn chan bool) {

	beaconHold := false
	if debug {
		logger.Debug().Msg("Beacon is on")
	}
	for {
		select {
		case beaconHold = <-beaconHoldIn:
			if debug {
				status := "on"
				if beaconHold {
					status = "off"
				}
				logger.Debug().Msg("Beacon is " + status)
			}
		default:
		}
		if !beaconHold {
			packetsOut <- BuildBeacon()
			time.Sleep(beaconInterval)
		}
	}
}

// BadgeHandlePackets - packet handler for the badge simulator
func BadgeHandlePackets(packetsIn chan *irp.Packet, packetsOut chan *irp.Packet, gameData *GameData) {
	if debug {
		logger.Debug().Msg("Start handling packets")
	}
	var opcode uint8

	for {
		packet := <-packetsIn
		opcode = packet.Opcode()
		desc, data := packet.PayloadDescData()

		switch opcode {
		case C.OPCODE_REQUEST_BADGE_DUMP:
			gameData.TransmitBadgeDump(packetsOut, irErr)
		// Game Start Time
		case C.OPCODE_SET_GAME_START_TIME:
			logger.Debug().Msgf("[%s] packet received, payload: %d", desc, data)
		// Game Duration
		case C.OPCODE_SET_GAME_DURATION:
			logger.Debug().Msgf("[%s] packet received, payload: %d", desc, data)
		// Game Variant
		case C.OPCODE_SET_GAME_VARIANT:
			logger.Debug().Msgf("[%s] packet received, payload: %d", desc, data)
		// Game Team
		case C.OPCODE_SET_BADGE_TEAM:
			logger.Debug().Msgf("[%s] packet received, payload: %d", desc, data)
		// Game ID
		case C.OPCODE_GAME_ID:
			logger.Debug().Msgf("[%s] packet received, payload: %d", desc, data)
		default:
			if debug {
				logger.Debug().Uint8("opcode", opcode).Msgf("[%s] packet not handled yet.", desc)
			}
		}
	}
}
