package message

type MessageType string

const (
	MsgText MessageType = "text"
	MsgName MessageType = "name"
)

type Message struct {
	Type    MessageType
	Content string
	Name    string
}
