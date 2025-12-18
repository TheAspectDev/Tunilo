package protocol

type MsgType byte

const (
	MsgReady MsgType = iota
	MsgPing
	MsgPong

	MsgError

	MsgRequest
	MsgResponse
)

type Message struct {
	Type      MsgType
	RequestID uint64
	Payload   []byte
}
