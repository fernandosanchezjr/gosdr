package buffers

import "container/ring"

type ByteRing struct {
	buffers *ring.Ring
}

func NewByteRing(size, count int) *ByteRing {
	var r = &ByteRing{
		buffers: ring.New(count),
	}
	r.init(size)
	return r
}

func (r *ByteRing) init(size int) {
	for i := 0; i < r.buffers.Len(); i++ {
		r.buffers.Value = make([]byte, size)
		r.buffers = r.buffers.Next()
	}
}

func (r *ByteRing) Next() []byte {
	r.buffers = r.buffers.Next()
	return r.buffers.Value.([]byte)
}
