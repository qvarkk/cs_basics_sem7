package main

import (
	"fmt"
	"os"
	"qvarkk/huffman/internal"
)

func main() {
	data, err := os.ReadFile("confession.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error. %v\n", err)
		os.Exit(1)
	}

	strData := string(data)

	internal.Code(strData)
}