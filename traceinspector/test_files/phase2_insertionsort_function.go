//go:build ignore

package main

func insertionSort(arr []int) []int {
	i := 1

	for i < len(arr) {
		key := arr[i]
		j := i - 1
		// j := i

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
	n := 10
	arr := make_array(n, 0)
	for i := 0; i < n; i++ {
		arr[i] = 10 - i
	}

	Print(insertionSort(arr))
}
