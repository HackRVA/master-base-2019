package utility

import (
	"fmt"
	"math"
	"testing"
)

func TestInspectProcess(t *testing.T) {
	starttime := int16(-42)
	ustarttime := Int12fromInt16toUint16(starttime)
	custarttime := Int12fromUint16toInt16(ustarttime)
	fmt.Printf(" 42:                            %#7[1]x - %016[1]b\n", uint16(42))
	fmt.Printf("-42:                            %#7[1]x - %016[1]b\n", starttime)
	fmt.Printf("-42 as Int12 payload in Uint16: %#7[1]x - %016[1]b\n", ustarttime)
	fmt.Printf("-42 as Int16 from payload:      %#7[1]x - %016[1]b\n", custarttime)
	fmt.Printf("MaxUint12:             %6d - %#7[1]x - %016[1]b\n", 0x0fff)
	fmt.Printf("MaxInt16:              %6d - %#7[1]x - %016[1]b\n", math.MaxInt16)
}

func TestRoundTrip(t *testing.T) {
	x := int16(-42)
	payloadX := Int12fromInt16toUint16(x)

	y := Int12fromUint16toInt16(payloadX)

	if x != y {
		t.Errorf("round trip of int12 payload does not equal start")
	}
}
