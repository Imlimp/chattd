package message

type MessageType string

const (
	MsgText        MessageType = "text"
	MsgName        MessageType = "name"
	MsgLobbyCreate MessageType = "create"
	MsgLobbyJoin   MessageType = "join"
	MsgLobbyList   MessageType = "list"
	MsgLobbyPrompt MessageType = "prompt"
	MsgLobbyJoined MessageType = "joined"
	MsgError       MessageType = "error"
)

type Message struct {
	Type    MessageType
	Content string
	Name    string
	LobbyID string
}

type LobbyListMessage struct {
	Type    MessageType
	Lobbies []string
}
