package server

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
)

var inflight = make(chan struct{}, 200)

func (srv *Server) HandleHTTP(w http.ResponseWriter, r *http.Request) {
	var RequestBuffer bytes.Buffer

	// Heavy load
	select {
	case inflight <- struct{}{}:
		defer func() { <-inflight }()
	default:
		http.Error(w, "Tunnel busy", http.StatusServiceUnavailable)
		return
	}

	if err := r.Write(&RequestBuffer); err != nil {
		http.Error(w, "Failed to encode request", http.StatusInternalServerError)
		return
	}

	// read this functions comment
	session := srv.getAnySession()
	if session == nil {
		http.Error(w, "no tunnel client connected", http.StatusServiceUnavailable)
		return
	}

	payload, err := session.Forward(RequestBuffer.Bytes())
	if err != nil {
		http.Error(w, "tunnel error", http.StatusBadGateway)
		return
	}

	reader := bufio.NewReader(bytes.NewReader(payload))
	response, err := http.ReadResponse(reader, r)
	if err != nil {
		return
	}

	CopyResponseHeaders(w, response)
	w.WriteHeader(response.StatusCode)

	_, err = io.Copy(w, response.Body)
	if err != nil {
		return
	}
}
