package server

import (
	"sync"
)

type Server struct {
	ServerAddress  string
	password       string
	ControlAddress string

	SessionsMu sync.RWMutex
	Sessions   map[string]*ControlSession
}

func NewServer(serverAddress string, controlAddress string, password string) *Server {
	return &Server{
		ServerAddress:  serverAddress,
		ControlAddress: controlAddress,
		Sessions:       make(map[string]*ControlSession),
		password:       password,
	}
}

// picks the first server
// reserved for later use
// ( like for sending certain requests to certain clients )
func (srv *Server) getAnySession() *ControlSession {
	srv.SessionsMu.RLock()
	defer srv.SessionsMu.RUnlock()

	for _, s := range srv.Sessions {
		return s
	}
	return nil
}
