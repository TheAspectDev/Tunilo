package server

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync/atomic"

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

	id := atomic.AddUint64(&srv.counter, 1)

	respChan := make(chan []byte, 1)

	srv.pendingMu.Lock()
	srv.pending[id] = respChan
	srv.pendingMu.Unlock()

	srv.writeMu.Lock()

	// todo
	_ = protocol.Write(clientConn, protocol.Message{
		Type:      protocol.MsgRequest,
		RequestID: id,
		Payload:   RequestBuffer.Bytes(),
	})
	srv.writeMu.Unlock()

	payload := <-respChan

	srv.pendingMu.Lock()
	delete(srv.pending, id)
	srv.pendingMu.Unlock()

	reader := bufio.NewReader(bytes.NewReader(payload))
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
