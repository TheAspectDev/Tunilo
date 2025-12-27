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
	_ = tcp.SetKeepAlive(true)
	_ = tcp.SetKeepAlivePeriod(10 * time.Second)
	_ = tcp.SetNoDelay(true)

}
