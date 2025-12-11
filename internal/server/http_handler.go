package server

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/TheAspectDev/tunio/internal/protocol"
)

func (srv *Server) HandleHTTP(w http.ResponseWriter, r *http.Request) {
	srv.clientMu.Lock()
	clientConn := srv.client
	srv.clientMu.Unlock()

	var RequestBuffer bytes.Buffer

	if err := r.Write(&RequestBuffer); err != nil {
		log.Printf("Failed to serialize HTTP request: %v", err)
		http.Error(w, "Failed to encode request", http.StatusInternalServerError)
		return
	}

	protocol.Write(clientConn, protocol.Message{
		Type:    protocol.MsgRequest,
		Payload: RequestBuffer.Bytes(),
	})

	msg, err := protocol.Read(clientConn)

	if msg.Type == protocol.MsgResponse {
		if err != nil {
			fmt.Println("error reading response", err)
			return
		}

		reader := bufio.NewReader(bytes.NewReader(msg.Payload))
		response, err := http.ReadResponse(reader, r)

		if err != nil {
			log.Printf("error parsing response %v", err)
			return
		}

		fmt.Println(response.Status)

		CopyResponseHeaders(w, response)
		w.WriteHeader(response.StatusCode)

		_, err = io.Copy(w, response.Body)

		if err != nil {
			log.Printf("Error streaming response body: %v", err)
			return
		}
	}
}
