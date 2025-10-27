package huffman

import (
	"bytes"
	"container/heap"
	"encoding/binary"
	"fmt"
	"qvarkk/huffman/internal/minheap"
	"strconv"
)

func Code(data string, buf *bytes.Buffer) (error) {
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