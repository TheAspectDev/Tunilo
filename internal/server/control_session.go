package server

import (
	"bufio"
	"errors"
	"net"
	"sync"
	"sync/atomic"

	"github.com/TheAspectDev/tunio/protocol"
)

// Each session gets a single ControlSession
type ControlSession struct {
	conn net.Conn

	pending   map[uint64]chan []byte
	pendingMu sync.Mutex

	writeMu sync.Mutex
	counter atomic.Uint64
	closed  chan struct{}
}

func NewControlSession(conn net.Conn) *ControlSession {
	return &ControlSession{
		conn:    conn,
		pending: make(map[uint64]chan []byte),
		closed:  make(chan struct{}),
	}
}

func (s *ControlSession) Close() {
	s.conn.Close()
}

// Listen for request ids
func (s *ControlSession) Run() {
	reader := bufio.NewReader(s.conn)
	defer s.cleanup()

	for {
		msg, err := protocol.Read(reader)
		if err != nil {
			return
		}

		switch msg.Type {
		case protocol.MsgPing:
			s.writeMu.Lock()
			protocol.Write(s.conn, protocol.Message{
				Type:      protocol.MsgPong,
				RequestID: 0,
			})
			s.writeMu.Unlock()

		case protocol.MsgResponse:
			s.pendingMu.Lock()
			ch := s.pending[msg.RequestID]
			s.pendingMu.Unlock()

			if ch != nil {
				ch <- msg.Payload
			}
		}

	}
}

func (s *ControlSession) cleanup() {
	s.pendingMu.Lock()
	defer s.pendingMu.Unlock()

	for _, ch := range s.pending {
		close(ch)
	}
	close(s.closed)
	s.conn.Close()
}

func (s *ControlSession) Forward(payload []byte) ([]byte, error) {
	id := s.counter.Add(1)
	ch := make(chan []byte, 1)

	// add the request to pending list
	// ControlSession.Run() goroutine will get data from the client and put it inside the reserved s.pending[id] channel
	s.pendingMu.Lock()
	s.pending[id] = ch
	s.pendingMu.Unlock()

	// remove after the response has arrived
	defer func() {
		s.pendingMu.Lock()
		delete(s.pending, id)
		s.pendingMu.Unlock()
	}()

	s.writeMu.Lock()
	err := protocol.Write(s.conn, protocol.Message{
		Type:      protocol.MsgRequest,
		RequestID: id,
		Payload:   payload,
	})
	s.writeMu.Unlock()

	if err != nil {
		return nil, err
	}

	select {
	case resp := <-ch:
		return resp, nil
	case <-s.closed:
		return nil, errors.New("session closed")
	}
}
