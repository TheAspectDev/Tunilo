package protocol

import (
	"net"
	"time"
)

func EnableTCPKeepalive(conn net.Conn) {
	tcp, ok := conn.(*net.TCPConn)
	if !ok {
		return
	}
	tcp.SetKeepAlive(true)
	tcp.SetKeepAlivePeriod(10 * time.Second)
	tcp.SetNoDelay(true)
}
