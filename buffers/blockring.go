package buffers

import (
	"container/ring"
)

type BlockRing[T BlockType] struct {
	buffers *ring.Ring
	Size    int
	Count   int
}

func NewBlockRing[T BlockType](size, count int) *BlockRing[T] {
	var r = &BlockRing[T]{
		buffers: ring.New(count),
		Size:    size,
		Count:   count,
	}
	r.init(size)
	return r
}

func (r *BlockRing[T]) init(size int) {
	for i := 0; i < r.buffers.Len(); i++ {
		r.buffers.Value = NewBlock[T](size)
		r.buffers = r.buffers.Next()
	}
}

func (r *BlockRing[T]) Next() *Block[T] {
	r.buffers = r.buffers.Next()
	return r.buffers.Value.(*Block[T])
}

func (r *BlockRing[T]) ReverseCopy(destination []*Block[T]) {
	var current = r.buffers
	var block *Block[T]
	for i := len(destination) - 1; i >= 0; i-- {
		block = current.Value.(*Block[T])
		block.Reset()
		destination[i] = block
		current = current.Prev()
	}
}
