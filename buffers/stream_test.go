package buffers

import (
	log "github.com/sirupsen/logrus"
	"testing"
	"time"
)

const testStreamSize = 64

func createTestBlock() *Block[byte] {
	b := NewBlock[byte](0xffff)
	for pos := range b.Data {
		b.Data[pos] = byte(pos % 0xff)
	}
	return b
}

func testStreamHandler[T BlockType](*Block[T]) {

}

func testAsync(test func() bool) chan bool {
	testChan := make(chan bool)
	resultChan := make(chan bool)
	go func() {
		testChan <- test()
	}()
	go func() {
		select {
		case testResult := <-testChan:
			resultChan <- testResult
		case <-time.After(1 * time.Second):
			resultChan <- false
		}
	}()
	return resultChan
}

func checkResult(t *testing.T, result bool, count *int, args ...any) {
	if !result {
		t.Fatal(args...)
	} else {
		*count += 1
	}
}

func TestNewStream(t *testing.T) {
	s := NewStream[byte](testStreamSize)
	block := createTestBlock()
	for i := 0; i < 1024; i++ {
		s.Send(block)
		if closed := s.Receive(testStreamHandler[byte]); closed {
			t.Fatal("stream closed early")
		}
	}
	closeChan := testAsync(func() bool {
		s.Close()
		return true
	})
	doneChan := testAsync(func() bool {
		if closed := s.Receive(testStreamHandler[byte]); !closed {
			return false
		}
		s.Done()
		return true
	})
	var messages int
	for messages < 2 {
		select {
		case result := <-closeChan:
			checkResult(t, result, &messages, "close failure")
		case result := <-doneChan:
			checkResult(t, result, &messages, "done failure")
		}
	}
}

func BenchmarkStream_Send(b *testing.B) {
	s := NewStream[byte](testStreamSize)
	block := createTestBlock()
	log.Println(block)
	go func() {
		for {
			if closed := s.Receive(testStreamHandler[byte]); closed {
				s.Done()
				return
			}
		}
	}()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Send(block)
	}
	b.StopTimer()
	s.Close()
}

func BenchmarkStream_Receive(b *testing.B) {
	s := NewStream[byte](testStreamSize)
	block := createTestBlock()
	go func() {
		for i := 0; i < b.N; i++ {
			s.Send(block)
		}
	}()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if closed := s.Receive(testStreamHandler[byte]); closed {
			s.Done()
		}
	}
	b.StopTimer()
}
