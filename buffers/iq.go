package buffers

type IQ struct {
	data []complex128
}

func NewIQ(size int) *IQ {
	var buf = &IQ{
		data: make([]complex128, size),
	}
	return buf
}

func convertByte(u byte) float64 {
	return (float64(u) - 127.5) / 128.0
}

func (buf *IQ) Read(raw []byte) (int, error) {
	var read int
	for pos := range buf.data {
		buf.data[pos] = complex(convertByte(raw[read]), convertByte(raw[read+1]))
		read += 2
	}
	return read, nil
}

func (buf *IQ) Data() []complex128 {
	return buf.data
}
