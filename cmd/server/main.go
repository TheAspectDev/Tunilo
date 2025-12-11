package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/TheAspectDev/tunio/internal/protocol"
	"github.com/TheAspectDev/tunio/internal/server"
)

const CONTROL_SERVER_ADDRESS = "0.0.0.0:9090"
const PUBLIC_SERVER_ADDRESS = "0.0.0.0:4311"

var (
	client   net.Conn
	clientMu sync.Mutex
)

func HandleHTTP(w http.ResponseWriter, r *http.Request) {
	clientMu.Lock()
	clientConn := client
	clientMu.Unlock()

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

		server.CopyResponseHeaders(w, response)
		w.WriteHeader(response.StatusCode)

		_, err = io.Copy(w, response.Body)

		if err != nil {
			log.Printf("Error streaming response body: %v", err)
			return
		}
	}

}

func main() {
	go startControlServer()

	http.HandleFunc("/", HandleHTTP)
	err := http.ListenAndServe(PUBLIC_SERVER_ADDRESS, nil)

	if err != nil {
		log.Fatal("HTTP srvr failed:", err)
	}

}

func startControlServer() {
	ln, err := net.Listen("tcp", CONTROL_SERVER_ADDRESS)

	if err != nil {
		log.Fatalf("Failed to start control server: %v", err)
	}

	log.Println("Control server listening on", CONTROL_SERVER_ADDRESS)

	for {
		conn, err := ln.Accept()
		clientMu.Lock()
		client = conn
		clientMu.Unlock()

		if err != nil {
			log.Println("Error accepting client:", err)
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	if err := waitForClientReady(conn); err != nil {
		(conn).Close()
		log.Println(err)
		return
	}

	fmt.Println("GOT READY!!!!!!!!!!!")
}

func waitForClientReady(conn net.Conn) error {
	reader := bufio.NewReader(conn)
	msg, err := protocol.Read(reader)

	if err != nil {
		return fmt.Errorf("failed to read READY: %w", err)
	}

	if !(msg.Type == protocol.MsgReady) {
		return fmt.Errorf("unexpected READY value: %q", msg.Type)
	}

	return nil
}
