package game

import (
	zl "github.com/rs/zerolog"
)

// Game - The game specification sent to the badge
type Game struct {
	BadgeID   uint16 // ID of badge receiving the game
	AbsStart  int64  // Unix time game starts
	StartTime int16  // The number of seconds from now game starts
	Duration  uint16 // 0x0fff
	Variant   uint8  // 0x0f
	Team      uint8  // 0x0f
	GameID    uint16 // 0x0fff
}

func (g Game) Logger(logger zl.Logger) zl.Logger {
	return logger.With().
		Uint16("BadgeID", g.BadgeID).
		Int64("AbsStart", g.AbsStart).
		Int16("StartTime", g.StartTime).
		Uint16("Duration", g.Duration).
		Uint8("Variant", g.Variant).
		Uint8("Team", g.Team).
		Uint16("GameID", g.GameID).Logger()
}
