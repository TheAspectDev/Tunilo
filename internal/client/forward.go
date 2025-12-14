package client

import (
	"bytes"
	"log"
	"net/http"
	"strings"

	"github.com/TheAspectDev/tunio/internal/protocol"
)

func (session *Session) ForwardRequest(req *http.Request, req_id uint64) {
	forwardData := strings.Split(session.forward, "://")

	// something.com
	req.URL.Host = forwardData[1]
	// http or https
	req.URL.Scheme = forwardData[0]
	// something.com
	req.Host = forwardData[1]

	req.RequestURI = ""

	localResp, err := session.localClient.Do(req)

	if err != nil {
		log.Printf("Error forwarding request to local app: %v", err)
		protocol.Write(session.controlConn, protocol.Message{
			Type:      protocol.MsgResponse,
			RequestID: req_id,
			Payload:   []byte("HTTP/1.1 503 Service Unavailable\r\nContent-Length: 0\r\n\r\n"),
		})
		return
	}

	defer localResp.Body.Close()

	var RequestBuffer bytes.Buffer
	if err := localResp.Write(&RequestBuffer); err != nil {
		log.Printf("Failed to serialize HTTP response: %v", err)
		return
	}

	protocol.Write(session.controlConn, protocol.Message{
		Type:      protocol.MsgResponse,
		RequestID: req_id,
		Payload:   RequestBuffer.Bytes(),
	})

}
