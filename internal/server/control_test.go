package server

import (
	"net"
	"testing"
	"time"

	"github.com/TheAspectDev/tunio/protocol"
)

func TestWaitForClientReady_AcceptsCorrectPassword(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	defer serverConn.Close()
	defer clientConn.Close()

	_ = serverConn.SetDeadline(time.Now().Add(2 * time.Second))
	_ = clientConn.SetDeadline(time.Now().Add(2 * time.Second))

	go func() {
		_ = protocol.Write(clientConn, protocol.Message{
			Type:      protocol.MsgReady,
			RequestID: 0,
			Payload:   []byte("password12345"),
		})
	}()

	srv := &Server{password: "password12345"}
	if err := srv.waitForClientReady(serverConn); err != nil {
		t.Fatalf("expected nil err got %v", err)
	}
}

func TestWaitForClientReady_RejectsWrongPassword(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	defer serverConn.Close()
	defer clientConn.Close()

	_ = serverConn.SetDeadline(time.Now().Add(2 * time.Second))
	_ = clientConn.SetDeadline(time.Now().Add(2 * time.Second))

	go func() {
		_ = protocol.Write(clientConn, protocol.Message{
			Type:      protocol.MsgReady,
			RequestID: 0,
			Payload:   []byte("wrongpassword12345"),
		})
	}()

	srv := &Server{password: "password12345"}
	if err := srv.waitForClientReady(serverConn); err == nil {
		t.Fatalf("expected error got nil")
	}
}
