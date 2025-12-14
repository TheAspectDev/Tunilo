package server

import (
	"sync"
)

type Server struct {
	serverAddress  string
	password       string
	controlAddress string

	sessionsMu sync.RWMutex
	sessions   map[string]*ControlSession
}

func NewServer(serverAddress string, controlAddress string, password string) *Server {
	return &Server{
		serverAddress:  serverAddress,
		controlAddress: controlAddress,
		sessions:       make(map[string]*ControlSession),
		password:       password,
	}
}

// picks the first server
// reserved for later use
// ( like for sending certain requests to certain clients )
func (srv *Server) getAnySession() *ControlSession {
	srv.sessionsMu.RLock()
	defer srv.sessionsMu.RUnlock()

	for _, s := range srv.sessions {
		return s
	}
	return nil
}
