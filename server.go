package main

import (
	"fmt"
	"io"
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

	var packet []byte
	temp := make([]byte, 4096)

	defer connection.Close()
	for {
		n, err := connection.Read(temp)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			println("END OF FILE")
			break
		}
		packet = append(packet, temp[:n]...)
		num, _ := connection.Write(packet)
		fmt.Printf("Answered back %d, the payload is %s\n", num, string(packet))
	}
}
