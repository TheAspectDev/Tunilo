package protocol

type MsgType byte

const (
	MsgReady MsgType = 1

	MsgPing MsgType = 2
	MsgPong MsgType = 3

	MsgError MsgType = 4

	MsgRequest  MsgType = 6
	MsgResponse MsgType = 7
)

type Message struct {
	Type      MsgType
	RequestID uint64
	Payload   []byte
}
