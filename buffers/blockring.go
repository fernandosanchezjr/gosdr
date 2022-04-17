package buffers

import "container/ring"

type BlockRing struct {
	buffers *ring.Ring
}

func NewBlockRing(size, count int) *BlockRing {
	var r = &BlockRing{
		buffers: ring.New(count),
	}
	r.init(size)
	return r
}

func (r *BlockRing) init(size int) {
	for i := 0; i < r.buffers.Len(); i++ {
		r.buffers.Value = NewBlock(size)
		r.buffers = r.buffers.Next()
	}
}

func (r *BlockRing) Next() *Block {
	r.buffers = r.buffers.Next()
	return r.buffers.Value.(*Block)
}
