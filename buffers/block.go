package buffers

import "fmt"

type Block[T BlockType] struct {
	Timestamp *Timestamp
	Data      []T
}

func NewBlock[T BlockType](size int) *Block[T] {
	return &Block[T]{
		Timestamp: NewTimestamp(),
		Data:      make([]T, size),
	}
}

func (b *Block[T]) WriteRaw(data []T, ts *Timestamp) {
	ts.Copy(b.Timestamp)
	copy(b.Data, data)
}

func (b *Block[T]) WriteBlock(other *Block[T]) {
	other.Timestamp.Copy(b.Timestamp)
	copy(b.Data, other.Data)
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
