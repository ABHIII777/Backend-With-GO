package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
) 

type HttpRequestStruct struct {
	Method string;
	Target string;
	Version string;
	Headers map[string] string;
	Body string;
}

func HttpRequest(method string, target string, version string, headers map[string] string, body string) HttpRequestStruct {
	return HttpRequestStruct {
		Method: method,
		Target: target,
		Version: version,
		Headers: headers,
		Body: body,
	}
}

func HttpRequestParser (in *bufio.Reader) (*HttpRequestStruct, error) {
	requestLine, err := in.ReadString('\n');
	if err != nil {
		return nil, fmt.Errorf("read request line: %w", err)
	}

	requestLine = strings.TrimRight(requestLine, "\r\n")

	if requestLine == "" {
		return nil, errors.New("empty request line")
	}

	parts := strings.Fields(requestLine)

	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid request line %q: expected METHOD TARGET VERSION", requestLine)
	}

	method := parts[0]
	target := parts[1]
	version := parts[2]

	if method == "" {
		return nil, errors.New("missing method")
	}

	if !strings.HasPrefix(target, "/") {
		return nil, fmt.Errorf("invalid target %q: must start with /", target)
	}

	if !strings.HasPrefix(version, "HTTP/") {
		return nil, fmt.Errorf("invalid version %q: expected HTTP/x.y", version)
	}

	headers := make(map[string] string)

	for {
		line, err := in.ReadString('\n')

		if err != nil {
			return nil, fmt.Errorf("read header: %w", err)
		}

		line = strings.TrimRight(line, "\r\n")
		if line == "" {
			break
		}

		parts := strings.SplitN(line, ":", 2)

		if len(parts) != 2 {
			return nil, errors.New("Invalid header")
		}

		keys := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		headers[keys] = value
	}

	var body string

	contentLength := -1
	for k, v := range headers {
		if strings.EqualFold(k, "Content-Length") {
			v = strings.TrimSpace(v)
			length, err := strconv.Atoi(v)

			if err != nil || length < 0 {
				return nil, errors.New("invalid Content-Length")
			}

			const maxBody = 1 << 20 
			if length > maxBody {
				return nil, errors.New("body too large")
			}

			contentLength = length
			break
		}
	}

	if contentLength > 0 {
		data := make([]byte, contentLength)
		_, err = io.ReadFull(in, data) 
		if err != nil {
			return nil, err
		}

		body = string(data)
	}

	return &HttpRequestStruct{
		Method: method,
		Target: target,
		Version: version,
		Headers: headers,
		Body: body,
	}, nil

}

func main() {
	ln, err := net.Listen("tcp", ":8080")

	if err != nil {
		log.Fatal(err)
	}

	defer ln.Close()

	for {
		conn, err := ln.Accept()

		if err != nil {
			log.Printf("accept error: %v", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	if err := conn.SetReadDeadline(time.Now().Add(30 * time.Second)); err != nil {
		log.Printf("set read deadline: %v", err)
		return
	}

	reader := bufio.NewReader(conn)

	req, err := HttpRequestParser(reader)
	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			log.Printf("read timeout, closing conn from %s", conn.RemoteAddr())
			return
		}
		log.Printf("bad request: %v", err)
		_ = conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
		body := "Bad Request\n"
		fmt.Fprintf(conn, "HTTP/1.1 400 Bad Request\r\nContent-Length: %d\r\nContent-Type: text/plain\r\nConnection: close\r\n\r\n%s", len(body), body)
		return
	}

	log.Printf("Received: %s %s %s", req.Method, req.Target, req.Version)

	if err := conn.SetWriteDeadline(time.Now().Add(5 * time.Second)); err != nil {
		log.Printf("set write deadline: %v", err)
		return
	}

	switch req.Method {
	case "GET":
		handleGET(conn, req.Target)

	case "POST":
		handlePOST(conn, req.Target, req.Body)

	default:
		body := fmt.Sprintf("Unknown method %q\n", req.Method)
		fmt.Fprintf(conn, "HTTP/1.1 405 Method Not Allowed\r\nContent-Length: %d\r\nContent-Type: text/plain\r\nConnection: close\r\n\r\n%s", len(body), body)
	}

}

func handleGET(conn net.Conn, resource string) {
	body := fmt.Sprintf("GET command received for resource: %s\n", resource)
	fmt.Fprintf(conn, "HTTP/1.1 200 OK\r\nContent-Length: %d\r\nContent-Type: text/plain\r\nConnection: close\r\n\r\n%s", len(body), body)
}

func handlePOST(conn net.Conn, target string, reqBody string) {
	respBody := fmt.Sprintf("POST received for: %s %s\n", target, reqBody)
	fmt.Fprintf(conn, "HTTP/1.1 200 OK\r\nContent-Length: %d\r\nContent-Type: text/plain\r\nConnection: close\r\n\r\n%s", len(respBody), respBody)
}
