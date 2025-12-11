package server

import (
	"bufio"
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

	log.Println("Control server listening on", srv.controlAddress)

	for {
		conn, err := ln.Accept()
		srv.clientMu.Lock()
		srv.client = conn
		srv.clientMu.Unlock()

		if err != nil {
			log.Println("Error accepting client:", err)
			continue
		}
		go srv.initClient(conn)
	}
}

func (srv *Server) initClient(conn net.Conn) {
	if err := waitForClientReady(conn); err != nil {
		conn.Close()
		fmt.Println("Handshake failed:", err)
		return
	}

	fmt.Println("Client ready")

	srv.clientMu.Lock()
	srv.client = conn
	srv.clientMu.Unlock()
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
