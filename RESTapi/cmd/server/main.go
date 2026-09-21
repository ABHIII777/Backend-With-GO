package main

import (
	"log"
	"net"

	"restapi/internal/repository"
)

func main() {
	// Default store. Postgres replaces this when DATABASE_URL is wired
	// in the next phase; memory keeps `go run ./cmd/server` working now.
	store := repository.NewMemoryStore()

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

		go HandleConnection(conn, store)
	}
}
