package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "192.214.173.4:8080")
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("failed to close connection: %v", err)
		}
	}(conn)

	_, err = conn.Write([]byte("/W \r\n"))
	//_, err = conn.Write([]byte("touqan\r\n"))
	//_, err = conn.Write([]byte("tou\r\n"))
	//_, err = conn.Write([]byte("user@4.158.190.163\r\n"))
	//_, err = conn.Write([]byte("user@4.158.190.163"))
	if err != nil {
		log.Fatalf("failed to write: %v", err)
	}

	by, err := io.ReadAll(conn)
	if err != nil {
		log.Fatalf("failed to read data from connection: %v", err)
	}

	fmt.Println(string(by))
}
