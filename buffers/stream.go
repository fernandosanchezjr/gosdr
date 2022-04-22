package buffers

import (
	"github.com/fernandosanchezjr/gosdr/utils"
	"runtime"
	"sync"
)

type Stream[T BlockType] struct {
	ch        chan *Block[T]
	wg        sync.WaitGroup
	ring      *BlockRing[T]
	receiving uint32
	closed    uint32
	Count     int
	Size      int
}

func NewStream[T BlockType](size int, count int) *Stream[T] {
	return &Stream[T]{
		ch:    make(chan *Block[T], count-1),
		ring:  NewBlockRing[T](size, count),
		Count: count,
		Size:  size,
	}
}

func (s *Stream[T]) markReceiving() {
	if utils.SetFlag(&s.receiving) {
		s.wg.Add(1)
	}
}

func (s *Stream[T]) Send(data *Block[T]) {
	if utils.GetFlag(&s.closed) {
		return
	}
	s.ch <- data
}

func (s *Stream[T]) Receive(handler BlockHandler[T]) (closed bool) {
	if utils.GetFlag(&s.closed) {
		return
	}
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

func (s *Stream[T]) Close() {
	if utils.GetFlag(&s.closed) {
		return
	}
	utils.SetFlag(&s.closed)
	close(s.ch)
	s.wg.Wait()
	runtime.GC()
}

func (s *Stream[T]) Done() {
	if !utils.GetFlag(&s.receiving) {
		return
	}
	s.wg.Done()
	runtime.GC()
}

func (s *Stream[T]) Next() *Block[T] {
	var next = s.ring.Next()
	next.Reset()
	return next
}
