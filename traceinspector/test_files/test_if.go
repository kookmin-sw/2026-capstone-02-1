//go:build ignore

package main

import "fmt"

func main() {
	var x int
	fmt.Scan(&x)
	if x < 0 {
		x = -x + 100
		fmt.Println("got negative")
	} else {
		x = -x
	}
	fmt.Println(x)
}
