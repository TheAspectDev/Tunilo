package client

import (
	"context"
	"net"
	"net/http"

	"github.com/TheAspectDev/tunio/internal/logging"
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

	s.Logger.Logf("Listening to requests...")

	for {
		if err := s.handleControlMessage(); err != nil {
			return err
		}
	}
}
