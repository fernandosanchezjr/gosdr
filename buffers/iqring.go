package buffers

import "container/ring"

type IQRing struct {
	buffers *ring.Ring
}

func NewIQRing(size, count int) *IQRing {
	var r = &IQRing{
		buffers: ring.New(count),
	}
	r.init(size)
	return r
}

func (r *IQRing) init(size int) {
	for i := 0; i < r.buffers.Len(); i++ {
		r.buffers.Value = NewIQ(size)
		r.buffers = r.buffers.Next()
	}
}

func (r *IQRing) Next() *IQ {
	r.buffers = r.buffers.Next()
	return r.buffers.Value.(*IQ)
}
