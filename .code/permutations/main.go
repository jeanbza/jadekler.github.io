package main

import "fmt"

var cnt = 0
var cntMemo = 0
var cache = make(map[string]int)

const str = "ABCDEFGH"

func main() {
	fmt.Println("Permuting:", str)
	fmt.Println("Non-memoized answer:", permutations([]rune(str), 0))
	fmt.Println("Non-memoized calls:", cnt)
	fmt.Println("Non-memoized space: 0")
	fmt.Println("Memoized answer:", permutationsMemo([]rune(str), 0))
	fmt.Println("Memoized calls:", cntMemo)
	fmt.Println("Memoized space:", len(cache))
}

func permutations(s []rune, start int) (combinations int) {
	cnt++
	if start == len(s) {
		return 1
	}
	for i := start; i < len(s); i++ {
		s[start], s[i] = s[i], s[start] // swap
		combinations += permutations(s, start+1)
		s[start], s[i] = s[i], s[start] // undo the swap
	}
	return
}

func permutationsMemo(s []rune, start int) (combinations int) {
	if cached, hit := cache[string(s[start:])]; hit {
		return cached
	}

	cntMemo++
	if start == len(s) {
		return 1
	}
	for i := start; i < len(s); i++ {
		s[start], s[i] = s[i], s[start] // swap
		res := permutationsMemo(s, start+1)
		cache[string(s[start:])] = res
		combinations += res
		s[start], s[i] = s[i], s[start] // undo the swap
	}
	return
}
