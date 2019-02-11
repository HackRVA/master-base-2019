package main

import (
	"net"
	"os"
)

func main() {
	connType := "unix" // or "unixgram" or "unixpacket"
	laddr := net.UnixAddr{"/tmp/unixdomaincli", connType}
	conn, err := net.DialUnix(connType, &laddr /*can be nil*/, &net.UnixAddr{"/tmp/unixdomain", connType})
	if err != nil {
		panic(err)
	}
	defer os.Remove("/tmp/unixdomaincli")

	_, err = conn.Write([]byte("hello"))
	if err != nil {
		panic(err)
	}
	conn.Close()
}
