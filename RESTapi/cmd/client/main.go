package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

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
			fmt.Println("Try again", err)
			continue
		}

		conn, err := net.DialTimeout("tcp", "localhost:8080", 5*time.Second)

		if err != nil {
			fmt.Println(err)
			continue
		}

		wired := parsed.Encode()

		if err := conn.SetWriteDeadline(time.Now().Add(5 * time.Second)); err != nil {
			fmt.Println("Try Again", err)
			conn.Close()
			continue
		}
		if _, err := fmt.Fprint(conn, wired); err != nil {
			fmt.Println("Try Again", err)
			conn.Close()
			continue
		}

		if err := conn.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
			fmt.Println("Try Again", err)
			conn.Close()
			continue
		}

		bs, err := io.ReadAll(conn)
		conn.Close()
		if err != nil {
			fmt.Println("Try Again", err)
			continue
		}

		fmt.Println(string(bs))
	}
	if err := userInput.Err(); err != nil {
		fmt.Println("input error:", err)
	}
}
