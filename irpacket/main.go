package irpacket

// #include <badge-ir-game-protocol.h>
import "C"
import "fmt"

func opcodes() {
	fmt.Println("BADGE_IR_GAME_ADDRESS", C.BADGE_IR_GAME_ADDRESS)
}
