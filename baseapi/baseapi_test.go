package baseapi

import (
	"fmt"
	"testing"
	"time"
)

func TestUntil(t *testing.T) {
	gameTime, _ := time.Parse(time.RFC822, "22 May 19 13:30 EST")
	fmt.Println(gameTime)
	currentTime, _ := time.Parse(time.RFC822, "22 May 19 12:33 EST")
	fmt.Println(currentTime)
	secondsUntil := Until(gameTime, currentTime)
	fmt.Println(secondsUntil)

}
