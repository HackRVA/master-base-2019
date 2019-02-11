package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	l, err := net.ListenUnix("unix", &net.UnixAddr{"/tmp/unixdomain", "unix"})
	if err != nil {
		panic(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			if sig == syscall.SIGINT {
				fmt.Println(sig)
				os.Remove("/tmp/unixdomain")
				os.Exit(0)
			}
		}
	}()

	for {
		conn, err := l.AcceptUnix()
		if err != nil {
			panic(err)
		}
		var buf [1024]byte
		n, err := conn.Read(buf[:])
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", string(buf[:n]))
		conn.Close()
	}
}
