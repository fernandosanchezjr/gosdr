package buffers

import (
	"fmt"
	"github.com/fernandosanchezjr/gosdr/utils"
)

type Block[T BlockType] struct {
	Timestamp *Timestamp
	Data      []T
	Size      int
	Pos       int
}

func NewBlock[T BlockType](size int) *Block[T] {
	return &Block[T]{
		Timestamp: NewTimestamp(),
		Data:      make([]T, size),
		Size:      size,
	}
}

func (b *Block[T]) WriteRaw(data []T, ts *Timestamp) {
	ts.Copy(b.Timestamp)
	copy(b.Data, data)
	b.Pos = 0
}

func (b *Block[T]) Write(source *Block[T], ts *Timestamp) int {
	if b.End() || source.End() {
		return 0
	}
	ts.Copy(b.Timestamp)
	var sampleCount = utils.MinInt(b.Remainder(), source.Remainder())
	copy(b.Data[b.Pos:b.Pos+sampleCount], source.Data[source.Pos:source.Pos+sampleCount])
	b.Pos += sampleCount
	source.Pos += sampleCount
	return sampleCount
}

func (b *Block[T]) String() string {
	return fmt.Sprintf("[%d]%T @ %s", len(b.Data), b.Data[0], b.Timestamp)
}

func (b *Block[T]) Less(other *Block[T]) bool {
	return b.Timestamp.Less(other.Timestamp)
}

func (b *Block[T]) CopyTimestamp(ts *Timestamp) *Timestamp {
	if ts == nil {
		ts = b.Timestamp.Child()
	} else {
		b.Timestamp.Copy(ts)
	}
	return ts
}

func (b *Block[T]) End() bool {
	return b.Pos == b.Size
}

func (b *Block[T]) Remainder() int {
	return b.Size - b.Pos
}

func (b *Block[T]) Reset() {
	b.Pos = 0
}
