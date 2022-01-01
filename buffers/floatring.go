package buffers

import "container/ring"

type FloatRing struct {
	buffers *ring.Ring
}

func NewFloatRing(size, count int) *FloatRing {
	var r = &FloatRing{
		buffers: ring.New(count),
	}
	r.init(size)
	return r
}

func (r *FloatRing) init(size int) {
	for i := 0; i < r.buffers.Len(); i++ {
		r.buffers.Value = make([]float32, size)
		r.buffers = r.buffers.Next()
	}
}

func (r *FloatRing) Next() []float32 {
	r.buffers = r.buffers.Next()
	return r.buffers.Value.([]float32)
}
