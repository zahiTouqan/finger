package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
)

func main() {
	serverHost := os.Getenv("SERVER_HOST")
	if serverHost == "" {
		serverHost = "localhost"
	}
	log.Printf("Sever host is %s", serverHost)
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "79"
	}
	log.Printf("Server port is %s", serverPort)
	serverAddress := fmt.Sprintf("%s:%s", serverHost, serverPort)

	wg := new(sync.WaitGroup)

	data := []string{"test /W \r\n", "touqan\r\n", "tou\r\n", "user@4.158.190.163\r\n", "user@4.158.190.163"}
	for i := 0; i < 5; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			makeRequest(data[i], serverAddress)
		}()
	}
	wg.Wait()
}

func makeRequest(data string, serverAddress string) {
	conn, err := net.Dial("tcp4", serverAddress)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()
	log.Printf("Connection is established")

	//_, err = conn.Write([]byte("test /W \r\n"))
	//_, err := conn.Write([]byte("touqan\r\n"))
	//_, err = conn.Write([]byte("tou\r\n"))
	//_, err = conn.Write([]byte("user@4.158.190.163\r\n"))
	//_, err = conn.Write([]byte("user@4.158.190.163"))

	_, err = conn.Write([]byte(data))
	if err != nil {
		log.Fatalf("failed to write: %v", err)
	}

	by, err := io.ReadAll(conn)
	if err != nil {
		log.Fatalf("failed to read data from connection: %v", err)
	}

	fmt.Println(string(by))
}
