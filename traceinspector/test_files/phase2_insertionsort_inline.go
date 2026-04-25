//go:build ignore

package main

func main() {
	n := 10
	arr := make_array(n, 0)
	for i := 0; i < n; i++ {
		arr[i] = 10 - i
	}

	i := 1

	arr_len := len(arr)
	for i < arr_len {
		key := arr[i]
		// j := i - 1 // correct
		j := i // index error

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

	Print(arr)
}
