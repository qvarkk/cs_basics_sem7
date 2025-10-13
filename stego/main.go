package main

import (
	"bufio"
	"fmt"
	"image"
	"image/png"
	"os"
	"strconv"
	"strings"

	"stego/internal"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Select an option:")
	fmt.Println("0 - Embed message")
	fmt.Println("1 - Extract message")

	optionStr, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't read user input: %v\n", err)
		os.Exit(1)
	}
	optionStr = strings.TrimSpace(optionStr)

	option, err := strconv.Atoi(optionStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error. %v\n", err)
		os.Exit(1)
	}

	switch option {
	case 0:
		handleEmbed(reader)
	case 1:
		handleExtract(reader)
	default:
		fmt.Fprintf(os.Stderr, "Error. No such option: %v\n", err)
		os.Exit(1)
	}
}

func handleExtract(reader *bufio.Reader) {
	img, err := ReadImageByPath(reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't read PNG image: %v\n", err)
		os.Exit(1)
	}

	msg, err := internal.DecodeMessage(img)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't decode message: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Message is: %s\n", msg)
}

func handleEmbed(reader *bufio.Reader) {
	msg, err := ReadMessage(reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't read message: %v\n", err)
		os.Exit(1)
	}
  
	img, err := ReadImageByPath(reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't read PNG image: %v\n", err)
		os.Exit(1)
	}

	encodedImg := internal.EncodeMessage(msg, img)
	err = CreatePngImage(encodedImg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't create PNG image: %v\n", err)
		os.Exit(1)
	}
}

func ReadMessage(reader *bufio.Reader) (string, error) {
	fmt.Print("Enter message to encode: ")
	msg, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(msg), nil
}

func ReadImageByPath(reader *bufio.Reader) (image.Image, error) {
	fmt.Print("Enter file path: ")
	fileInput, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	filePath := strings.TrimSpace(fileInput)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func CreatePngImage(img image.Image) error {
	outputFile, err := os.Create("output.png")
	if err != nil {
		return err
	}
	defer outputFile.Close()

	err = png.Encode(outputFile, img)
	if err != nil {
		return err
	}

	return nil
}

func BoolSliceToBinaryString(slice []bool) string {
	var sb strings.Builder
	sb.Grow(len(slice))

	for _, b := range slice {
		if b {
			sb.WriteByte('1')
		} else {
			sb.WriteByte('0')
		}
	}

	return sb.String()
}