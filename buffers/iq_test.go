package buffers

import (
	"testing"
)

func BenchmarkIQ_Read(b *testing.B) {
	var iq = NewIQ(testSampleRate / 2)
	var rawBuf = make([]byte, testSampleRate)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := iq.Read(rawBuf); err != nil {
			b.Fail()
		}
	}
}
