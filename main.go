package main

import (
	"bufio"
	"fmt"
	// "log"
	"net"
	// "net/http"
	// _ "net/http/pprof"
	"os"
)

func Listen(streams *StreamMap) {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		fmt.Println("new connection")
		go streams.Handle(conn)
	}
}

func main() {
	// go func() {
	// log.Println(http.ListenAndServe(":6060", nil))
	// }()

	streams := NewStreamMap()
	go streams.Attach(bufio.NewReader(os.Stdin))
	go Listen(streams)
	streams.Serve()
}
