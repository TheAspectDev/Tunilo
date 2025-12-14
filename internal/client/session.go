package client

import (
	"context"
	"net"
	"net/http"
)

type Session struct {
	localClient *http.Client
	controlConn net.Conn
	forward     string
}

func NewSession(controlConn net.Conn, localClient *http.Client, forward string) *Session {
	client := &Session{
		controlConn: controlConn,
		localClient: localClient,
		forward:     forward,
	}
	return client
}

func (s *Session) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()

		// close connection on ctx close
		s.controlConn.Close()
	}()

	for {
		if err := s.handleControlMessage(); err != nil {
			return err
		}
	}
}
