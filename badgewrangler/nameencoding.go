package badgewrangler

//func encode_username_slice()

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

func CompressNameBytes(bytes []byte) uint16 {
	var compressed uint16
	compressed = compressed | (uint16(bytes[0]) << 5)
	compressed = compressed | uint16(bytes[1])
	return compressed
}

func ExpandNameBytes(fragment uint16) []byte {
	bs := make([]byte, 2)
	bs[0] = byte((fragment >> 5) & 0x1f)
	bs[1] = byte(fragment & 0x1f)
	return bs
}

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
