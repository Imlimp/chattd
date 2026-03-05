package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

var (
	connections []net.Conn
	mu          sync.Mutex
)

func main() {
	port := ":5012"
	listner, err := net.Listen("tcp4", port)

	if err != nil {
		log.Fatal(err)
	}
	defer listner.Close()

	fmt.Printf("Opened server on port: %s.\n", port[1:])

	for {
		connection, err := listner.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(connection)
	}
}

func handleConnection(connection net.Conn) {
	mu.Lock()
	connections = append(connections, connection)
	mu.Unlock()
	fmt.Printf("Joined: %s\n", connection.RemoteAddr().String())

	temp := make([]byte, 4096)

	defer connection.Close()
	for {
		_, err := connection.Read(temp)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			mu.Lock()
			for i, c := range connections {
				if c == connection {
					fmt.Printf("Left: %s\n", connection.RemoteAddr().String())
					connections = append(connections[:i], connections[i+1:]...)
					break
				}
			}
			mu.Unlock()
			break
		}
		mu.Lock()
		connectionsCopy := make([]net.Conn, len(connections))
		copy(connectionsCopy, connections)
		mu.Unlock()

		for _, c := range connectionsCopy {
			if c != connection {
				c.Write(temp)
			}
		}
	}
}
