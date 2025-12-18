package protocol

import (
	"encoding/binary"
	"io"
)

func Read(r io.Reader) (Message, error) {
	var msg Message

	// Read the first byte
	typeBuf := make([]byte, 1)
	_, err := io.ReadFull(r, typeBuf)
	if err != nil {
		return msg, err
	}

	msg.Type = MsgType(typeBuf[0])

	// read req id ( 8 bytes )
	reqIDBytes := make([]byte, 8)
	_, err = io.ReadFull(r, reqIDBytes)
	if err != nil {
		return msg, err
	}

	msg.RequestID = binary.BigEndian.Uint64(reqIDBytes)

	// Read the next 4 bytes ( payload length )
	lenBuf := make([]byte, 4)
	_, err = io.ReadFull(r, lenBuf)
	if err != nil {
		return msg, err
	}
	length := binary.BigEndian.Uint32(lenBuf)

	// Read the next n=length bytes
	if length > 0 {
		msg.Payload = make([]byte, length)
		_, err = io.ReadFull(r, msg.Payload)
		if err != nil {
			return msg, err
		}
	}

	return msg, nil
}
