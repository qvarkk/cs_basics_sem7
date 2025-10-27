package minheap

type HuffmanTreeNode struct {
	Key rune
	Value int
	Left, Right *HuffmanTreeNode
}

type MinHeap []HuffmanTreeNode

func (h MinHeap) Len() int						{ return len(h) }
func (h MinHeap) Less(i, j int) bool	{ return h[i].Value < h[j].Value }
func (h MinHeap) Swap(i, j int)				{ h[i], h[j] = h[j], h[i] }

func (h *MinHeap) Push(x any) {
	*h = append(*h, x.(HuffmanTreeNode))
}

func (h *MinHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}