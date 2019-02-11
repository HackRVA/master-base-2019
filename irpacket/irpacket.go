package irpacket

// Packet structure for badge messages
type Packet struct {
	Start   uint8
	Command uint8
	Address uint8
	BadgeID uint16
	Payload uint16
}

func buildPacket(start uint8,
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

func readPacket(rawPacket uint32) *Packet {
	return &Packet{
		Start:   uint8((rawPacket >> 31) & 0x01),
		Command: uint8((rawPacket >> 30) & 0x01),
		Address: uint8((rawPacket >> 25) & 0x1f),
		BadgeID: uint16((rawPacket >> 16) & 0x1ff),
		Payload: uint16(rawPacket & 0x0ffff)}
}

func startBits(start uint8) uint32 {
	return ((uint32(start) & 0x01) << 31)
}

func commandBits(command uint8) uint32 {
	return ((uint32(command) & 0x01) << 30)
}

func addressBits(address uint8) uint32 {
	return ((uint32(address) & 0x01f) << 25)
}

func badgeidBits(badgeid uint16) uint32 {
	return ((uint32(badgeid) & 0x1ff) << 16)
}

func payloadBits(payload uint16) uint32 {
	return (uint32(payload) & 0x0ffff)
}
func writePacket(packet *Packet) uint32 {
	return startBits(packet.Start) |
		commandBits(packet.Command) |
		addressBits(packet.Address) |
		badgeidBits(packet.BadgeID) |
		payloadBits(packet.Payload)
}
