package response

import (
	"fmt"
	"net"
)

func HandleGETRequest(conn net.Conn, resource string) {
	body := fmt.Sprintf("GET command received for resource: %s\n", resource)
	fmt.Fprintf(conn, "HTTP/1.1 200 OK\r\nContent-Length: %d\r\nContent-Type: text/plain\r\nConnection: close\r\n\r\n%s", len(body), body)
}

func HandlePOSTRequest(conn net.Conn, target string, reqBody string) {
	respBody := fmt.Sprintf("POST received for: %s %s\n", target, reqBody)
	fmt.Fprintf(conn, "HTTP/1.1 200 OK\r\nContent-Length: %d\r\nContent-Type: text/plain\r\nConnection: close\r\n\r\n%s", len(respBody), respBody)
}
