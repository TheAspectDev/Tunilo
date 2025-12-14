package server

import (
	"net"
	"sync"
)

type Server struct {
	serverAddress  string
	password       string
	controlAddress string
	client         net.Conn
	clientMu       sync.RWMutex

	pending   map[uint64]chan []byte
	pendingMu sync.Mutex
	counter   uint64
}

func NewServer(serverAddress string, controlAddress string, password string) *Server {
	return &Server{
		serverAddress:  serverAddress,
		controlAddress: controlAddress,
		pending:        make(map[uint64]chan []byte),
		password:       password,
	}
}
