package main

import (
	"bytes"
	"fmt"
	"os"
	"qvarkk/huffman/internal/huffman"
)

type ByteSequence []byte

func main() {
	data, err := os.ReadFile("confession.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	var buf bytes.Buffer
	err = huffman.Code(string(data), &buf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	outputFile, err := os.Create("output.bin")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error. %v\n", err)
		os.Exit(1)
	}
	defer outputFile.Close()

	_, err = buf.WriteTo(outputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to file: %v\n", err)
    os.Exit(1)
	}
}