package main

import (
	"fmt"
	"log"
	"net"
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
	fmt.Printf("Connection: %s\n", connection.RemoteAddr().String())
}
