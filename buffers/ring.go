package buffers

import "container/ring"

type ByteRing struct {
	buffers *ring.Ring
}

func NewByteRing(byteSize, bufferCount int) *ByteRing {
	var byteRing = &ByteRing{
		buffers: ring.New(bufferCount),
	}
	byteRing.init(byteSize)
	return byteRing
}

func (br *ByteRing) init(byteSize int) {
	for i := 0; i < br.buffers.Len(); i++ {
		br.buffers.Value = make([]byte, byteSize)
		br.buffers = br.buffers.Next()
	}
}

func (br *ByteRing) Next() []byte {
	var nextBuffer = br.buffers.Value
	br.buffers = br.buffers.Next()
	return nextBuffer.([]byte)
}
