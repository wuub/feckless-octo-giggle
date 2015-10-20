package main

import (
	"fmt"
	"io"
	"sync"
	"sync/atomic"
)

type Stream struct {
	id        string
	input     chan []byte
	delivered uint64
	dropped   uint64
	processed uint64
	filtered  uint64
	sent      uint64
}

type StreamMap struct {
	stop    chan bool
	input   chan []byte
	streams map[string]*Stream
	sync.RWMutex
}

func NewStreamMap() *StreamMap {
	sm := new(StreamMap)
	sm.input = make(chan []byte)
	sm.stop = make(chan bool)
	sm.streams = make(map[string]*Stream)
	return sm
}

func (s *StreamMap) Deliver(line []byte) {
	s.RLock()
	defer s.RUnlock()
	fmt.Println(line)
	for _, stream := range s.streams {
		select {
		case stream.input <- line:
			atomic.AddUint64(&stream.delivered, 1)
		default:
			atomic.AddUint64(&stream.dropped, 1)
		}
	}
}

func (s *StreamMap) Connect(stream *Stream) {
	s.Lock()
	defer s.Unlock()
	s.streams[stream.id] = stream
}

func (s *StreamMap) Disconnect(stream *Stream) {
	s.Lock()
	defer s.Unlock()
	delete(s.streams, stream.id)
}

func (s *StreamMap) Serve() {
	var line []byte
	for {
		select {
		case line = <-s.input:
			s.Deliver(line)
		case <-s.stop:
			break
		}
	}
}

func (s *StreamMap) Stop() {
	s.stop <- true
}

func (s *StreamMap) Attach(r Readliner) {
	defer s.Stop()
	for {
		line, _, err := r.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		s.input <- line
	}
}
