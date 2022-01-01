package buffers

import (
	"bytes"
	"container/ring"
)

type BufferRing struct {
	buffers *ring.Ring
}

func NewBufferRing(count int) *BufferRing {
	var r = &BufferRing{
		buffers: ring.New(count),
	}
	r.init()
	return r
}

func (r *BufferRing) init() {
	for i := 0; i < r.buffers.Len(); i++ {
		r.buffers.Value = bytes.NewBuffer(make([]byte, 0))
		r.buffers = r.buffers.Next()
	}
}

func (r *BufferRing) Next() *bytes.Buffer {
	r.buffers = r.buffers.Next()
	var next = r.buffers.Value.(*bytes.Buffer)
	next.Reset()
	return next
}
