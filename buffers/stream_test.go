package buffers

import (
	"testing"
	"time"
)

const testStreamSize = 64

func createData() []byte {
	data := make([]byte, 0xffff)
	for i := range data {
		data[i] = byte(i)
	}
	return data
}

func testStreamHandler([]byte) {

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
	s := NewStream(testStreamSize)
	data := createData()
	for i := 0; i < 1024; i++ {
		s.Send(data)
		closed := s.Receive(testStreamHandler)
		if closed {
			t.Fatal("stream closed early")
		}
	}
	closeChan := testAsync(func() bool {
		s.Close()
		return true
	})
	doneChan := testAsync(func() bool {
		closed := s.Receive(testStreamHandler)
		if !closed {
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
	s := NewStream(testStreamSize)
	data := createData()
	go func() {
		for {
			closed := s.Receive(testStreamHandler)
			if closed {
				s.Done()
				return
			}
		}
	}()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Send(data)
	}
	b.StopTimer()
	s.Close()
}
