package client

import (
	"bufio"
	"bytes"

	"github.com/TheAspectDev/tunio/internal/protocol"
)

func (client *Client) Authenticate(password *string) {
	var passBuffer bytes.Buffer
	writer := bufio.NewWriter(&passBuffer)
	writer.WriteString(*password)
	writer.Flush()

	protocol.Write(client.controlServer, protocol.Message{
		Type:      protocol.MsgReady,
		Payload:   passBuffer.Bytes(),
		RequestID: 0,
	})
}
