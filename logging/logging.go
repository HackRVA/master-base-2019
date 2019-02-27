package logging

import (
	"os"

	zl "github.com/rs/zerolog"
)

var Ger zl.Logger

func init() {
	file, err := os.OpenFile("test.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err == nil {
		Ger = zl.New(file).With().Timestamp().Logger()
	} else {
		Ger.Debug().Msg("Failed to log to file, using default stderr")
	}
}

func SetGlobalLevel(level zl.Level) {
	zl.SetGlobalLevel(level)
}
