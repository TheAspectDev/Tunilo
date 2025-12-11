package main

import (
	"bufio"
	"fmt"
	"log"
	"net"

	"github.com/TheAspectDev/tunio/internal/protocol"
)

const CONTROL_SERVER_ADDRESS = "0.0.0.0:9090"

func main() {
	startControlServer()
}

func startControlServer() {
	ln, err := net.Listen("tcp", CONTROL_SERVER_ADDRESS)

	if err != nil {
		log.Fatalf("Failed to start control server: %v", err)
	}

	log.Println("Control server listening on", CONTROL_SERVER_ADDRESS)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Error accepting client:", err)
			continue
		}
		handleClient(&conn)
	}
}

func handleClient(conn *net.Conn) {
	if err := waitForClientReady(*conn); err != nil {
		(*conn).Close()
		log.Println(err)
		return
	}

	fmt.Println("GOT READY!!!!!!!!!!!")

}

func waitForClientReady(conn net.Conn) error {
	reader := bufio.NewReader(conn)
	msg, err := protocol.Read(reader)

	if err != nil {
		return fmt.Errorf("failed to read READY: %w", err)
	}

	if !(msg.Type == protocol.MsgReady) {
		return fmt.Errorf("unexpected READY value: %q", msg.Type)
	}

	return nil
}
