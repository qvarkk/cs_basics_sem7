package internal

import (
	"fmt"
	"sort"
)

func Code(data string) {
	counts := countRunes(data)
	fmt.Printf("%c\n", counts)
}

func countRunes(data string) []rune {
	counts := make(map[rune]int)

	for _, char := range data {
		counts[char]++
	}

	keys := make([]rune, 0, len(counts))
	for key := range counts {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool { return counts[keys[i]] > counts[keys[j]] })	
	
	return keys
}