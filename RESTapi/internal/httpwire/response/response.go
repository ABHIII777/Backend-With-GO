package response

import (
	"encoding/json"
	"fmt"
	"net"
)

func statusText(status int) string {
	switch status {
	case 200:
		return "OK"
	case 201:
		return "Created"
	case 204:
		return "No Content"
	case 400:
		return "Bad Request"
	case 404:
		return "Not Found"
	case 405:
		return "Method Not Allowed"
	case 409:
		return "Conflict"
	case 500:
		return "Internal Server Error"
	default:
		return "OK"
	}
}

func writeHead(conn net.Conn, status int, contentType string, contentLength int) {
	fmt.Fprintf(conn, "HTTP/1.1 %d %s\r\nContent-Length: %d\r\nContent-Type: %s\r\nConnection: close\r\n\r\n",
		status, statusText(status), contentLength, contentType)
}

func WriteJSON(conn net.Conn, status int, v any) {
	if status == 204 {
		writeHead(conn, status, "application/json", 0)
		return
	}
	bs, err := json.Marshal(v)
	if err != nil {
		WriteError(conn, 500, "encode error")
		return
	}
	bs = append(bs, '\n')
	writeHead(conn, status, "application/json", len(bs))
	_, _ = conn.Write(bs)
}

func WriteError(conn net.Conn, status int, msg string) {
	WriteJSON(conn, status, map[string]string{"error": msg})
}

func WriteText(conn net.Conn, status int, body string) {
	writeHead(conn, status, "text/plain", len(body))
	_, _ = conn.Write([]byte(body))
}
