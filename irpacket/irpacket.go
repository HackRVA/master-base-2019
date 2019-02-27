package irpacket

// #include <../badge-ir-game-protocol.h>
import "C"

import (
	"fmt"
	"strconv"

	log "github.com/HackRVA/master-base-2019/logging"
	utl "github.com/HackRVA/master-base-2019/utility"
	zl "github.com/rs/zerolog"
)

// Start - default Start value
const Start = 1

// Command - default Command value
const Command = 1

// PayloadSpec - The data describing a payload type
type PayloadSpec struct {
	Opcode      uint8
	Description string
	Mask        uint16
}

// PayloadSpecList - Slice of Payload Structures
var payloadSpecList = []PayloadSpec{
	{C.OPCODE_SET_GAME_START_TIME, "Set Game Start Time", 0x0fff},
	{C.OPCODE_SET_GAME_DURATION, "St Game Duration", 0x0fff},
	{C.OPCODE_HIT, "Badge Hit (team)", 0x0f},
	{C.OPCODE_SET_BADGE_TEAM, "Set Badge Team", 0x0f},
	{C.OPCODE_REQUEST_BADGE_DUMP, "Request Badge Dump", 0x0},
	{C.OPCODE_SET_GAME_VARIANT, "Set Game Variant", 0x0f},
	{C.OPCODE_BADGE_RECORD_COUNT, "Badge Record Count", 0x0fff},
	{C.OPCODE_BADGE_UPLOAD_HIT_RECORD_BADGE_ID, "Badge Upload Hit Record Badge ID", 0x01ff},
	{C.OPCODE_GAME_ID, "Game ID", 0x0fff},
	{C.OPCODE_BADGE_UPLOAD_HIT_RECORD_TIMESTAMP, "Badge Upload Hit Record Timestamp", 0x0fff}}

var payloadSpecMap map[uint8]PayloadSpec

func init() {
	payloadSpecMap = make(map[uint8]PayloadSpec)

	for _, payload := range payloadSpecList {
		payloadSpecMap[payload.Opcode] = payload
	}
}

// GetPayloadSpecs - returns a PayloadData
func GetPayloadSpecs(opcode uint8) PayloadSpec {
	return payloadSpecMap[opcode]
}

// RawPacket - the unsigned integer representing the raw badge packet
type RawPacket uint32

// ReadPacket - read a packet from a rawPacket
func ReadPacket(rawPacket RawPacket) *Packet {
	return &Packet{
		Start:   uint8((rawPacket >> 31) & 0x01),
		Command: uint8((rawPacket >> 30) & 0x01),
		Address: uint8((rawPacket >> 25) & 0x1f),
		BadgeID: uint16((rawPacket >> 16) & 0x1ff),
		Payload: uint16(rawPacket & 0x0ffff)}
}

// RawPacketToBytes - convert a rawPacket to a four byte array
func RawPacketToBytes(rawPacket RawPacket) []byte {
	return []byte{uint8(rawPacket & 0x0ff),
		uint8((rawPacket >> 8) & 0x0ff),
		uint8((rawPacket >> 16) & 0x0ff),
		uint8((rawPacket >> 24) & 0x0ff)}
}

// Packet - return a Packet
func (r RawPacket) Packet() *Packet {
	return ReadPacket(r)
}

// Bytes - return the PacketBytes
func (r RawPacket) Bytes() PacketBytes {
	return r.Packet().Bytes()
}

// PacketBytes - the bytes containing the raw packet
type PacketBytes []byte

// BytesToRawPacket - convert four byte array to a RawPacket
func BytesToRawPacket(bytes PacketBytes) RawPacket {
	return RawPacket(uint32(bytes[0]) | uint32(bytes[1])<<8 | uint32(bytes[2])<<16 | uint32(bytes[3])<<24)
}

// RawPacket - return the RawPacket
func (b PacketBytes) RawPacket() RawPacket {
	return BytesToRawPacket(b)
}

// Packet - return the Packet
func (b PacketBytes) Packet() *Packet {
	return b.RawPacket().Packet()
}

// Packet structure for badge messages
type Packet struct {
	Start   uint8
	Command uint8
	Address uint8
	BadgeID uint16
	Payload uint16
}

// BuildPacket - Build a packet
func BuildPacket(
	badgeid uint16, payload uint16) *Packet {
	return &Packet{
		Start:   Start,
		Command: Command,
		Address: C.BADGE_IR_GAME_ADDRESS,
		BadgeID: badgeid,
		Payload: payload}
}

// StartBit - return a rawPacket with placed Start bit
func StartBit(start uint8) RawPacket {
	return ((RawPacket(start) & 0x01) << 31)
}

// CommandBit - return a rawPacket with placed Command bit
func CommandBit(command uint8) RawPacket {
	return ((RawPacket(command) & 0x01) << 30)
}

// AddressBits - return a rawPacket with placed Address bits
func AddressBits(address uint8) RawPacket {
	return ((RawPacket(address) & 0x01f) << 25)
}

// BadgeIDBits - return a rawPacket with placed BadgeID bits
func BadgeIDBits(badgeid uint16) RawPacket {
	return ((RawPacket(badgeid) & 0x1ff) << 16)
}

// PayloadBits - return a rawPacket with placed Payload bits
func PayloadBits(payload uint16) RawPacket {
	return (RawPacket(payload) & 0x0ffff)
}

// WritePacket - return a rawPacket from a packet
func WritePacket(packet *Packet) RawPacket {
	return RawPacket(StartBit(packet.Start) |
		CommandBit(packet.Command) |
		AddressBits(packet.Address) |
		BadgeIDBits(packet.BadgeID) |
		PayloadBits(packet.Payload))
}

// PrintPacket - print out a packet's contents
func PrintPacket(packet *Packet) {
	fmt.Printf("  packet: %x\n", WritePacket(packet))
	fmt.Printf("     cmd: %#6x - %6[1]d\n", packet.Command)
	fmt.Printf("   start: %#6x - %6[1]d\n", packet.Start)
	fmt.Printf(" address: %#6x - %6[1]d\n", packet.Address)
	fmt.Printf("badge ID: %#6x - %6[1]d\n", packet.BadgeID)
	fmt.Printf(" payload: %#6x - %6[1]d\n", packet.Payload)
}

// Logger - return a packet sublogger
func Logger(packet *Packet) zl.Logger {
	return log.Ger.With().
		Uint8("cmd", packet.Command).
		Uint8("start", packet.Start).
		Uint8("address", packet.Address).
		Uint16("badge ID", packet.BadgeID).
		Uint16("payload", packet.Payload).
		Str("opcode", packet.OpcodeDescription()).
		Int16("payload data", packet.PayloadData()).Logger()
}

// PrintPayload - Print out the payload particulars
func PrintPayload(packet *Packet) {
	opcode := uint8(packet.Payload >> 12)
	pd := payloadSpecMap[opcode]
	fmt.Println(strconv.Itoa(int(opcode))+":"+pd.Description+":", packet.Payload&pd.Mask)
}

// GetPayload - Returns the description and value of a packet's payload
func GetPayload(packet *Packet) (string, int16) {
	opcode := uint8(packet.Payload >> 12)
	pd := payloadSpecMap[opcode]
	uintPayload := packet.Payload & pd.Mask
	if opcode == C.OPCODE_SET_GAME_START_TIME {
		return pd.Description, utl.Int12fromUint16toInt16(uintPayload)
	}
	return pd.Description, int16(uintPayload)
}

// Print - method to print a packet's contents
func (p Packet) Print() { PrintPacket(&p) }

// Opcode - Return Opcode of a packet's payload
func (p Packet) Opcode() uint8 {
	return uint8(p.Payload >> 12)
}

// PayloadData - Return payload data from a packet
func (p *Packet) PayloadData() int16 {
	_, payload := GetPayload(p)
	return payload
}

// OpcodeDescription - Return the opcode descriptin from a packet
func (p *Packet) OpcodeDescription() string {
	return GetPayloadSpecs(p.Opcode()).Description
}

// RawPacket - return a packet's rawPacket
func (p Packet) RawPacket() RawPacket {
	return WritePacket(&p)
}

// Bytes - return a packet's raw bytes
func (p Packet) Bytes() []byte {
	return RawPacketToBytes(WritePacket(&p))
}

// PrintPayload - prints a packets' payload opcode and value
func (p Packet) PrintPayload() {
	PrintPayload(&p)
}
