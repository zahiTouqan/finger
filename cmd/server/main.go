package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"stackit.de/zahiTouqan/finger-daemon"
	"strings"
	"time"
)

func main() {
	database := finger.NewUserDatabase()

	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen to TCP port 8080: %v", err)
	}
	defer l.Close()

	log.Println("listening on port 8080")

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %v", err)
			continue
		}

		go handleConn(conn, database)
	}
}

func handleConn(c net.Conn, database *finger.UserDatabase) {
	defer c.Close()

	err := c.SetReadDeadline(time.Now().Add(time.Second * 10))
	if err != nil {
		log.Printf("failed to set read deadline: %v", err)
		return
	}

	// Read query from client and notify it if the read deadline is exceeded
	reader := bufio.NewReader(c)
	query, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("failed to read data from connection: %v", err)

		var opError *net.OpError
		if errors.As(err, &opError) {
			_, err := c.Write([]byte("Read deadline exceeded, ensure that you are including the \\r\\n at the end"))
			if err != nil {
				log.Printf("failed to write error message to connection: %v", err)
				return
			}
		}
		return
	}

	// Parse the arguments
	args, err := parseArgs([]byte(query))
	if err != nil {
		log.Printf("failed to parse args: %v", err)
		return
	}

	// Set verbosity
	verbose := false
	if len(args) != 0 && args[0] == "/W" {
		verbose = true
		args = args[1:]
	}

	// Call the library with the arguments
	if len(args) == 0 { //\r\n
		users := database.GetAllUsers(verbose)
		for _, user := range users {
			err := writeUserToConn(c, &user, verbose)
			if err != nil {
				log.Printf("failed to write user to conn: %v", err)
				return
			}
		}
	}

	// Read queries and forward if needed
	for _, fullWithHost := range args {
		username, host, found := strings.Cut(fullWithHost, "@") // user@somewhere user
		if found {
			fmt.Printf("Forwarding request for user %s to host %s\n", username, host)
			res, err := forwardQuery(username, host, verbose)
			if err != nil {
				log.Printf("failed to forward query: %v", err)
				return
			}

			_, err = c.Write([]byte(res))
			if err != nil {
				log.Printf("failed to write data to connection: %v", err)
				return
			}
			return
		}

		// Unambiguous user lookup
		/*
			fmt.Printf("Looking up user: %s\n", username)
			user, exists := database.GetUser(username)
			if exists {
				err := writeUserToConn(c, &user, verbose)
				if err != nil {
					log.Printf("failed to write user to conn: %v", err)
					return
				}
			} else {
				_, err := c.Write([]byte("Could not find User: " + username + "\r\n"))
				if err != nil {
					log.Printf("failed to write data to connection: %v", err)
					return
				}
			}
		*/

		// Ambiguous user lookup
		users := database.GetUserAmbiguous(username)
		if users == nil {
			log.Printf("user %s not found", username)
			_, err := c.Write([]byte("Could not find User: " + username + "\r\n"))
			if err != nil {
				log.Printf("failed to write data to connection: %v", err)
				return
			}
		}
		for _, user := range users {
			err := writeUserToConn(c, &user, verbose)
			if err != nil {
				log.Printf("failed to write user to conn: %v", err)
				return
			}
		}
	}
}

// helper function to write a specific user to a connection
func writeUserToConn(c net.Conn, user *finger.User, verbose bool) error {
	var err error
	fmt.Printf("Writing data: %s\n", user)
	if verbose {
		_, err = c.Write([]byte(user.String() + "\r\n"))
	} else {
		_, err = c.Write([]byte(user.PartialString() + "\r\n"))
	}
	if err != nil {
		log.Printf("failed to write data to connection: %v", err)
		return err
	}
	return nil
}

func forwardQuery(username string, host string, verbose bool) (string, error) {
	add := host + ":79"
	conn, err := net.DialTimeout("tcp", add, time.Second*10)
	if err != nil {
		log.Printf("failed to connect to server: %v", err)
		return "", err
	}
	defer conn.Close()

	if verbose {
		username = "/W " + username
	}
	username = username + "\r\n"

	_, err = conn.Write([]byte(username + "\r\n"))
	if err != nil {
		log.Printf("failed to write data to connection: %v", err)
		return "", err
	}

	by, err := io.ReadAll(conn)
	if err != nil {
		log.Printf("failed to read data from connection: %v", err)
		return "", err
	}

	return string(by), nil
}

func parseArgs(args []byte) ([]string, error) {
	fmt.Println("Started parsing args")
	var actualArgs []string
	if args == nil {
		log.Printf("Arguments are empty")
		return nil, errors.New("query arguments are empty")
	} else if string(args[len(args)-2:]) != "\r\n" {
		log.Printf("Command must end with <CRLF>")
		return nil, errors.New("command must end with CRLF")
	} else {
		args = args[:len(args)-2]
		lastSpace := -1
		var inArg bool
		for i := 0; i < len(args); i++ {
			if args[i] == ' ' {
				if inArg {
					actualArgs = append(actualArgs, string(args[lastSpace+1:i]))
					inArg = false
				}
				lastSpace = i
			} else {
				inArg = true
				if i == len(args)-1 {
					actualArgs = append(actualArgs, string(args[lastSpace+1:]))
				}
			}
		}
	}
	return actualArgs, nil
}
