package baseapi

import (
	"fmt"
	"testing"

	bw "github.com/HackRVA/master-base-2019/badgewrangler"
)

var testGameData = &bw.GameData{
	BadgeID: uint16(332),
	GameID:  uint16(1234),
	Hits: []*bw.Hit{
		{BadgeID: uint16(101), Timestamp: uint16(33), Team: uint8(2)},
		{BadgeID: uint16(100), Timestamp: uint16(103), Team: uint8(2)},
		{BadgeID: uint16(101), Timestamp: uint16(203), Team: uint8(2)},
		{BadgeID: uint16(101), Timestamp: uint16(303), Team: uint8(2)},
		{BadgeID: uint16(101), Timestamp: uint16(403), Team: uint8(2)},
		{BadgeID: uint16(101), Timestamp: uint16(503), Team: uint8(2)},
		{BadgeID: uint16(101), Timestamp: uint16(603), Team: uint8(2)},
		{BadgeID: uint16(101), Timestamp: uint16(703), Team: uint8(2)},
	},
}

func TestSaveGameData(t *testing.T) {
	SaveGameData(testGameData)
}

func TestGetGameData(t *testing.T) {
	TestSaveGameData(t)

	var retrieved []string
	for _, c := range GetGameData() {
		fmt.Println(c.ToString())

		retrieved = append(retrieved, c.ToString())
	}
	fmt.Println(retrieved)
}

func TestStrGameData(t *testing.T) {
	fmt.Println(StrGameData())
}

func TestGetNext(t *testing.T) {
	GetNext()
}

func TestZeroGameData(t *testing.T) {
	zeroGameData()
}
