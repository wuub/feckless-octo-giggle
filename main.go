package main

import (
	"bufio"
	"log"
	"net"
	"os"
)

func Listen(streams *StreamMap) {
	defer streams.Stop()
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Started on: %s", ln.Addr().String())
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go streams.Handle(conn)
	}
}

func main() {
	streams := NewStreamMap()
	go streams.Attach(bufio.NewReader(os.Stdin))
	go Listen(streams)
	streams.Serve()
}
