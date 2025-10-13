package internal

import (
	"fmt"
	"image"
	"unicode/utf8"
)

func DecodeMessage(img image.Image) (string, error) {
	bitMsg := ExtractBitsFromImage(img)
	return BoolSliceToUTF8String(bitMsg)
}

func ExtractBitsFromImage(img image.Image) []bool {
	msgHeader := make([]bool, 0)
	bitMsg := make([]bool, 0)

	count := 0
	var msgLen uint32 = 1
	isHeaderRead := false

	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y && msgLen > 0; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X && msgLen > 0; x++ {
			r, _, _, _ := img.At(x, y).RGBA()
			isBitSet := (r & 1) != 0

			if !isHeaderRead {
				if count < HeaderLen * 8 {
					msgHeader = append(msgHeader, isBitSet)
					count++
				} else {
					x--
					isHeaderRead = true
					msgLen = GetMsgLenFromHeader(msgHeader) * 8
				}
			} else {
				bitMsg = append(bitMsg, isBitSet)
				msgLen--
			}
		}
	}
	return bitMsg
}

func GetMsgLenFromHeader(header []bool) uint32 {
	var value uint32
	for _, bit := range header {
		value <<= 1
		if bit {
			value |= 1
		}
	}
	return value
}

func BoolSliceToUTF8String(slice []bool) (string, error) {
	if len(slice) % 8 != 0 {
		return "", fmt.Errorf("bit slice has to be multiple of 8")
	}

	bytes := make([]byte, len(slice) / 8)
	for i := range len(slice) / 8 {
		var b byte
		for j := range 8 {
			if slice[i*8 + j] {
				b |= (1 << (7 - j))
			}
		}
		bytes[i] = b
	}
	
	if !utf8.Valid(bytes) {
		return "", fmt.Errorf("invalid UTF8 sequence")
	}

	return string(bytes), nil
}