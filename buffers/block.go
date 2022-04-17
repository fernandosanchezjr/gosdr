package buffers

type Block struct {
	Timestamp *Timestamp
	Data      []byte
}

func NewBlock(size int) *Block {
	return &Block{
		Timestamp: NewTimestamp(),
		Data:      make([]byte, size),
	}
}

func (b *Block) WriteBytes(data []byte, ts *Timestamp) {
	ts.Copy(b.Timestamp)
	copy(b.Data, data)
}

func (b *Block) WriteBlock(other *Block) {
	other.Timestamp.Copy(b.Timestamp)
	copy(b.Data, other.Data)
}
