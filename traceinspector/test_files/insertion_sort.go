//go:build ignore

package main

import "fmt"

func insertionSort(arr []int) []int {
	i := 1

	arr_len := len(arr)
	for i < arr_len {
		key := arr[i]
		j := i - 1

		for j >= 0 {
			cur := arr[j]
			if cur > key {
				arr[j+1] = arr[j]
				j = j - 1
				continue
			}

			break
		}

		arr[j+1] = key
		i = i + 1
	}

	return arr
}

func main() {
	n := 0
	fmt.Scanf("%d", n)
	arr := make_array(n, 0)
	for i := 0; i < n; i++ {
		x := 0
		fmt.Scanf("%d", x)
		arr[i] = x
	}
	fmt.Print(insertionSort(arr), "\n")
}
