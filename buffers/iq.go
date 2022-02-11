package buffers

import (
	"errors"
	"github.com/fernandosanchezjr/gosdr/utils"
)

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
	return (float32(u) - 127.4) / 128
}

func (buf *IQ) Read(raw []byte) (int, error) {
	var read int
	for read < len(raw) && buf.pos < buf.size {
		var sample = complex(convertByte(raw[read]), convertByte(raw[read+1]))
		buf.data[buf.pos] = sample
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

func (buf *IQ) Write(out []complex64) (int, error) {
	if buf.Full() {
		return 0, errors.New("end of buffer")
	}
	var copyLen = utils.MinInt(len(out), buf.size-buf.pos)
	copy(out, buf.data[buf.pos:])
	buf.pos += copyLen
	return copyLen, nil
}
