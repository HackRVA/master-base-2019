package badgewrangler

import (
	"fmt"
	"strings"
	"testing"
)

func roundTripName(nameIn string) string {

	bsIn := EncodeNameBytes(nameIn)

	var bsOut []byte

	//var bx []byte

	bx := make([]byte, 2)
	pkts := make([]uint16, 5)

	for i := 0; i < 5; i++ {
		bx, bsIn = bsIn[0:2], bsIn[2:]
		pkts[i] = CompressNameBytes(bx)
	}

	for i := 0; i < 5; i++ {
		bsOut = append(bsOut, ExpandNameBytes(pkts[i])...)
	}

	return strings.TrimSpace(DecodeNameBytes(bsOut))
}

func TestNameEncoding(t *testing.T) {

	name := "BARNY_FIFE"
	fmt.Println("name:", name)

	bs := EncodeNameBytes(name)
	fmt.Println("name:", []byte(name))
	fmt.Println("bs:", bs)

	name2 := "AARON"
	fmt.Println("name2", name2)

	bs2 := EncodeNameBytes(name2)
	fmt.Println("name2:", []byte(name2))
	fmt.Println("bs:", bs2)

	firstTwo := bs[0:2]
	fmt.Println("firstTwo:", firstTwo)

	fmt.Printf("letter mask: %#7x - %016[1]b\n", 0x01f)

	fmt.Println("decoded name1:", DecodeNameBytes(bs), ":")
	fmt.Println("decoded name2:", DecodeNameBytes(bs2), ":")

	if roundTripName(name) != name {
		t.Errorf("roundTripName(name) == name")
	}

	if roundTripName(name2) != name2 {
		t.Errorf("roundTripName(name2) == name2")
	}

}
