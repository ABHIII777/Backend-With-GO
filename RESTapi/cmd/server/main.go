package main

import (
	"log"
	"net"
	"os"

	"restapi/internal/repository"
)

func main() {
	// Load ./RESTapi/.env when present so `go run ./cmd/server` works
	// from any shell (bash, fish, zsh). Real environment variables
	// take precedence over file values.
	loadDotEnv(".env")

	// Postgres is the real backend when DATABASE_URL is set.
	var store repository.Store = repository.NewMemoryStore()

	if url := os.Getenv("DATABASE_URL"); url != "" {
		db, err := repository.NewPostgresDB(url)
		if err != nil {
			log.Fatalf("Postgres unreachable: %v", err)
		}
		defer db.Close()

		store = repository.NewPostgresStore(db)
		log.Println("Store: postgres")
	} else {
		log.Println("Store: memory (DATABASE_URL unset)")
	}

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
