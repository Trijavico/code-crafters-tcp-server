package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	OK_RESPONSE = "HTTP/1.1 200 OK"
	NOT_FOUND   = "HTTP/1.1 404 Not Found"
)

func handleConn(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	r_size, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading from connection: ", err.Error())
		os.Exit(1)
	}

	request := string(buffer[:r_size])

	lines := strings.Split(request, "\r\n")
	path := strings.Split(lines[0], " ")[1]

	if path == "/" {
		response := OK_RESPONSE + "\r\n\r\n"
		_, err = conn.Write([]byte(response))

	} else if strings.HasPrefix(path, "/echo/") {
		body := path[6:]
		response := fmt.Sprintf("%s\r\nContent-Type: text/plain\r\nContent-Length: %v\r\n\r\n%s\r\n\r\n", OK_RESPONSE, len(body), body)
		_, err = conn.Write([]byte(response))

	} else if path == "/user-agent" {
		body := strings.Split(lines[2], " ")[1]
		response := fmt.Sprintf("%s\r\nContent-Type: text/plain\r\nContent-Length: %v\r\n\r\n%s\r\n\r\n", OK_RESPONSE, len(body), body)
		fmt.Println(response)
		_, err = conn.Write([]byte(response))

	} else {
		response := NOT_FOUND + "\r\n\r\n"
		_, err = conn.Write([]byte(response))

	}

	if err != nil {
		fmt.Println("Failed to write to connection")
		os.Exit(1)
	}
}

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for i := 0; i < 10; i++ {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConn(conn)
	}

}
