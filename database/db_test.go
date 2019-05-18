package database

import (
	"strings"
	"testing"
	"time"

	bw "github.com/HackRVA/master-base-2019/badgewrangler"
	gm "github.com/HackRVA/master-base-2019/game"
)

var testGameData = &bw.GameData{
	BadgeID: uint16(332),
	GameID:  uint16(1234),
	Hits: []*bw.Hit{
		{GameID: uint16(101), BadgeID: uint16(101), Timestamp: uint16(33), Team: uint8(2)},
		{GameID: uint16(101), BadgeID: uint16(100), Timestamp: uint16(103), Team: uint8(2)},
		{GameID: uint16(101), BadgeID: uint16(101), Timestamp: uint16(203), Team: uint8(2)},
		{GameID: uint16(101), BadgeID: uint16(101), Timestamp: uint16(303), Team: uint8(2)},
		{GameID: uint16(101), BadgeID: uint16(101), Timestamp: uint16(403), Team: uint8(2)},
		{GameID: uint16(101), BadgeID: uint16(101), Timestamp: uint16(503), Team: uint8(2)},
		{GameID: uint16(101), BadgeID: uint16(101), Timestamp: uint16(603), Team: uint8(2)},
		{GameID: uint16(101), BadgeID: uint16(101), Timestamp: uint16(703), Team: uint8(2)},
	},
}

// TestSaveGameData -- Tests SaveGameData and GetGameData
func TestSaveGameData(t *testing.T) {
	SaveGameData(testGameData)
	gData := GetGameData()

	if gData[0].BadgeID != testGameData.BadgeID {
		t.Error("BadgeID does not match test data")
	}
	if gData[0].GameID != testGameData.GameID {
		t.Error("GameID does not match test data")
	}

	for i, h := range gData[0].Hits {
		if h.BadgeID != testGameData.Hits[i].BadgeID {
			t.Error("BadgeId - hit data doesn't match testData")
		}
		if h.Timestamp != testGameData.Hits[i].Timestamp {
			t.Error("Timestamp - hit data doesn't match testData")
		}
		if h.Team != testGameData.Hits[i].Team {
			t.Error("Team - hit data doesn't match testData")
		}
	}
}

func TestStrGameData(t *testing.T) {
	d := &GameDataWithSent{}
	d.GameData = *testGameData

	var s []string

	s = append(s, d.ToString())

	if strings.Compare(StrGameData()[0], s[0]) != 0 {
		t.Error("the strings don't match -- StrGameData isn't working correctly")
	}
}

// TestSchedule -- tests scheduling and GetNext
func TestSchedule(t *testing.T) {
	twoMin := time.Now().Local().Add(time.Minute * time.Duration(2))
	testGame := &gm.Game{
		BadgeID:  555,
		Team:     55,
		Duration: 12,
		Variant:  1,
		AbsStart: twoMin.Unix(),
	}

	ScheduleGame(*testGame)

	next := GetNext()
	another := GetNext()

	if next.BadgeID != testGame.BadgeID {
		t.Error("Test BadgeID does not match with next scheduled game")
	}
	if next.Team != 2 {
		t.Error("Team assignment is not working correctly")
	}
	if another.Team != 1 {
		t.Error("Team assignment is not working correctly")
	}
	if next.Duration != testGame.Duration {
		t.Error("Test Duration does not match with next scheduled game")
	}
	if next.Variant != testGame.Variant {
		t.Error("Test Variant does not match with next scheduled game")
	}
}

// TestZombie -- testing that zombie is created and then it should switch to human
func TestZombie(t *testing.T) {
	if determineTeam(2, 2) != 1 {
		t.Error("determine team should send PatientZero")
	}
	if determineTeam(2, 2) != 2 {
		t.Error("determine team should start sending humans")
	}
	if determineTeam(2, 2) != 2 {
		t.Error("determine team should still be sending humans")
	}
	if determineTeam(2, 3) != 1 {
		t.Error("determine team should send PatientZero")
	}
	if determineTeam(2, 3) != 2 {
		t.Error("determine team should start sending humans")
	}
	if determineTeam(2, 3) != 2 {
		t.Error("determine team should still be sending humans")
	}
}
