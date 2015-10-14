package main

import (
	"io"
	"testing"
)

type MockReadliner struct {
	idx   int
	lines [][]byte
}

func (m *MockReadliner) ReadLine() ([]byte, bool, error) {
	if m.idx >= len(m.lines) {
		return nil, false, io.EOF
	}
	m.idx = m.idx + 1
	return m.lines[m.idx-1], false, nil
}

type MockSender struct {
	lines [][]byte
}

func (m *MockSender) Send(line []byte) error {
	m.lines = append(m.lines, line)
	return nil
}

var testLines [][]byte = [][]byte{[]byte("LOG line1 uuid1 Hello"),
	[]byte("WARN line2 uuid2 hi"),
	[]byte("ERR line3 uuid3 there")}

func TestTrueFilter(t *testing.T) {
	mr := MockReadliner{lines: testLines}
	ms := MockSender{}
	all := func([]byte) bool { return true }
	_ = Pipe(&mr, &ms, all)
	if len(ms.lines) != len(testLines) {
		t.Fatal("expecting two lines")
	}
}

type BenchReadliner struct {
	line  []byte
	count int
}

func (b *BenchReadliner) ReadLine() ([]byte, bool, error) {
	if b.count <= 0 {
		return nil, false, io.EOF
	}
	b.count -= 1
	return b.line, false, nil
}

type BenchSender struct{}

func (b *BenchSender) Send([]byte) error {
	return nil
}

func BenchmarkBenchReadliner(b *testing.B) {
	br := BenchReadliner{line: testLines[0], count: b.N}
	b.SetBytes(int64(len(br.line)))
	for i := 0; i < b.N; i++ {
		br.ReadLine()
	}
}

func BenchmarkPipe(b *testing.B) {
	br := BenchReadliner{line: testLines[0], count: b.N}
	bs := BenchSender{}
	b.SetBytes(int64(len(br.line)))

	Pipe(&br, &bs, func([]byte) bool { return true })
}
