//go:build ignore

package main

import "fmt"

func add1(a int) int {
	return a + 1
}

func composite(a int, b int, c int) bool {
	g := (a-b)-(0-1) <= 2
	return 5*a+4/2-2+(2+a)*b+11+add1(a) == 1 //5a + 2 - 2 + (2 + a)b + 11 = 1 -> 5a + (2 + a)b + 10
}

func main() {
	a := add1(1)
	b := add1(a)
	fmt.Print(b)
}
