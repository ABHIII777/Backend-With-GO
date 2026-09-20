package main

import (
	"strings"
	"sort"
)

func (input *UserInputStruct) Encode() string {
	var sb strings.Builder

	target := input.Path
	if len(input.Query) > 0 {
		keys := make([]string, 0, len(input.Query))
		for k := range input.Query {
			keys = append(keys, k)
		}
		sort.Strings(keys) 
		pairs := make([]string, 0, len(keys))
		for _, k := range keys {
			pairs = append(pairs, k+"="+input.Query[k])
		}
		target += "?" + strings.Join(pairs, "&")
	}
	sb.WriteString(input.Method + " " + target + " " + input.Version + "\r\n")
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
