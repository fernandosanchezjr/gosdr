package buffers

import (
	"sync"
	"sync/atomic"
)

type Stream struct {
	ch        chan []byte
	wg        sync.WaitGroup
	receiving int32
}

func NewStream(size int) *Stream {
	return &Stream{
		ch: make(chan []byte, size-1),
	}
}

func (s *Stream) Send(data []byte) {
	s.ch <- data
}

func (s *Stream) markReceiving() {
	if atomic.CompareAndSwapInt32(&s.receiving, 0, 1) {
		s.wg.Add(1)
	}
}

func (s *Stream) Receive(handler StreamHandler) (closed bool) {
	s.markReceiving()
	select {
	case data, ok := <-s.ch:
		if !ok {
			closed = true
			return
		}
		handler(data)
	}
	return
}

func (s *Stream) Close() {
	close(s.ch)
	s.wg.Wait()
}

func (s *Stream) Done() {
	if atomic.LoadInt32(&s.receiving) == 0 {
		return
	}
	s.wg.Done()
}
