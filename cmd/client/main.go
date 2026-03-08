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

	scanner := bufio.NewScanner(os.Stdin)
	temp := make([]byte, 4096)

	lobbyID := handleLobbyFlow(connection, scanner, temp)
	if lobbyID == "" {
		return
	}

	println("What is your name?")
	var name string
	if scanner.Scan() {
		name = scanner.Text()
	}

	go handleMessage(connection)

	msg := message.Message{
		Type: message.MsgText,
		Name: name,
	}
	for scanner.Scan() {
		input := scanner.Text()
		msg.Content = input
		data, _ := json.Marshal(msg)
		connection.Write(data)
	}
}

func handleLobbyFlow(connection *net.TCPConn, scanner *bufio.Scanner, temp []byte) string {
	for {
		n, _ := connection.Read(temp)
		var msg message.Message
		json.Unmarshal(temp[:n], &msg)

		switch msg.Type {
		case message.MsgLobbyPrompt:
			println("Do you want to (c)reate, (j)oin, or (l)ist lobbies?")
			if !scanner.Scan() {
				return ""
			}
			input := scanner.Text()

			respondMsg := message.Message{}
			switch input {
			case "c", "create":
				println("Enter lobby name:")
				if !scanner.Scan() {
					return ""
				}
				lobbyName := scanner.Text()
				respondMsg.Type = message.MsgLobbyCreate
				respondMsg.LobbyID = lobbyName
			case "j", "join":
				println("Enter lobby ID:")
				if !scanner.Scan() {
					return ""
				}
				lobbyName := scanner.Text()
				respondMsg.Type = message.MsgLobbyJoin
				respondMsg.LobbyID = lobbyName
			case "l", "list":
				respondMsg.Type = message.MsgLobbyList
			default:
				println("Invalid choice, try again.")
				continue
			}

			data, _ := json.Marshal(respondMsg)
			connection.Write(data)

		case message.MsgLobbyList:
			var listMsg message.LobbyListMessage
			json.Unmarshal(temp[:n], &listMsg)
			if len(listMsg.Lobbies) == 0 {
				println("No lobbies available. Create one!")
			} else {
				println("Available lobbies:")
				for _, id := range listMsg.Lobbies {
					fmt.Printf("  - %s\n", id)
				}
			}

		case message.MsgError:
			fmt.Printf("Error: %s\n", msg.Content)

		case message.MsgLobbyJoined:
			return msg.LobbyID

		default:
			return msg.LobbyID
		}
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
			var msg message.Message
			json.Unmarshal(temp[:n], &msg)
			fmt.Printf("%s: %s\n", msg.Name, msg.Content)
		}
	}
}
