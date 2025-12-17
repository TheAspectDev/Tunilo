package server

import (
	"errors"
	"log"
	"net/http"
	"sync"
)

type TLSConfig struct {
	Cert string
	Key  string
}

type Server struct {
	ServerAddress  string
	password       string
	ControlAddress string

	SessionsMu sync.RWMutex
	Sessions   map[string]*ControlSession

	tls *TLSConfig
}

type ServerConfig struct {
	serverAddress  string
	password       string
	controlAddress string
	tls            *TLSConfig
}

func NewServerBuilder() *ServerConfig {
	return &ServerConfig{}
}

func (b *ServerConfig) SetAddress(address string) *ServerConfig {
	b.serverAddress = address
	return b
}

func (b *ServerConfig) SetPassword(password string) *ServerConfig {
	b.password = password
	return b
}

func (b *ServerConfig) SetControlAddress(address string) *ServerConfig {
	b.controlAddress = address
	return b
}

func (b *ServerConfig) SetTLS(config TLSConfig) *ServerConfig {
	b.tls = &TLSConfig{Cert: config.Cert, Key: config.Key}
	return b
}

func (b *ServerConfig) Build() (*Server, error) {
	if b.serverAddress == "" {
		return nil, errors.New("server address required")
	}

	if b.controlAddress == "" {
		return nil, errors.New("control address required")
	}

	if b.tls != nil && (b.tls.Cert == "" || b.tls.Key == "") {
		return nil, errors.New("cert or key missing, use -insecure flag to run without tls")
	}

	return &Server{
		ServerAddress:  b.serverAddress,
		ControlAddress: b.controlAddress,
		password:       b.password,
		tls:            b.tls,
		Sessions:       make(map[string]*ControlSession),
	}, nil
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

// ListenAndServeTLS must be done through reverse-proxies like nginx/...
// the tls is only used for secure connection through the tunnel
func (srv *Server) StartPublicServer(httpServer *http.Server) {
	// if srv.tls == nil {
	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server failed: %v", err)
		return
	}
	// } else {
	// 	if err := httpServer.ListenAndServeTLS(srv.tls.Cert, srv.tls.Key); err != http.ErrServerClosed {
	// 		log.Fatalf("HTTP server failed: %v", err)
	// 		return
	// 	}
	// }
}
