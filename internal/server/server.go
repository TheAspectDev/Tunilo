package server

import (
	"net"
	"sync"
)

type Server struct {
	serverAddress  string
	controlAddress string
	client         net.Conn
	clientMu       sync.Mutex
}

func NewServer(serverAddress string, controlAddress string) *Server {
	return &Server{
		serverAddress:  serverAddress,
		controlAddress: controlAddress,
	}
}
