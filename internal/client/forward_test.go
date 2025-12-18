package client

import (
	"bytes"
	"math/rand/v2"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TheAspectDev/tunio/protocol"
)

func TestSession_ForwardRequest_Success(t *testing.T) {
	// basic local backend
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))
	defer backend.Close()

	a, b := net.Pipe()
	defer a.Close()
	defer b.Close()

	s := &Session{
		controlConn: b,
		localClient: backend.Client(),
		forward:     backend.URL,
	}

	req := httptest.NewRequest("GET", "http://gettingreplaced.com", nil)

	id := rand.Uint64N(999999)

	go s.ForwardRequest(req, id)

	msg, err := protocol.Read(a)
	if err != nil {
		t.Fatalf("read: %v", err)
	}

	if msg.Type != protocol.MsgResponse {
		t.Fatalf("type %v", msg.Type)
	}

	if msg.RequestID != id {
		t.Fatalf("request id %d", msg.RequestID)
	}

	if !bytes.Contains(msg.Payload, []byte("200 OK")) {
		t.Fatalf("bad response")
	}
}
