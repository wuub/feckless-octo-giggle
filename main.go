package main

import (
	"bufio"
	"os"
)

func main() {
	streams := NewStreamMap()
	go streams.Attach(bufio.NewReader(os.Stdin))
	streams.Serve()
}
