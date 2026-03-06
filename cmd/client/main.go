package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/Imlimp/chattd/message"
)

const (
	HOST = "localhost"
	PORT = "5012"
	TYPE = "tcp4"
)

func main() {
	println("What is your name?")
	scanner := bufio.NewScanner(os.Stdin)
	var name string
	if scanner.Scan() {
		input := scanner.Text()
		name = string(input)
	}

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

	message := message.Message{
		Type: message.MsgText,
		Name: name,
	}
	for scanner.Scan() {
		input := scanner.Text()
		message.Content = string(input)
		data, _ := json.Marshal(message)
		connection.Write(data)
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
			var message message.Message
			json.Unmarshal(temp[:n], &message)
			fmt.Printf("%s: %s\n", message.Name, message.Content)
		}
	}
}
