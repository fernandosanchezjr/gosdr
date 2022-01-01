package buffers

import "container/ring"

type ByteRing struct {
	buffers *ring.Ring
}

func NewByteRing(size, count int) *ByteRing {
	var byteRing = &ByteRing{
		buffers: ring.New(count),
	}
	byteRing.init(size)
	return byteRing
}

func (br *ByteRing) init(size int) {
	for i := 0; i < br.buffers.Len(); i++ {
		br.buffers.Value = make([]byte, size)
		br.buffers = br.buffers.Next()
	}
}

func (br *ByteRing) Next() []byte {
	br.buffers = br.buffers.Next()
	return br.buffers.Value.([]byte)
}
