//go:build ignore

package main

import (
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

	_, err = connection.Write([]byte("Wooohooo"))
	if err != nil {
		log.Fatal(err)
	}

	var received []byte
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
			println("Reading data...")
			received = append(received, temp[:n]...)
			println("Received message:", string(received))
		}
	}
}
