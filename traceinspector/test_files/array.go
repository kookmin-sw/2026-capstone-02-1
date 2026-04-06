//go:build ignore

package main

import "fmt"

func triple_array() {
	a := make_array(5, []int{1, 2, 3})
	a[0][1] = 999
	fmt.Print(a, "\n")
}

func resize_array(arr []int, len int, defval int) []int {
	x := make_array(len, defval)
	y := 0
	for i := 0; i < len(arr); i++ {
		x[i] = arr[i]
		y = i
	}
	return x
}

func main() {
	a := []int{111, 222, 333}
	fmt.Print("a:", a, "\n")
	fmt.Print("Enter length to resize a:")
	x := len(a)
	fmt.Scanf("%d", x)
	fmt.Print("Enter fill value:")
	v := 0
	fmt.Scanf("%d", v)
	a = resize_array(a, x, v)
	fmt.Print("resized a:", a, "\n")
	triple_array()
}
