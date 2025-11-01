package huffman

import (
	"bufio"
	"bytes"
	"container/heap"
	"encoding/binary"
	"fmt"
	"io"
	"qvarkk/huffman/internal/minheap"
	"strconv"
)

func Code(data string, buf *bytes.Buffer) error {
	h := heapifyData(data)
	root := createHuffmanTree(h)
	codes := make(map[rune]string)
	generateHuffmanCodes(&root, "", codes)

	codedData := ""
	for _, char := range string(data) {
		codedData += codes[char]
	}

	if mod := len(codedData) % 8; mod != 0 {
		for range 8 - mod {
			codedData += "0"
		}
	}

	err := writeHuffmanTable(buf, codes)
	if err != nil {
		return err
	}

	writeHuffmanCode(buf, codedData)

	return nil
}

func Decode(reader *bufio.Reader, buf *bytes.Buffer) error {
	codes := make(map[string]rune)
	err := readHuffmanTableCodes(reader, codes)
	if err != nil {
		return fmt.Errorf("couldn't read huffman table: %w", err)
	}

	dataBytes, err := readHuffmanCodes(reader)
	if err != nil {
		return fmt.Errorf("couldn't read huffman table: %w", err)
	}

	err = decodeHuffman(codes, dataBytes, buf)
	if err != nil {
		return fmt.Errorf("couldn't read huffman table: %w", err)
	}

	return nil
}

func heapifyData(data string) *minheap.MinHeap {
	counts := make(map[rune]int)

	for _, char := range data {
		counts[char]++
	}

	h := &minheap.MinHeap{}
	for key, value := range counts {
		h.Push(minheap.HuffmanTreeNode{Key: key, Value: value})
	}

	heap.Init(h)

	return h
}

func createHuffmanTree(h *minheap.MinHeap) minheap.HuffmanTreeNode {
	for h.Len() > 1 {
		left := heap.Pop(h).(minheap.HuffmanTreeNode)
		right := heap.Pop(h).(minheap.HuffmanTreeNode)
		parent := minheap.HuffmanTreeNode{
			Value: left.Value + right.Value,
			Left: &left,
			Right: &right,
		}
		heap.Push(h, parent)
	}
	
	return heap.Pop(h).(minheap.HuffmanTreeNode)
}

func generateHuffmanCodes(node *minheap.HuffmanTreeNode, prefix string, codes map[rune]string) {
	if node == nil {
		return
	}
	if node.Left == nil && node.Right == nil {
		codes[node.Key] = prefix
	}
	generateHuffmanCodes(node.Left, prefix + "0", codes)
	generateHuffmanCodes(node.Right, prefix + "1", codes)
}

func writeHuffmanTable(buf *bytes.Buffer, codes map[rune]string) error {
	if err := binary.Write(buf, binary.BigEndian, uint16(len(codes))); err != nil {
    return fmt.Errorf("couldn't write table size: %v", err)
  }

	for r, code := range codes {
    if err := binary.Write(buf, binary.BigEndian, uint32(r)); err != nil {
      return fmt.Errorf("couldn't write rune: %v", err)
    }

    if err := binary.Write(buf, binary.BigEndian, int8(len(code))); err != nil {
      return fmt.Errorf("couldn't write code length: %v", err)
    }
		if _, err := buf.WriteString(code); err != nil {
			return fmt.Errorf("couldn't write code: %v", err)
		}
	}

	return nil
}

func writeHuffmanCode(buf *bytes.Buffer, data string) error {
	for i := 0; i < len(data); i += 8 {
		code, err := strconv.ParseUint(data[i:i+8], 2, 8)
		if err != nil {
			return fmt.Errorf("couldn't parse code: %v", err)
		}

		if err := binary.Write(buf, binary.BigEndian, byte(code)); err != nil {
    	return fmt.Errorf("couldn't write code byte: %v", err)
    }
	}

	return nil
}

func readHuffmanTableSize(reader *bufio.Reader) (uint16, error) {
	sizeBytes := make([]byte, 0, 2)

	for range 2 {
		b, err := reader.ReadByte()
		if err != nil {
			return 0, fmt.Errorf("couldn't read table size: %w", err)
		}

		sizeBytes = append(sizeBytes, b)
	}

	return binary.BigEndian.Uint16(sizeBytes), nil
}

func readHuffmanTableCodes(reader *bufio.Reader, codes map[string]rune) error {
	size, err := readHuffmanTableSize(reader)
	if err != nil {
		return fmt.Errorf("couldn't read huffman table: %w", err)
	}
	fmt.Printf("Size:\t%d\n\n", size)

	for range size {
		var runeBuf [4]byte
    n, err := reader.Read(runeBuf[:])
    if err != nil || n != 4 {
      return fmt.Errorf("couldn't read 4 bytes for rune: %w", err)
    }
    r := rune(binary.BigEndian.Uint32(runeBuf[:]))

		lenByte, err := reader.ReadByte()
		if err != nil {
			return fmt.Errorf("couldn't read rune in table: %w", err)
		}
		
		code := ""
		for range uint8(lenByte) {
			r, size, err :=  reader.ReadRune()
			if err != nil || size != 1 {
				return fmt.Errorf("couldn't read rune code or rune code is incorrect: %w", err)
			}
			code += string(r)
		}

		codes[code] = r
	}

	return nil
}

func readHuffmanCodes(reader *bufio.Reader) ([]byte, error) {
    var dataBytes []byte
    buf := make([]byte, 1024)

    for {
        n, err := reader.Read(buf)
        if n > 0 {
            dataBytes = append(dataBytes, buf[:n]...)
        }
				if n == 0 && err == nil {
            break
        }
        if err != nil {
            if err == io.EOF {
                break
            }
            return nil, fmt.Errorf("couldn't read compressed file: %w", err)
        }
    }
    return dataBytes, nil
}

func boolSliceToString(slice []bool) string {
	res := ""

	for _, b := range slice {
		if b {
			res += "1"
		} else {
			res += "0"
		}
	}

	return res
}

func decodeHuffman(codes map[string]rune, data []byte, buf *bytes.Buffer) error {
	currCode := make([]bool, 0, 16)

	for _, b := range data {
		for i := range 8 {
			currCode = append(currCode, b&(1<<(7-i)) != 0)
			if r, ok := codes[boolSliceToString(currCode)]; ok {
				buf.WriteRune(r)
				currCode = nil
			}
		}
	}

	return nil
}