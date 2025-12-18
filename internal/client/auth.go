package client

import (
	"bufio"
	"bytes"

	"github.com/TheAspectDev/tunio/protocol"
)

func (session *Session) Authenticate(password string) error {
	var passBuffer bytes.Buffer
	writer := bufio.NewWriter(&passBuffer)
	writer.WriteString(password)
	writer.Flush()

	session.Logger.Logf("Authenticating...")

	return protocol.Write(session.controlConn, protocol.Message{
		Type:      protocol.MsgReady,
		Payload:   passBuffer.Bytes(),
		RequestID: 0,
	})
}
