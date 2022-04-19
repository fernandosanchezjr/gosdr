package buffers

type BlockType interface {
	byte | complex64 | complex128 | float32 | float64 | uint16 | uint32 | uint64 | uint | int8 | int16 | int32 | int64 | int | bool
}
