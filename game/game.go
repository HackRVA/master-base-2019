package game

// Game - The game specification sent to the badge
type Game struct {
	AbsStart  int64  // Unix time game starts
	StartTime int16  // The number of seconds from now game starts
	Duration  uint16 // 0x0fff
	Variant   uint8  // 0x0f
	Team      uint8  // 0x0f
	GameID    uint16 // 0x0fff
}
