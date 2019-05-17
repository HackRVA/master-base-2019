package utility

import "time"

// Int12fromUint16toInt16 - Convert int12 payload to int16
func Int12fromUint16toInt16(x uint16) int16 {
	if (x & 0x0800) > 0 {
		x = x | 0x0f000
	}
	return int16(x)
}

// Int12fromInt16toUint16 - convert int16 to int12 payload
func Int12fromInt16toUint16(x int16) uint16 {
	if x < 0 {
		x = x & 0x0fff
	}
	return uint16(x)
}

// MicroTime - Timestamp in microseconds
func MicroTime() string {
	return time.Now().Format("2006-01-02 15:04:05.000000")
}
