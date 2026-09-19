package main

import (
	"log"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")

	if err != nil {
		log.Fatal(err)
	}

	defer ln.Close()

	for {
		conn, err := ln.Accept()

		if err != nil {
			log.Printf("Accept error: %v", err)
			continue
		}

		go HandleConnection(conn)
	}
}
