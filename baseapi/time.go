package baseapi

import (
	"fmt"
	"math"
	"time"
)

// Until -- takes in sceduled gametime and
// returns how many seconds away from the scheduled game time
func Until(gameTime time.Time, fromTime time.Time) uint16 {
	fmt.Println(gameTime)
	return uint16(math.RoundToEven(gameTime.Sub(fromTime).Seconds()))
}
