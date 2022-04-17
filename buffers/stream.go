package buffers

import (
	"sync"
	"sync/atomic"
)

type Stream struct {
	ch        chan *Block
	wg        sync.WaitGroup
	receiving int32
	closed    int32
}

func NewStream(size int) *Stream {
	return &Stream{
		ch: make(chan *Block, size-1),
	}
}

func (s *Stream) setFlag(flag *int32) bool {
	return atomic.CompareAndSwapInt32(flag, 0, 1)
}

func (s *Stream) getFlag(flag *int32) bool {
	return atomic.LoadInt32(flag) == 1
}

func (s *Stream) markReceiving() {
	if s.setFlag(&s.receiving) {
		s.wg.Add(1)
	}
}

func (s *Stream) Send(data *Block) {
	if s.getFlag(&s.closed) {
		return
	}
	s.ch <- data
}

func (s *Stream) Receive(handler BlockHandler) (closed bool) {
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
	s.setFlag(&s.closed)
	close(s.ch)
	s.wg.Wait()
}

func (s *Stream) Done() {
	if atomic.LoadInt32(&s.receiving) == 0 {
		return
	}
	s.wg.Done()
}
