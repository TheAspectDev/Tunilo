package client

import (
	"bufio"
	"bytes"
	"log"
	"net/http"

	"github.com/TheAspectDev/tunio/internal/protocol"
)

func (client *Client) HandleMessage() {
	msg, err := protocol.Read(client.controlServer)
	if err != nil {
		log.Println("error reading message", err)
		return
	}

	if msg.Type == protocol.MsgRequest {
		reader := bufio.NewReader(bytes.NewReader(msg.Payload))
		request, err := http.ReadRequest(reader)

		if err != nil {
			log.Println("error processing request:", err)
			return
		}

		client.ForwardRequest(request, msg.RequestID)
	}
}
