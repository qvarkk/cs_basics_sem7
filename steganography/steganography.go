package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"image/png"
	"os"
	"strings"
)

func main() {
  fmt.Print("Enter file path: ")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
    os.Exit(1)
	}

	filePath := strings.TrimSpace(input)
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot open file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	_, err = png.Decode(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "File is not a valid PNG: %v\n", err)
		os.Exit(1)
	}

	binData, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read file content: %v\n", err)
		os.Exit(1)
	}

	doTheThing(binData, "test")
}

func doTheThing(binData []byte, message string) {
	binData = binData[8:] // skip png header

	for range 2 {
		chunkLength := binary.BigEndian.Uint32(binData[:4])
		chunkType := binData[4:8]

		fmt.Printf("% d\n", chunkLength)
		fmt.Printf("% x\n", chunkType)

		if bytes.Equal(chunkType, []byte {0x49, 0x48, 0x44, 0x52}) {
			fmt.Println("IHDR chunk, skipping")
		}
	}
}