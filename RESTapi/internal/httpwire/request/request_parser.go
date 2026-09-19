package request

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func HTTPRequestParser(input *bufio.Reader) (*HTTPRequestStruct, error) {
	requestLine, err := input.ReadString('\n')

	if err != nil {
		return nil, fmt.Errorf("Read request line: %w", err)
	}

	requestLine = strings.TrimRight(requestLine, "\r\n")

	if requestLine == "" {
		return nil, errors.New("Empty request line")
	}

	parts := strings.Fields(requestLine)

	if len(parts) != 3 {
		return nil, fmt.Errorf("Invalid Request line %q: expected METHOD TARGET VERSION", requestLine)
	}

	method := parts[0]
	target := parts[1]
	version := parts[2]

	if method == "" {
		return nil, errors.New("Missing Method")
	}

	if !strings.HasPrefix(target, "/") {
		return nil, fmt.Errorf("Invalid target %q: must start with /", target)
	}

	if !strings.HasPrefix(version, "HTTP") {
		return nil, fmt.Errorf("Invalid version %q: expected HTTP/x.y", version)
	}

	headers := make(map[string]string)

	for {
		line, err := input.ReadString('\n')

		if err != nil {
			return nil, fmt.Errorf("Read header: %w", err)
		}

		line = strings.TrimRight(line, "\r\n")

		if line == "" {
			break
		}

		parts := strings.SplitN(line, ":", 2)

		if len(parts) != 2 {
			return nil, errors.New("Invalid header")
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		headers[key] = value
	}

	var body string

	contentLength := -1

	for key, value := range headers {
		if strings.EqualFold(key, "Content-Length") {
			value = strings.TrimSpace(value)
			length, err := strconv.Atoi(value)

			if err != nil || length < 0 {
				return nil, errors.New("Invalid Content-Length")
			}

			const maxBody = 1 << 20

			if length > maxBody {
				return nil, errors.New("Body length is too large")
			}

			contentLength = length
			break
		}
	}

	if contentLength > 0 {
		data := make([]byte, contentLength)

		_, err := io.ReadFull(input, data)
		if err != nil {
			return nil, err
		}

		body = string(data)
	}

	return &HTTPRequestStruct{
		Method:  method,
		Target:  target,
		Version: version,
		Headers: headers,
		Body:    body,
	}, nil

}
