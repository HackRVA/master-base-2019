package badgewrangler

// EncodeNameBytes - Encode the ascii User Name as a byte array of 5 bit characters
func EncodeNameBytes(name string) []byte {
	bs := make([]byte, 10)
	length := len(name)
	for i, c := range name {
		if c >= 'A' && c <= 'Z' {
			bs[i] = byte(c) - 'A'
		} else if c == '_' {
			bs[i] = byte(26)
		} else {
			bs[i] = byte(27)
		}
	}
	for j := length; j < 10; j++ {
		bs[j] = byte(27)
	}
	return bs
}

// CompressNameBytes - Compress the two element array of character bytes into 5 bit uints
func CompressNameBytes(bytes []byte) uint16 {
	var compressed uint16
	compressed = compressed | (uint16(bytes[0]) << 5)
	compressed = compressed | uint16(bytes[1])
	return compressed
}

// ExpandNameBytes - Expand the 5 bit uints into an 2 element array of regular bytes
func ExpandNameBytes(fragment uint16) []byte {
	bs := make([]byte, 2)
	bs[0] = byte((fragment >> 5) & 0x1f)
	bs[1] = byte(fragment & 0x1f)
	return bs
}

// DecodeNameBytes - Convert the array of 5 bit characters into a trimmed ascii string
func DecodeNameBytes(bytes []byte) string {
	var decoded string
	for _, b := range bytes {
		if b >= 0 && b <= 25 {
			decoded += string(b + 'A')
		} else if b == 26 {
			decoded += "_"
		} else {
			decoded += " "
		}
	}
	return decoded
}
