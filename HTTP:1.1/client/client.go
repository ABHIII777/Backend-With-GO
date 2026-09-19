package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type UserInputStruct struct {
	Method string
	Target string
	Version string
	Headers map[string] string
	Body string
}

func UserInputParser(line string) (*UserInputStruct, error) {
	line = strings.TrimSpace(line)

	if line == "" {
		return nil, errors.New("Empty input...")
	}

	parts := strings.Fields(line)

	if len(parts) < 2 {
		return nil, fmt.Errorf("Invalid Request line %q: expected METHOD TARGET [BODY]", line)
	}

	method := strings.ToUpper(parts[0])
	target := parts[1]

	body := strings.Join(parts[2:], " ")

	if method != "GET" && method != "POST" {
		return nil, errors.New("Only GET and POST Methods are supported for now")
	}

	if !strings.HasPrefix(target, "/") {
		return nil, fmt.Errorf("Invalid target %q: must start with /", target)
	}

	if method == "GET" && body != "" {
		return nil, errors.New("GET must not have body")
	}

	headers := map[string] string {
		"Host" : "localhost:8080",
		"Connection" : "close",
	}

	if body != "" {
		headers["Content-Length"] = strconv.Itoa(len(body))
		headers["Content-Type"] = "text/plain"
	}

	if strings.ContainsAny(target, "\r\n") || strings.ContainsAny(body, "\r\n") {
		return nil, errors.New("Target/Body must not contain CR or LF")
	}

	return &UserInputStruct{
		Method: method,
		Target: target,
		Version: "HTTP/1.1",
		Headers: headers,
		Body: body,
	}, nil

}

func (input *UserInputStruct) Encode() string {
	var sb strings.Builder
	sb.WriteString(input.Method + " " + input.Target + " " + input.Version + "\r\n")
	sb.WriteString("Host: " + input.Headers["Host"] + "\r\n")
	sb.WriteString("Connection: " + input.Headers["Connection"] + "\r\n")

	if v, ok := input.Headers["Content-Length"]; ok {
		sb.WriteString("Content-Length: " + v + "\r\n")
		sb.WriteString("Content-Type: " + input.Headers["Content-Type"] + "\r\n")
	}

	sb.WriteString("\r\n")
	sb.WriteString(input.Body)
	return sb.String()
}

func main() {
	
	userInput := bufio.NewScanner(os.Stdin)
	
	for {
		fmt.Print("> ")

		if !userInput.Scan() {
			break
		}

		line := userInput.Text()
		if line == "exit" {
			break
		}

		parsed, err := UserInputParser(line)
		if err != nil {
			fmt.Println("Try again", err); continue
		}

		conn, err := net.DialTimeout("tcp", "localhost:8080", 5*time.Second)

		if err != nil {
			fmt.Println(err); continue
		}

		wired := parsed.Encode()

		if err := conn.SetWriteDeadline(time.Now().Add(5 * time.Second)); err != nil {
			fmt.Println("Try Again", err); conn.Close(); continue
		}
		if _, err := fmt.Fprint(conn, wired); err != nil {
			fmt.Println("Try Again", err); conn.Close(); continue
		}

		if err := conn.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
			fmt.Println("Try Again", err); conn.Close(); continue
		}

		bs, err := io.ReadAll(conn)
		conn.Close()
		if err != nil {
			fmt.Println("Try Again", err); continue
		}

		fmt.Println(string(bs))
	}
	if err := userInput.Err(); err != nil {
		fmt.Println("input error:", err)
	}
}
