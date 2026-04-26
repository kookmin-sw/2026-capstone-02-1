//go:build ignore

package main

import "fmt"

func add1(a int) int {
	return a + 1
}

func main() {
	a := add1(1)
	b := add1(a)
	// c := add1(add1(add1(add1(add1(add1(add1(1)))))))
	c := add1(add1(1))
	fmt.Print(c)
}
