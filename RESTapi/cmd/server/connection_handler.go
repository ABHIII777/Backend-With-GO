package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"restapi/internal/httpwire/request"
	"restapi/internal/httpwire/response"
)

func HandleConnection(conn net.Conn) {
	defer conn.Close()

	if err := conn.SetReadDeadline(time.Now().Add(30 * time.Second)); err != nil {
		log.Printf("Set read deadline: %v", err)
		return
	}

	reader := bufio.NewReader(conn)

	req, err := request.HTTPRequestParser(reader)

	if err != nil {
		var netErr net.Error

		if errors.As(err, &netErr) && netErr.Timeout() {
			log.Printf("Read timeout, closing conn from %s", conn.RemoteAddr())
			return
		}

		log.Printf("Bad request: %v", err)

		_ = conn.SetWriteDeadline(time.Now().Add(5 * time.Second))

		body := "Bad Request\n"
		fmt.Fprintf(conn, "HTTP/1.1 400 Bad Request\r\nContent-Length: %d\r\nContent-Type: text/plain\r\nConnection: close\r\nr\n%s", len(body), body)
		return
	}

	log.Printf("Received: %s %s %s", req.Method, req.Target, req.Version)

	if err := conn.SetWriteDeadline(time.Now().Add(5 * time.Second)); err != nil {
		log.Printf("set write deadline: %v", err)
		return
	}

	switch req.Method {
	case "GET":
		response.HandleGETRequest(conn, req.Target)

	case "POST":
		response.HandlePOSTRequest(conn, req.Target, req.Body)

	default:
		body := fmt.Sprintf("Unknown method %q\n", req.Method)
		fmt.Fprintf(conn, "HTTP/1.1 405 Method Not Allowed\r\nContent-Length: %d\r\nContent-Type: text/plain\r\nConnection: close\r\n\r\n%s", len(body), body)
	}
}
