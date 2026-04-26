//go:build ignore

package main

func main() {
	a := 5
	arr := make_array(a, 1337)
	// len_Arr := len(arr)
	for i := 0; i < 10; i++ {
		arr[i] = i
	}
	Print(arr)
}
