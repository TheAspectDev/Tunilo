package client

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/TheAspectDev/tunio/internal/logging"
	"github.com/TheAspectDev/tunio/internal/protocol"
)

type Session struct {
	localClient *http.Client
	controlConn net.Conn
	forward     string

	Logger logging.Logger
}

func NewSession(controlConn net.Conn, localClient *http.Client, forward string) *Session {
	client := &Session{
		controlConn: controlConn,
		localClient: localClient,
		forward:     forward,
	}
	return client
}

func (s *Session) Close() error {
	return s.controlConn.Close()
}

func (s *Session) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		// close connection on ctx close
		s.controlConn.Close()
	}()

	go s.startPingLoop(ctx)

	s.Logger.Logf("Listening to requests...")

	for {
		if err := s.handleControlMessage(); err != nil {
			return err
		}
	}
}

func (s *Session) startPingLoop(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			_ = protocol.Write(s.controlConn, protocol.Message{
				Type:      protocol.MsgPing,
				RequestID: 0,
			})
		}
	}
}
