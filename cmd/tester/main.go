package main

import (
	"fmt"
)

type A struct {
	A int
}

type B A

func main() {

	C := &B{A: 1}
	fmt.Printf("%#v", C)
}
