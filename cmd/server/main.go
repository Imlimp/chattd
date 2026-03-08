package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	"github.com/Imlimp/chattd/message"
)

type Lobby struct {
	ID      string
	Name    string
	Members map[net.Conn]bool
}

var (
	lobbies     = make(map[string]*Lobby)
	lobbyMu     sync.Mutex
	connections = make(map[net.Conn]string)
	connMu      sync.Mutex
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
	fmt.Printf("Joined: %s sending lobby prompt\n", connection.RemoteAddr().String())

	lobbyPromptMessage := message.Message{
		Type: message.MsgLobbyPrompt,
	}
	data, _ := json.Marshal(lobbyPromptMessage)
	connection.Write(data)

	temp := make([]byte, 4096)
	n, err := connection.Read(temp)
	if err != nil {
		if err == io.EOF {
			return
		}
		println("Failed to receive data.", err.Error())
		return
	}
	if n > 0 {
		var msg message.Message
		json.Unmarshal(temp[:n], &msg)

		switch msg.Type {
		case message.MsgLobbyJoin:
			lobbyID := msg.LobbyID
			lobbyMu.Lock()
			lobby, exists := lobbies[lobbyID]
			lobbyMu.Unlock()
			if !exists {
				errorMsg := message.Message{
					Type:    message.MsgError,
					Content: "Lobby not found",
				}
				data, _ := json.Marshal(errorMsg)
				connection.Write(data)
				return
			}
			lobby.Members[connection] = true
			connMu.Lock()
			connections[connection] = lobbyID
			connMu.Unlock()
			joinedMsg := message.Message{
				Type:    message.MsgLobbyJoined,
				LobbyID: lobbyID,
			}
			data, _ := json.Marshal(joinedMsg)
			connection.Write(data)
			handleChat(connection, lobbyID)
		case message.MsgLobbyCreate:
			lobbyID := msg.LobbyID
			newLobby := &Lobby{
				ID:      lobbyID,
				Name:    lobbyID,
				Members: make(map[net.Conn]bool),
			}
			newLobby.Members[connection] = true
			lobbyMu.Lock()
			lobbies[lobbyID] = newLobby
			lobbyMu.Unlock()
			connMu.Lock()
			connections[connection] = lobbyID
			connMu.Unlock()
			joinedMsg := message.Message{
				Type:    message.MsgLobbyJoined,
				LobbyID: lobbyID,
			}
			data, _ := json.Marshal(joinedMsg)
			connection.Write(data)
			handleChat(connection, lobbyID)
		case message.MsgLobbyList:
			lobbyMu.Lock()
			var lobbyIDs []string
			for id := range lobbies {
				lobbyIDs = append(lobbyIDs, id)
			}
			lobbyMu.Unlock()
			listMsg := message.LobbyListMessage{
				Type:    message.MsgLobbyList,
				Lobbies: lobbyIDs,
			}
			data, _ := json.Marshal(listMsg)
			connection.Write(data)
			handleConnection(connection)
		}
	}

	connection.Close()
}

func handleChat(connection net.Conn, lobbyID string) {
	temp := make([]byte, 4096)
	defer func() {
		lobbyMu.Lock()
		lobby := lobbies[lobbyID]
		if lobby != nil {
			delete(lobby.Members, connection)
			if len(lobby.Members) == 0 {
				delete(lobbies, lobbyID)
			}
		}
		lobbyMu.Unlock()
		connMu.Lock()
		delete(connections, connection)
		connMu.Unlock()
		fmt.Printf("Left: %s\n", connection.RemoteAddr().String())
		connection.Close()
	}()

	for {
		n, err := connection.Read(temp)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			break
		}

		var msg message.Message
		json.Unmarshal(temp[:n], &msg)
		msg.LobbyID = lobbyID

		broadcastToLobby(lobbyID, msg, connection)
	}
}

func broadcastToLobby(lobbyID string, msg message.Message, exclude net.Conn) {
	lobbyMu.Lock()
	lobby := lobbies[lobbyID]
	lobbyMu.Unlock()

	if lobby == nil {
		return
	}

	data, _ := json.Marshal(msg)
	for conn := range lobby.Members {
		if conn != exclude {
			conn.Write(data)
		}
	}
}
