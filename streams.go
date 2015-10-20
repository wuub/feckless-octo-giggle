package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type Stream struct {
	id        string
	input     chan []byte
	delivered uint64
	dropped   uint64
	processed uint64
	filtered  uint64
	sent      uint64
	filter    Filter
	sender    Sender
	conn      net.Conn
}

func NewStream(conn net.Conn) *Stream {
	s := new(Stream)
	s.input = make(chan []byte, 100)
	s.id = conn.RemoteAddr().String()
	s.conn = conn
	f, err := PrefixFilter("")
	if err != nil {
		panic(err)
	}
	s.filter = f
	return s
}

func (s *Stream) Serve() (err error) {
	var line []byte
	ticker := time.NewTicker(5 * time.Second)
	var now time.Time
	for {
		select {
		case line = <-s.input:
			if s.filter == nil {
				continue
			}
			atomic.AddUint64(&s.processed, 1)
			if !s.filter(line) {
				atomic.AddUint64(&s.filtered, 1)
				continue
			}
		case now = <-ticker.C:
			line = []byte(fmt.Sprintf("%s heartbeat %d %d\n", now, s.delivered, s.dropped))
		}
		// TODO: for now we're ignoring short writes
		_, err = s.conn.Write(line)
		if err != nil {
			return
		}
	}
	return
}

type StreamMap struct {
	stop           chan bool
	processed      uint64
	processedBytes uint64
	input          chan []byte
	streams        map[string]*Stream
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
	s.processed += 1
	s.processedBytes += uint64(len(line))
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
	log.Printf("Connected %s. Active streams: %d\n", stream.id, len(s.streams))
}

func (s *StreamMap) Disconnect(stream *Stream) {
	s.Lock()
	defer s.Unlock()
	delete(s.streams, stream.id)
	log.Printf("Disconnected %s. Active streams: %d\n", stream.id, len(s.streams))
}

func (s *StreamMap) Handle(conn net.Conn) {
	stream := NewStream(conn)
	s.Connect(stream)
	defer s.Disconnect(stream)
	stream.Serve()
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

type ReadByteser interface {
	ReadBytes(byte) ([]byte, error)
}

func (s *StreamMap) Attach(r ReadByteser) {
	defer s.Stop()
	var line []byte
	var err error
	for {
		line, err = r.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		s.input <- line
	}
}
