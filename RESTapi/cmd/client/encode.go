package main

import (
	"strings"
)

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