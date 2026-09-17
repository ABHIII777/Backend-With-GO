package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
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
			log.Fatal(err)
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	line, err := reader.ReadString('\n')

	if err != nil {
		log.Fatal(err)
	}

	parts := strings.SplitN(strings.TrimSpace(line), " ", 2)

	if len(parts) != 2 {
		fmt.Println("Invalid format...")
		return
	}

	command := parts[0]
	resource := parts[1]

	log.Printf("Received: %s%s", command, resource)

}
