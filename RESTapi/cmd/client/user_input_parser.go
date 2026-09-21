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

	// Preserve raw body exactly (JSON spacing matters inside strings).
	// parts[0]=METHOD, parts[1]=TARGET, remainder of line = BODY.
	body := ""
	if idx := strings.Index(line, target); idx != -1 {
		body = strings.TrimSpace(line[idx+len(target):])
		// If METHOD appears inside body, Index finds TARGET first
		// only when line starts with "METHOD TARGET". parts[1] is the
		// second field, so slicing after its first occurrence is safe
		// for lines shaped "METHOD TARGET [BODY]".
		// Edge: `POST /a POST /b` -> body "POST /b", intended.
		_ = parts
	}

	switch method {
	case "GET", "POST", "PUT", "PATCH", "DELETE":
	default:
		return nil, errors.New("Only GET, POST, PUT, PATCH, DELETE are supported")
	}

	if !strings.HasPrefix(target, "/") {
		return nil, fmt.Errorf("Invalid target %q: must start with /", target)
	}

	if method == "GET" && body != "" {
		return nil, errors.New("GET must not have body")
	}

	headers := map[string]string{
		"Host":       "localhost:8080",
		"Connection": "close",
	}

	if body != "" {
		headers["Content-Length"] = strconv.Itoa(len(body))
		trimmed := strings.TrimSpace(body)
		if strings.HasPrefix(trimmed, "{") {
			headers["Content-Type"] = "application/json"
		} else {
			headers["Content-Type"] = "text/plain"
		}
	}

	if strings.ContainsAny(target, "\r\n") || strings.ContainsAny(body, "\r\n") {
		return nil, errors.New("Target/Body must not contain CR or LF")
	}

	path := target
	query := make(map[string]string)

	if i := strings.Index(target, "?"); i != -1 {
		path = target[:i]
		raw := target[i+1:]
		if raw != "" {
			for _, pair := range strings.Split(raw, "&") {
				if pair == "" {
					continue
				}
				kv := strings.SplitN(pair, "=", 2)
				if len(kv) == 1 {
					query[kv[0]] = ""
					continue
				}
				if kv[0] == "" {
					return nil, errors.New("empty query key")
				}
				query[kv[0]] = kv[1]
			}
		}
	}
	if path == "" {
		path = "/"
	}

	return &UserInputStruct{
		Method:  method,
		Path:    path,
		Query:   query,
		Version: "HTTP/1.1",
		Headers: headers,
		Body:    body,
	}, nil

}
