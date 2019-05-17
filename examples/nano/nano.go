package main

import (
	"fmt"
	"time"
)

func main() {
	a := makeTimeStamp()
	b := makeTimeStamp2()

	fmt.Println(time.Now().Format("2006-01-02 15:04:05.000000"))
	fmt.Printf("%d \n", time.Now().UnixNano())
	fmt.Printf("%d \n", a)
	fmt.Printf("%d \n", b)

}

func makeTimeStamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func makeTimeStamp2() int64 {
	return time.Now().UnixNano() / 1e6
}
