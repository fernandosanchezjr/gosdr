package buffers

import (
	"testing"
)

const testSampleRate = 2400256
const testByteSize = testSampleRate * 2

func TestNewByteRing(t *testing.T) {
	NewByteRing(testByteSize, 16)
}

func BenchmarkByteRing_Next(b *testing.B) {
	var br = NewByteRing(testByteSize, 16)
	b.ResetTimer()
	var buf []byte
	for i := 0; i < b.N; i++ {
		buf = br.Next()
	}
	if buf == nil {
		b.Fail()
	}
}
