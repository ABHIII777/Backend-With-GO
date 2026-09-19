package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

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