package main

import "fmt"

func main() {

	name := "BARNY_FIFE"
	for i, c := range name {
		fmt.Println(i, " => ", string(c))
	}

	var ba [10]byte

	copy(ba[:], name)
	fmt.Println("name:", []byte(name), "ba:", ba)

	bs := make([]byte, 10)

	copy(bs[:], name)

	fmt.Println("name:", []byte(name), "bs:", bs)

	bx := make([]byte, 2)

	bx, bs = bs[0:2], bs[2:]
	fmt.Println("bx:", bx, "bs:", bs)

	bx, bs = bs[0:2], bs[2:]
	fmt.Println("bx:", bx, "bs:", bs)

	bx, bs = bs[0:2], bs[2:]
	fmt.Println("bx:", bx, "bs:", bs)

	bx, bs = bs[0:2], bs[2:]
	fmt.Println("bx:", bx, "bs:", bs)

	bx, bs = bs[0:2], bs[2:]
	fmt.Println("bx:", bx, "bs:", bs)

	var bs3 []byte
	bs3 = append(bs3, 14)

}
