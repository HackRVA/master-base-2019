package irpacket

import "testing"
import "fmt"

const start uint8 = 1
const command uint8 = 1
const address uint8 = 0x1B
const badgeid uint16 = 0x0101
const payload uint16 = 0xffff

func testRawPacket() uint32 {
	return startBits(start) |
		commandBits(command) |
		addressBits(address) |
		badgeidBits(badgeid) |
		payloadBits(payload)
}

func TestBitShifting(t *testing.T) {
	fmt.Println("See how the bits pack")
	fmt.Printf("start   - %#5x - %6[1]d - %08[1]b\n", start)
	fmt.Printf("command - %#5x - %6[1]d - %08[1]b\n", command)
	fmt.Printf("address - %#5x - %6[1]d - %08[1]b\n", address)
	fmt.Printf("badgeid - %#5x - %6[1]d - %016[1]b\n", badgeid)
	fmt.Println()
	fmt.Printf("(start   & 0x01)  >> 32  - %032b - %#[1]x\n", startBits(start))
	fmt.Printf("(command & 0x01)  << 30  - %032b - %#[1]x\n", commandBits(command))
	fmt.Printf("(address & 0x01f) << 25  - %032b - %#[1]x\n", addressBits(address))
	fmt.Printf("(badgeid & 0x1ff) << 16  - %032b - %#[1]x\n", badgeidBits(badgeid))
	fmt.Printf("(payload & 0x0ffff)      - %032b - %#[1]x\n", payloadBits(payload))
	fmt.Println()
	fmt.Printf("bits ored together       - %032b - %#[1]x\n", testRawPacket())
}

func TestReadPacket(t *testing.T) {

	testPacket := readPacket(testRawPacket())

	if testPacket.Start != start {
		t.Errorf("readPacket(testRawPacket()).Start = start")
	}

	if testPacket.Command != command {
		t.Errorf("readPacket(testRawPacket()).Command = command")
	}

	if testPacket.Address != address {
		t.Errorf("readPacket(testRawPacket()).Address = address")
	}

	if testPacket.BadgeID != badgeid {
		t.Errorf("readPacket(testRawPacket()).BadgeID = badgeid")
	}

	if testPacket.Payload != payload {
		t.Errorf("readPacket(testRawPacket()).Payload = payload")
	}

}

func TestBuildPacket(t *testing.T) {

	testPacket := buildPacket(start, command, address, badgeid, payload)

	if testPacket.Start != start {
		t.Errorf("testPacket.Start = start")
	}

	if testPacket.Command != command {
		t.Errorf("testPacket.Command = command")
	}

	if testPacket.Address != address {
		t.Errorf("testPacket.Address = address")
	}

	if testPacket.BadgeID != badgeid {
		t.Errorf("testPacket.BadgeID = badgeid")
	}

	if testPacket.Payload != payload {
		t.Errorf("testPacket.Payload = payload")
	}

}

func TestWritePacket(t *testing.T) {

	testPacket := buildPacket(start, command, address, badgeid, payload)

	if writePacket(testPacket) != testRawPacket() {
		t.Errorf("writePacket(testPacket()) = rawTestPacket")
	}
}
