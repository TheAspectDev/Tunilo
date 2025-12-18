package server

import (
	"net"
	"testing"
	"time"

	"github.com/TheAspectDev/tunio/protocol"
)

func TestControlSession_RunRespondsToPing(t *testing.T) {
	a, b := net.Pipe()
	defer a.Close()
	defer b.Close()

	a.SetDeadline(time.Now().Add(5 * time.Second))
	b.SetDeadline(time.Now().Add(5 * time.Second))

	controlSession := NewControlSession(a)
	go controlSession.Run()

	if err := protocol.Write(b, protocol.Message{
		Type:      protocol.MsgPing,
		RequestID: 0,
	}); err != nil {
		t.Fatalf("b write error: %v", err)
		return
	}

	data, err := protocol.Read(b)
	if err != nil {
		t.Fatalf("b read error: %v", err)
		return
	}

	if data.Type != protocol.MsgPong {
		t.Fatalf("didnt get pong %v", data.Type)
	}
}

func TestControlSession_ForwardMatchesResponse(t *testing.T) {
	a, b := net.Pipe()
	defer a.Close()
	defer b.Close()

	a.SetDeadline(time.Now().Add(4 * time.Second))
	b.SetDeadline(time.Now().Add(4 * time.Second))

	controlSession := NewControlSession(a)
	go controlSession.Run()

	resultCh := make(chan []byte, 1)
	errCh := make(chan error, 1)

	go func() {
		resp, err := controlSession.Forward([]byte("req"))
		if err != nil {
			errCh <- err
			return
		}
		resultCh <- resp
	}()

	reqMsg, err := protocol.Read(b)
	if err != nil {
		t.Fatalf("read request: %v", err)
	}
	if reqMsg.Type != protocol.MsgRequest || string(reqMsg.Payload) != "req" {
		t.Fatalf("invalid request msg: %+v", reqMsg)
	}

	if err := protocol.Write(b, protocol.Message{
		Type:      protocol.MsgResponse,
		RequestID: reqMsg.RequestID,
		Payload:   []byte("resp"),
	}); err != nil {
		t.Fatalf("write response: %v", err)
	}

	select {
	case err := <-errCh:
		t.Fatalf("Forward gave an err: %v", err)
	case got := <-resultCh:
		if string(got) != "resp" {
			t.Fatalf("expected response got %q", string(got))
		}
	case <-time.After(4 * time.Second):
		t.Fatalf("timed out")
	}
}
