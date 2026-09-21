package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"sort"
	"strings"
	"time"

	"restapi/internal/httpwire/request"
	"restapi/internal/httpwire/response"
	"restapi/internal/repository"
)

func HandleConnection(conn net.Conn, store repository.Store) {
	defer conn.Close()
	_ = store

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

		response.WriteText(conn, 400, "Bad Request\n")
		return
	}

	target := req.Path

	if len(req.Query) > 0 {
		keys := make([]string, 0, len(req.Query))
		for k := range req.Query {
			keys = append(keys, k)
		}

		sort.Strings(keys)

		pairs := make([]string, 0, len(keys))
		for _, k := range keys {
			pairs = append(pairs, k+"="+req.Query[k])
		}

		target += "?" + strings.Join(pairs, "&")
	}

	log.Printf("Received: %s %s %s", req.Method, target, req.Version)

	if err := conn.SetWriteDeadline(time.Now().Add(5 * time.Second)); err != nil {
		log.Printf("set write deadline: %v", err)
		return
	}

	switch req.Method {
	case "GET":
		response.HandleGETRequest(conn, target)

	case "POST", "PUT", "PATCH", "DELETE":
		response.HandlePOSTRequest(conn, target, req.Body)

	default:
		response.WriteText(conn, 405, fmt.Sprintf("Unknown method %q\n", req.Method))
	}
}
