package server

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"

	"github.com/TheAspectDev/tunio/internal/protocol"
)

func (srv *Server) StartControlServer() {
	ln, err := net.Listen("tcp", srv.controlAddress)

	if err != nil {
		log.Fatalf("Failed to start control server: %v", err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Error accepting client:", err)
			continue
		}

		srv.clientMu.Lock()
		srv.client = conn
		srv.clientMu.Unlock()

		go srv.initClient(conn)
	}
}

func (srv *Server) initClient(conn net.Conn) {
	if err := srv.waitForClientReady(conn); err != nil {
		conn.Close()
		return
	}

	srv.clientMu.Lock()
	srv.client = conn
	srv.clientMu.Unlock()

	go srv.tunnelReader(conn)
}

func (srv *Server) tunnelReader(conn net.Conn) {
	reader := bufio.NewReader(conn)

	for {
		msg, err := protocol.Read(reader)
		if err != nil {
			return
		}

		if msg.Type == protocol.MsgResponse {
			srv.pendingMu.Lock()
			ch := srv.pending[msg.RequestID]
			srv.pendingMu.Unlock()

			if ch != nil {
				ch <- msg.Payload
			}
		}
	}
}

func (srv *Server) waitForClientReady(conn net.Conn) error {
	reader := bufio.NewReader(conn)
	msg, err := protocol.Read(reader)

	if err != nil {
		return fmt.Errorf("failed to read READY: %w", err)
	}

	if !(msg.Type == protocol.MsgReady) {
		return fmt.Errorf("unexpected READY value: %q", msg.Type)
	}

	var passBuffer bytes.Buffer
	writer := bufio.NewWriter(&passBuffer)

	if _, err := writer.WriteString(srv.password); err != nil {
		return fmt.Errorf("failed to serialize password: %w", err)
	}
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to serialize password: %w", err)
	}

	if !bytes.Equal(msg.Payload, passBuffer.Bytes()) {
		return fmt.Errorf("password mismatch")
	}

	return nil
}
