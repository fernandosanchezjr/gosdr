package buffers

import "container/ring"

type IQRing struct {
	buffers *ring.Ring
}

func NewIQRing(size, count int) *IQRing {
	var byteRing = &IQRing{
		buffers: ring.New(count),
	}
	byteRing.init(size)
	return byteRing
}

func (br *IQRing) init(size int) {
	for i := 0; i < br.buffers.Len(); i++ {
		br.buffers.Value = NewIQ(size)
		br.buffers = br.buffers.Next()
	}
}

func (br *IQRing) Next() *IQ {
	br.buffers = br.buffers.Next()
	return br.buffers.Value.(*IQ)
}
