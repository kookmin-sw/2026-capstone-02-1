//go:build ignore

package main

func main() {
	a := 0
	Scanf("%d", a)
	if a <= 0 {
		return 0
	}
	// len_Arr := len(arr)
	arr := make_array(a, 0)
	for i := len(arr) - 1; i >= 0; i-- {
		arr[i] = a - i
	}
	Print(a)
}
