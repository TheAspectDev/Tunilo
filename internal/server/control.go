package server

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"log"
	"net"

	"github.com/TheAspectDev/tunio/protocol"
)

func (srv *Server) StartControlServer() {
	var (
		ln  net.Listener
		err error
	)

	if srv.tls != nil {
		pair, err := tls.LoadX509KeyPair(srv.tls.Cert, srv.tls.Key)
		if err != nil {
			log.Fatalf("Failed to load keypair: %v", err)
		}

		cfg := &tls.Config{
			MinVersion:   tls.VersionTLS13,
			Certificates: []tls.Certificate{pair},
		}

		ln, err = tls.Listen("tcp", srv.ControlAddress, cfg)
	} else {
		ln, err = net.Listen("tcp", srv.ControlAddress)
	}

	if err != nil {
		log.Fatalf("Failed to start control server: %v", err)
		return
	}

	for {
		conn, err := ln.Accept()
		protocol.EnableTCPKeepalive(conn)
		if err != nil {
			log.Println("Error accepting client:", err)
			continue
		}

		go srv.handleNewClient(conn)
	}
}

func (srv *Server) handleNewClient(conn net.Conn) {
	if err := srv.waitForClientReady(conn); err != nil {
		conn.Close()
		return
	}

	session := NewControlSession(conn)
	key := conn.RemoteAddr().String()

	srv.SessionsMu.Lock()
	srv.Sessions[key] = session
	srv.SessionsMu.Unlock()

	go func(conn net.Conn) {
		defer func() {
			conn.Close()

			srv.SessionsMu.Lock()
			delete(srv.Sessions, key)
			srv.SessionsMu.Unlock()
		}()

		session.Run()
	}(conn)
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
