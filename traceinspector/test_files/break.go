//go:build ignore

package main

import "fmt"

func test_continue() {
	x := 0
	for i := 0; i < 10; i++ {
		if i == 5 {
			continue
		}
		fmt.Print(i)
	}
	fmt.Print("Done", x)
}

func test__nestedcontinue() {
	x := 0
	for i := 0; i < 10; i++ {
		for true {
			if i >= 0 {
				continue
			}
		}
		fmt.Print(i)
	}
	fmt.Print("Done", x)
}

func main() {
	x := 0
	for true {
		x++
		if x > 5 {
			break
		}
	}
	fmt.Print("Done", x)
}
