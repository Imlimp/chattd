//go:build ignore

package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
)

const (
	HOST = "localhost"
	PORT = "5012"
	TYPE = "tcp4"
)

func main() {
	address := HOST + ":" + PORT
	tcpServer, err := net.ResolveTCPAddr(TYPE, address)
	if err != nil {
		log.Fatal(err)
	}

	connection, err := net.DialTCP(TYPE, nil, tcpServer)
	if err != nil {
		log.Fatal(err)
	}

	defer connection.Close()

	go handleMessage(connection)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		connection.Write([]byte(input))
	}
}

func handleMessage(connection *net.TCPConn) {
	for {
		temp := make([]byte, 4096)
		n, err := connection.Read(temp)
		if err != nil {
			if err == io.EOF {
				break
			}
			println("Failed to receive data.", err.Error())
			os.Exit(1)
		}
		if n > 0 {
			println("Received message:", string(temp[:n]))
		}
	}
}
