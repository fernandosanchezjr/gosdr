package buffers

import "container/ring"

type BlockRing[T BlockType] struct {
	buffers *ring.Ring
}

func NewBlockRing[T BlockType](size, count int) *BlockRing[T] {
	var r = &BlockRing[T]{
		buffers: ring.New(count),
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
