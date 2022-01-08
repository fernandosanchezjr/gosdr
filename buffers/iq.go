package buffers

type IQ struct {
	Sequence uint64
	size     int
	data     []complex64
	pos      int
}

func NewIQ(size int) *IQ {
	var buf = &IQ{
		size: size,
		data: make([]complex64, size),
	}
	return buf
}

func convertByte(u byte) float32 {
	return (float32(u) - 127.5) / 128.0
}

func (buf *IQ) Read(raw []byte) (int, error) {
	var read int
	for read < len(raw) && buf.pos < buf.size {
		buf.data[buf.pos] = complex(convertByte(raw[read]), convertByte(raw[read+1]))
		read += 2
		buf.pos += 1
	}
	return read, nil
}

func (buf *IQ) Data() []complex64 {
	return buf.data
}

func (buf *IQ) Copy(source *IQ) {
	buf.Sequence = source.Sequence
	copy(buf.data, source.data)
}

func (buf *IQ) Full() bool {
	return buf.size == buf.pos
}

func (buf *IQ) Reset() {
	buf.pos = 0
}
