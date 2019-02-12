package irpacket

import "fmt"

// Packet structure for badge messages
type Packet struct {
	Start   uint8
	Command uint8
	Address uint8
	BadgeID uint16
	Payload uint16
}

func BuildPacket(start uint8,
	command uint8,
	address uint8,
	badgeid uint16,
	payload uint16) *Packet {
	return &Packet{
		Start:   start,
		Command: command,
		Address: address,
		BadgeID: badgeid,
		Payload: payload}
}

func ReadPacket(rawPacket uint32) *Packet {
	return &Packet{
		Start:   uint8((rawPacket >> 31) & 0x01),
		Command: uint8((rawPacket >> 30) & 0x01),
		Address: uint8((rawPacket >> 25) & 0x1f),
		BadgeID: uint16((rawPacket >> 16) & 0x1ff),
		Payload: uint16(rawPacket & 0x0ffff)}
}

func StartBits(start uint8) uint32 {
	return ((uint32(start) & 0x01) << 31)
}

func CommandBits(command uint8) uint32 {
	return ((uint32(command) & 0x01) << 30)
}

func AddressBits(address uint8) uint32 {
	return ((uint32(address) & 0x01f) << 25)
}

func BadgeidBits(badgeid uint16) uint32 {
	return ((uint32(badgeid) & 0x1ff) << 16)
}

func PayloadBits(payload uint16) uint32 {
	return (uint32(payload) & 0x0ffff)
}

func WritePacket(packet *Packet) uint32 {
	return StartBits(packet.Start) |
		CommandBits(packet.Command) |
		AddressBits(packet.Address) |
		BadgeidBits(packet.BadgeID) |
		PayloadBits(packet.Payload)
}

func BytesToRawPacket(bytes []byte) uint32 {
	return uint32(bytes[0]) | uint32(bytes[1])<<8 | uint32(bytes[2])<<16 | uint32(bytes[3])<<24
}

func RawPacketToBytes(rawPacket uint32) []byte {
	return []byte{uint8(rawPacket & 0x0ff),
		uint8((rawPacket >> 8) & 0x0ff),
		uint8((rawPacket >> 16) & 0x0ff),
		uint8((rawPacket >> 24) & 0x0ff)}
}

func PrintPacket(packet *Packet) {
	fmt.Printf("  packet: %x\n", WritePacket(packet))
	fmt.Printf("     cmd: %#x\n", packet.Command)
	fmt.Printf("   start: %#x\n", packet.Start)
	fmt.Printf(" address: %#x\n", packet.Address)
	fmt.Printf("badge ID: %#x\n", packet.BadgeID)
	fmt.Printf(" payload: %#x\n", packet.Payload)
}
