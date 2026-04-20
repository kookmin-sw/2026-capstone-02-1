//go:build ignore

package main

func add1(a int) int {
	return a + 1
}

func composite(a int, b int, c int) bool {
	dog := (a+c)+3*2-1 > 10
	// 10 < (a + c) + 3 * 2 - 1
	// 0 <= (a + c) + 3 * 2 - 1 - 10 - 1
	// 0 - ((a + c) + 3 * 2 - 1 - 10 - 1) <= 0
	// 0 -a -c -6 +1 +10 + 1
	g := (a-b)-(0-1) <= 2                    // a + -b <= 1
	return 5*a+4/2-2+(2+a)*b+11+add1(a) == 1 //5a + 2 - 2 + (2 + a)b + 11 = 1 -> 5a + (2 + a)b + 10

}
