package buffers

type BlockHandler[T BlockType] func(block *Block[T])
