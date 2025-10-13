package internal

import (
	"encoding/binary"
	"image"
	"image/color"
)

const (
	HeaderLen = 4
)

func EncodeMessage(msg string, img image.Image) image.Image {
	bitMessage := ConvertMessageToBits(msg)

	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rgba.Set(x, y, img.At(x, y))
		}
	}

	bitIndex := 0
	for y := rgba.Bounds().Min.Y; y < rgba.Bounds().Max.Y && bitIndex < len(bitMessage); y++ {
		for x := rgba.Bounds().Min.X; x < rgba.Bounds().Max.X && bitIndex < len(bitMessage); x++ {
			r, g, b, a := rgba.At(x, y).RGBA()

			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)
			a8 := uint8(a >> 8)

			if bitMessage[bitIndex] {
				r8 |= 0x01
			} else {
				r8 &^= 0x01 
			}

			rgba.SetRGBA(x, y, color.RGBA{r8, g8, b8, a8})

			bitIndex++
		}
	}

	return rgba
}

func ConvertMessageToBits(message string) []bool {
	messageBytes := []byte(message)
	length := len(messageBytes)

	header := make([]byte, HeaderLen)
	binary.BigEndian.PutUint32(header, uint32(length))

	payload := append(header, messageBytes...)

	bitMessage := make([]bool, 0, len(payload) * 8)
	for _, b := range payload {
		for i := range 8 {
			bitMessage = append(bitMessage, (b >> (7 - i)) & 1 == 1)
		}
	}

	return bitMessage
}