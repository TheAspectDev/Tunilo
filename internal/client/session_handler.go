package client

import (
	"bufio"
	"bytes"
	"net/http"

	"github.com/TheAspectDev/tunio/protocol"
)

func (s *Session) handleControlMessage() error {
	msg, err := protocol.Read(s.controlConn)
	if err != nil {
		return err
	}

	switch msg.Type {
	case protocol.MsgRequest:
		return s.handleRequest(&msg)
	default:
		return nil
	}
}

func (s *Session) handleRequest(msg *protocol.Message) error {
	reader := bufio.NewReader(bytes.NewReader(msg.Payload))
	req, err := http.ReadRequest(reader)
	if err != nil {
		return err
	}

	s.ForwardRequest(req, msg.RequestID)
	return nil
}
