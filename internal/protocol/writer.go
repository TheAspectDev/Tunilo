package protocol

import (
	"encoding/binary"
	"io"
)

// FULLY WRITTEN MESSAGE
// [MESSAGE_TYPE REQUEST_ID MESSAGE_LENGTH PAYLOAD]
func Write(w io.Writer, msg Message) error {
	// Write TYPE
	_, err := w.Write([]byte{byte(msg.Type)})
	if err != nil {
		return err
	}

	// Write req id
	reqIDBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(reqIDBytes, msg.RequestID)
	_, err = w.Write(reqIDBytes)
	if err != nil {
		return err
	}

	// Write LENGTH (4 bytes)
	length := uint32(len(msg.Payload))
	lenBytes := make([]byte, 4)

	binary.BigEndian.PutUint32(lenBytes, length)

	_, err = w.Write(lenBytes)
	if err != nil {
		return err
	}

	// Write PAYLOAD
	if len(msg.Payload) > 0 {
		_, err = w.Write(msg.Payload)
	}
	return err
}
