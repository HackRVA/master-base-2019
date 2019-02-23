package baseapi

import (
	"math"
	"time"
)

// Until -- takes in sceduled gametime and
// returns how many seconds away from the scheduled game time
func Until(gameTime time.Time) uint16 {
	return uint16(math.RoundToEven(time.Until(gameTime).Seconds()))
}
