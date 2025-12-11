package main

import (
	"log"
	"net"

	"github.com/TheAspectDev/tunio/internal/protocol"
)

const CONTROL_SERVER_ADDRESS = "0.0.0.0:9090"

func main() {
	conn, err := net.Dial("tcp", CONTROL_SERVER_ADDRESS)
	if err != nil {
		log.Println("error while connecting to ", CONTROL_SERVER_ADDRESS)
	}
	defer conn.Close()

	protocol.Write(conn, protocol.Message{
		Type: protocol.MsgRequest,
	})

}
