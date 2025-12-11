package protocol

type Type byte

const (
	MsgReady Type = 1

	MsgPing Type = 2
	MsgPong Type = 3

	MsgError Type = 4

	MsgRequest  Type = 6
	MsgResponse Type = 7
)

type Message struct {
	Type      Type
	RequestID uint64
	Payload   []byte
}
